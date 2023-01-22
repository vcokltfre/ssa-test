package ssatest

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

type OpType string

const (
	OpTypeAdd OpType = "add"
	OpTypeSub OpType = "sub"
	OpTypeMul OpType = "mul"
	OpTypeDiv OpType = "div"
	OpTypeUse OpType = "use"
)

var (
	lineRegex  = regexp.MustCompile(`^\s*([a-z])+\s*=\s*([a-z]+|[0-9]+)(\s*([\/\*+-])\s*([a-z]+|[0-9]+))?\s*$`)
	emptyRegex = regexp.MustCompile(`^\s*$`)

	opTypes = map[string]OpType{
		"+": OpTypeAdd,
		"-": OpTypeSub,
		"*": OpTypeMul,
		"/": OpTypeDiv,
	}
	opTypesInverse = map[OpType]string{
		OpTypeAdd: "+",
		OpTypeSub: "-",
		OpTypeMul: "*",
		OpTypeDiv: "/",
	}
)

type Operation struct {
	Type        OpType
	Destination string
	Values      []string
}

func (o Operation) String() string {
	var output string

	if o.Type == OpTypeUse {
		output = fmt.Sprintf("use %s", strings.Join(o.Values, " "))
	} else if len(o.Values) == 1 {
		output = fmt.Sprintf("%s = %s", o.Destination, o.Values[0])
	} else {
		output = fmt.Sprintf("%s = %s %s %s", o.Destination, o.Values[0], opTypesInverse[o.Type], o.Values[1])
	}

	return strings.TrimSpace(output)
}

func Parse(source string) ([]Operation, error) {
	lines := strings.Split(source, "\n")
	ops := []Operation{}

	for _, line := range lines {
		if strings.HasPrefix(line, "use") {
			parts := strings.Split(line[3:], " ")
			cleanParts := []string{}

			for _, part := range parts {
				cleanPart := strings.TrimSpace(part)
				if cleanPart == "" {
					continue
				}

				cleanParts = append(cleanParts, cleanPart)
			}

			ops = append(ops, Operation{
				Type:   OpTypeUse,
				Values: cleanParts,
			})
			break
		}

		if emptyRegex.MatchString(line) {
			continue
		}

		matches := lineRegex.FindAllStringSubmatch(line, -1)[0][1:]
		if len(matches) == 0 {
			return nil, fmt.Errorf("invalid line: %s", line)
		}

		ops = append(ops, Operation{
			Type:        opTypes[matches[3]],
			Destination: matches[0],
			Values:      []string{matches[1], matches[4]},
		})
	}

	if len(ops) == 0 || ops[len(ops)-1].Type != OpTypeUse {
		fmt.Println("Final operation must be use.")
		os.Exit(1)
	}

	return ops, nil
}
