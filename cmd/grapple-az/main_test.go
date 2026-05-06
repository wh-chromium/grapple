package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestTransformAZ_VariousMatches(t *testing.T) {
	sampleJSON := `{
		"results": [
			{
				"path": "path/no-slash.cs",
				"matches": {
					"content": [
						{"charOffset": 10, "length": 5}
					]
				},
				"repository": {"name": "RepoA"},
				"project": {"name": "ProjA"}
			},
			{
				"path": "/already/slashed.go",
				"matches": {
					"content": [
						{"charOffset": 20, "length": 10},
						{"charOffset": 40, "length": 2}
					]
				},
				"repository": {"name": "RepoB"},
				"project": {"name": "ProjB"}
			}
		]
	}`

	input := strings.NewReader(sampleJSON)
	output := &bytes.Buffer{}

	if err := Run(input, output); err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	// Verify leading slash is added where missing
	expected := "ProjA/RepoA/path/no-slash.cs:10:5\n" +
		"ProjB/RepoB/already/slashed.go:20:10\n" +
		"ProjB/RepoB/already/slashed.go:40:2\n"
	if output.String() != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, output.String())
	}
}

func TestTransformAZ_Empty(t *testing.T) {
	sampleJSON := `{"results": []}`
	input := strings.NewReader(sampleJSON)
	output := &bytes.Buffer{}

	if err := Run(input, output); err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if output.Len() != 0 {
		t.Errorf("Expected empty output, got %q", output.String())
	}
}
