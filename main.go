// Copyright (c) 2022, Max Fedorov <mail@skam.in>

package main

import (
	"fmt"
	"os"
)

const (
	Version = "1.0.0"
	Author  = "Max Fedorov <mail@skam.in>"
)

func main() {
	params := parseCliArguments()
	if params == nil {
		fmt.Fprintln(os.Stderr, "ERROR: missing required arguments")
		os.Exit(1)
	}

	configHcl := ReadHclFile(*params.input)
	templateHcl := ReadHclFile(*params.template)
	outputHcl := Process(configHcl, templateHcl)

	if *params.output == "" {
		fmt.Println(string(outputHcl.Bytes()))
	} else {
		err := os.WriteFile(*params.output, outputHcl.Bytes(), 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
