package ssatest

import (
	"fmt"
	"os"
)

/*
a = 1
b = 1
a = b
*/

func Optimise(ops []Operation) []Operation {
	assignMap := map[string]string{}
	assignments := map[string]int{}

	ssaOps := []Operation{}

	for _, op := range ops[:len(ops)-1] {
		ssaName := fmt.Sprintf("%s-%d", op.Destination, assignments[op.Destination])
		assignMap[ssaName] = op.Destination
		assignments[op.Destination]++

		// Handle single-value case
		if len(op.Values) == 1 {
			val := getValueSSA(op.Values[0])
			if val.Reference != "" {
				offset := 0
				if val.Reference == op.Destination {
					offset = 1
				}

				refAssigns, ok := assignments[val.Reference]
				if !ok {
					fmt.Println("Invalid assignment:", val.Reference)
					os.Exit(1)
				}
				refSsaName := fmt.Sprintf("%s-%d", val.Reference, refAssigns-1-offset)

				ssaOps = append(ssaOps, Operation{
					Type:        op.Type,
					Destination: ssaName,
					Values:      []string{refSsaName},
				})
			} else {
				ssaOps = append(ssaOps, Operation{
					Type:        op.Type,
					Destination: ssaName,
					Values:      []string{op.Values[0]},
				})
			}

			continue
		}

		// Handle dual-value case
		v1 := getValueSSA(op.Values[0])
		v2 := getValueSSA(op.Values[1])
		values := []string{}

		if v1.Reference != "" {
			offset := 0
			if v1.Reference == op.Destination {
				offset = 1
			}

			refAssigns, ok := assignments[v1.Reference]
			if !ok {
				fmt.Println("Invalid assignment:", v1.Reference)
				os.Exit(1)
			}
			values = append(values, fmt.Sprintf("%s-%d", v1.Reference, refAssigns-1-offset))
		} else {
			values = append(values, op.Values[0])
		}

		if v2.Reference != "" {
			offset := 0
			if v2.Reference == op.Destination {
				offset = 1
			}

			refAssigns, ok := assignments[v2.Reference]
			if !ok {
				fmt.Println("Invalid assignment:", v2.Reference)
				os.Exit(1)
			}
			values = append(values, fmt.Sprintf("%s-%d", v2.Reference, refAssigns-1-offset))
		} else {
			values = append(values, op.Values[1])
		}

		ssaOps = append(ssaOps, Operation{
			Type:        op.Type,
			Destination: ssaName,
			Values:      values,
		})
	}

	if os.Getenv("SSA_DUMP") == "1" {
		data := ""
		for _, op := range ssaOps {
			data += fmt.Sprintf("%s\n", op.String())
		}
		os.WriteFile("SSA_DUMP", []byte(data), 0644)
	}

	use := ops[len(ops)-1]
	usedNames := map[string]any{}

	for _, name := range use.Values {
		refAssigns, ok := assignments[name]
		if !ok {
			fmt.Println("Attempt to use name that doesn't exist:", name)
			os.Exit(1)
		}

		usedNames[fmt.Sprintf("%s-%d", name, refAssigns-1)] = true
	}

	for i := len(ssaOps) - 1; i >= 0; i-- {
		op := ssaOps[i]
		if _, ok := usedNames[op.Destination]; !ok {
			continue
		}

		for _, val := range op.Values {
			resolved := getValueSSA(val)
			if resolved.Reference != "" {
				usedNames[resolved.Reference] = true
			}
		}
	}

	optimisedOps := []Operation{}

	for _, op := range ssaOps {
		if _, ok := usedNames[op.Destination]; ok {
			values := []string{}

			for _, val := range op.Values {
				resolved := getValueSSA(val)
				if resolved.Reference != "" {
					values = append(values, assignMap[resolved.Reference])
				} else {
					values = append(values, val)
				}
			}

			optimisedOps = append(optimisedOps, Operation{
				Type:        op.Type,
				Destination: assignMap[op.Destination],
				Values:      values,
			})
		}
	}

	return append(optimisedOps, ops[len(ops)-1])
}
