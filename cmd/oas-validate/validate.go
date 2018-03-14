// CLI utility for validating OAS files.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/hypnoglow/oas2"
)

const help = `Validate OpenAPI specification

Usage:
    oas-validate [FLAGS] <SPEC_FILE>

Flags:
    -h, -help          Print help message
`

func main() {
	flagHelp := flag.Bool("help", false, "Print help message")
	flagHelpShort := flag.Bool("h", false, "Print help message")
	flag.Parse()

	if *flagHelp || *flagHelpShort {
		fmt.Println(help)
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) != 1 {
		fmt.Println(help)
		os.Exit(1)
	}
	specFile := args[0]

	if err := oas.ValidateSpec(specFile); err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("Successfully validated")
}
