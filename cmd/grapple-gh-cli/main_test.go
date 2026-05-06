package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestTransformGH_CLI(t *testing.T) {
	sampleJSON := `[
		{
			"path": "internal/logic.go",
			"repository": {"nameWithOwner": "org/repo-cli"},
			"textMatches": [
				{
					"matches": [
						{"indices": [100, 110]}
					]
				},
				{
					"matches": [
						{"indices": [200, 205]}
					]
				}
			]
		},
		{
			"path": "LICENSE",
			"repository": {"nameWithOwner": "org/repo-cli"},
			"textMatches": null
		}
	]`

	input := strings.NewReader(sampleJSON)
	output := &bytes.Buffer{}

	if err := Run(input, output); err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	expected := "org/repo-cli/internal/logic.go:100:10\norg/repo-cli/internal/logic.go:200:5\norg/repo-cli/LICENSE:0:0\n"
	if output.String() != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, output.String())
	}
}
