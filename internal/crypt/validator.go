package crypt

import (
	"errors"
	"fmt"
	"os"
)

// ValidateInputFile checks whether path exists, is a regular file, and is readable.
// It returns an error describing why the file is invalid, or nil if valid.
func ValidateInputFile(filename string) error {
	if filename == "" {
		return errors.New("input file path cannot be empty")
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

// ValidateOutputFile checks whether file exists (to avoid accidental overwriting),
// file name and path are proper and that the program has proper write permissions.
// It returns an error describing why the file is invalid, or nil if valid.
func ValidateOutputFile(filename string) error {
	panic("unimplemented")
}


func ValidatePassword(pass string) error {
	panic("unimplemented")
}