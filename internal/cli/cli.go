package cli

import (
	"fmt"
	"io"
	"os"
)

// GetInput returns an io.ReadCloser for the input source.
// It returns os.Stdin if no argument is provided, or opens the file if one is.
func GetInput(args []string) (io.ReadCloser, error) {
	if len(args) > 1 {
		file, err := os.Open(args[1])
		if err != nil {
			return nil, fmt.Errorf("error opening file: %w", err)
		}
		return file, nil
	}
	return io.NopCloser(os.Stdin), nil
}

// ExitOnError prints the error and exits with code 1.
func ExitOnError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
