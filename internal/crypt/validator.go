package crypt

import (
	"errors"
	"fmt"
	"os"
	"unicode"
)

// ValidateInputFile checks whether path exists, is a regular file, and is readable.
// It returns an error describing why the file is invalid, or nil if valid.
func ValidateInputFile(filename string) error {
	if filename == "" {
		return fmt.Errorf("input file path cannot be empty")
	}

	info, err := os.Stat(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("path %q does not exist", filename)
		} else {
			return fmt.Errorf("stat path: %w", err)
		}
	}

	if info.IsDir(){
		return fmt.Errorf("path %q is a directory, expected a regular file", filename)	
	}

	if !info.Mode().IsRegular() {
		return fmt.Errorf("path %q is not a regular file", filename)
	}

	file, err := os.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	file.Close()

	return nil
}

// ValidateOutputFile checks whether file exists (to avoid accidental overwriting) 
// and that the program has proper write permissions.
// It returns an error describing why the file is invalid, or nil if valid.
func ValidateOutputFile(filename string, overwrite bool) error {
	// No need to check if filename is empty as default value of the out flag is set to "output.enc"

	_, err := os.Stat(filename)
	// File exists and overwrite is disabled
	if err == nil && !overwrite {
		return fmt.Errorf("file %q already exists and the 'overwrite' flag is not set", filename)
	}

	// the pipe "|" is an OR operator
	flags := os.O_WRONLY | os.O_CREATE

	if overwrite {
		// O_TRUNC clears existing files
		flags |= os.O_TRUNC
	}

	// Attempts to create or open the output file for writing, truncating it if 
	// overwrite is enabled, to confirm permissions and path validity in one step.
	file, err := os.OpenFile(filename, flags, 0666)
	if err != nil {
		return fmt.Errorf("cannot write to output path %q: %w", filename, err)
	}

	file.Close()
	os.Remove(filename)
	
	return nil
}

func PromptForPassword(pass *string){
	for {
		if *pass == "" {
			fmt.Print("Enter Password: ")
			fmt.Scanln(pass)
		}

		err := validatePassword(*pass)
		if err == nil {
			break
		}

		// Print error and reset password so we
		// get a prompt on the next iteration
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		*pass = ""
	}

	// TODO: add a "retype password" prompt later on 
}

func validatePassword(pass string) error {
	if pass == "" {
		return fmt.Errorf("password cannot be empty")
	}

	passLength := len([]rune(pass))

	if passLength < 8 {
		return fmt.Errorf("Password needs to be longer")
	}

	if passLength > 128 {
		return fmt.Errorf("password needs to be shorter")
	}

	for _, char := range pass {
		if unicode.IsControl(char) {
			return fmt.Errorf("password contains invalid control characters")
		}
	}

	return nil
}