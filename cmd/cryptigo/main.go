package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Jouini-Mohamed-Chaker/cryptigo/internal/crypt"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: cli <encrypt|decrypt> -in <file> [-out <file>] [-pass <str>] [-overwrite]")
		os.Exit(1)
	}

	subcommand := os.Args[1]
	if subcommand != "encrypt" && subcommand != "decrypt" {
		fmt.Println("Error: expected 'encrypt' or 'decrypt' subcommand")
		os.Exit(1)
	}

	cmd := flag.NewFlagSet(subcommand, flag.ExitOnError)
	inFile := cmd.String("in", "", "Input file")
	outFile := cmd.String("out", "output.enc", "")
	pass := cmd.String("pass", "", "Password")
	overwrite := cmd.Bool("overwrite", false, "Overwrite exisitng file with same output file name")
	cleanup := cmd.Bool("cleanup", false, "Delete input after completion")

	cmd.Parse(os.Args[2:])

	if *inFile == "" {
		fatal(fmt.Errorf("-in is required"))
	}

	if err := crypt.ValidateInputFile(*inFile); err != nil {
		fatal(err)
	}

	if err := crypt.ValidateOutputFile(*outFile, *overwrite); err != nil {
		fatal(err)
	}

	crypt.PromptForPassword(pass)

	switch subcommand {
	case "encrypt":
		if err := crypt.Encrypt(*inFile, *outFile, *pass); err != nil {
			fatal(err)
		}
	case "decrypt":
		if err := crypt.Decrypt(*inFile, *outFile, *pass); err != nil {
			fatal(err)
		}
	}

	if *cleanup {
		err := os.Remove(*inFile)
		if err != nil {
			fmt.Printf("Failed to remove input file: %v\n", err)
		} else {
			fmt.Println("Input file removed.")
		}
	}
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}
