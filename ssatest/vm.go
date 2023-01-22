package ssatest

import (
	"fmt"
	"os"
)

func recurseValue(values map[string]Value, value Value) Value {
	if value.Reference == "" {
		return value
	}

	if _, ok := values[value.Reference]; !ok {
		fmt.Println("Reference not found:", value.Reference)
		os.Exit(1)
	}

	return recurseValue(values, values[value.Reference])
}

func Run(ops []Operation) {
	values := map[string]Value{}
	assignments := map[string]int{}

	use := ops[len(ops)-1]
	useMap := map[string]any{}

	for _, name := range use.Values {
		useMap[name] = nil
	}

	for _, op := range ops[:len(ops)-1] {
		assignments[op.Destination]++

		if op.Type == "" {
			values[op.Destination] = getValue(op.Values[0])
			continue
		}

		v1 := recurseValue(values, getValue(op.Values[0]))
		v2 := recurseValue(values, getValue(op.Values[1]))

		switch op.Type {
		case OpTypeAdd:
			values[op.Destination] = Value{Value: v1.Value + v2.Value}
		case OpTypeSub:
			values[op.Destination] = Value{Value: v1.Value - v2.Value}
		case OpTypeMul:
			values[op.Destination] = Value{Value: v1.Value * v2.Value}
		case OpTypeDiv:
			values[op.Destination] = Value{Value: v1.Value / v2.Value}
		}
	}

	fmt.Printf("Final output (%d assignments):\n", len(ops)-1)

	for k, v := range values {
		if _, ok := useMap[k]; !ok {
			continue
		}

		var rep string

		if v.Reference != "" {
			rv := recurseValue(values, v)
			rep = fmt.Sprintf("%s (%d)", v.Reference, rv.Value)
		} else {
			rep = fmt.Sprintf("%d", v.Value)
		}

		fmt.Printf("%s = %-8s with %d assignments\n", k, rep, assignments[k])
	}
}
