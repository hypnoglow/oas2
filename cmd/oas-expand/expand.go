// CLI utility that expands OAS file to reduce init time.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/hypnoglow/oas2"
)

const help = `Expand OpenAPI specification

Usage:
    oas-expand [FLAGS] <SPEC_FILE>

Flags:
    -h, -help          Print help message
    -t, -target-dir    Save expanded spec to diretory
`

func main() {
	flagHelp := flag.Bool("help", false, "Print help message")
	flagHelpShort := flag.Bool("h", false, "Print help message")
	flagTargetDir := flag.String("target-dir", "", "Output directory for expanded spec files")
	flagTargetDirShort := flag.String("t", "", "Output directory for expanded spec files")
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

	targetDir := *flagTargetDir
	if *flagTargetDirShort != "" {
		targetDir = *flagTargetDirShort
	}

	// Save expanded spec to file in dir
	if targetDir != "" {
		if _, err := oas.LoadSpec(specFile, oas.CacheDir(targetDir)); err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
		fmt.Println("Spec expanded successfully")
		os.Exit(0)
	}

	// Print expanded spec to stdout
	document, err := oas.LoadSpec(specFile)
	if err != nil {
		fmt.Printf("Error: %s\n\n", err)
		os.Exit(1)
	}
	if err := json.NewEncoder(os.Stdout).Encode(document.Spec()); err != nil {
		fmt.Printf("Error: %s\n\n", err)
		os.Exit(1)
	}
}
