package cli

import (
	"os"
	"testing"
)

func TestGetInput_Stdin(t *testing.T) {
	// Should default to stdin if < 2 args
	args := []string{"binary_name"}
	reader, err := GetInput(args)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if reader == nil {
		t.Errorf("Expected a valid reader, got nil")
	}
}

func TestGetInput_File(t *testing.T) {
	// Create a temp file
	tmpfile, err := os.CreateTemp("", "test_input.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	args := []string{"binary_name", tmpfile.Name()}
	reader, err := GetInput(args)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer reader.Close()
	
	if reader == os.Stdin {
		t.Errorf("Expected reader to be a file, not os.Stdin")
	}
}

func TestGetInput_File_Error(t *testing.T) {
	args := []string{"binary_name", "non_existent_file_12345.txt"}
	_, err := GetInput(args)
	if err == nil {
		t.Fatalf("Expected error for non-existent file, got nil")
	}
}
