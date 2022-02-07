// Copyright (c) 2022, Max Fedorov <mail@skam.in>

package main

import (
	"fmt"
	"github.com/akamensky/argparse"
	"os"
)

type Args struct {
	input      *string
	template   *string
	output     *string
	version    *bool
	updateArgs *bool
}

func parseCliArguments() *Args {
	args := new(Args)
	parser := argparse.NewParser("hclmergetool",
		`Works with HashiCorp Configuration Language (HCL) files. `+
			`Allows you to supplement the input file with blocks and attributes `+
			`from the template file. You can use any file name (file extension doesn't matter) `+
			`for the input and template files. The only thing required is valid HCL syntax.`)
	args.input = parser.String("i", "input", &argparse.Options{Required: false, Help: "path to HCL input file"})
	args.template = parser.String("t", "template", &argparse.Options{Required: false, Help: "path to HCL template file"})
	args.output = parser.String("o", "output",
		&argparse.Options{
			Required: false,
			Help:     "path to HCL output file. If not set, print to stdout",
		})
	args.updateArgs = parser.Flag("u", "update-args",
		&argparse.Options{
			Required: false,
			Help: `if set, the values for existing arguments will be updated from ` +
				`the template values, otherwise, the default value remains`,
		})
	args.version = parser.Flag("v", "version", &argparse.Options{Required: false, Help: "show version"})

	err := parser.Parse(os.Args)
	//fmt.Printf("%#v\n", args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	if *args.version {
		fmt.Printf("Version: %s\n", Version)
		os.Exit(0)
	}
	if *args.input == "" {
		fmt.Print(parser.Usage(err))
		fmt.Println("ERROR: required argument -i or --input not found")
		os.Exit(1)
	}
	if *args.template == "" {
		fmt.Print(parser.Usage(err))
		fmt.Println("ERROR: required argument -t or --template not found")
		os.Exit(1)
	}

	return args
}
