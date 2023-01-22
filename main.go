package main

import (
	"fmt"
	"os"
	"ssa-test/ssatest"
)

func main() {
	cmd := os.Args[1]
	file := os.Args[2]

	data, err := os.ReadFile(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	switch cmd {
	case "optimise":
		ops, err := ssatest.Parse(string(data))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		optimisedOps := ssatest.Optimise(ops)
		output := ""

		for _, op := range optimisedOps {
			output += op.String() + "\n"
		}

		err = os.WriteFile(os.Args[3], []byte(output), 0644)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case "run":
		ops, err := ssatest.Parse(string(data))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		ssatest.Run(ops)
	default:
		fmt.Println("Unknown command")
		os.Exit(1)
	}
}
