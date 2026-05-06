package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestTransformGH_API(t *testing.T) {
	sampleJSON := `{
		"items": [
			{
				"path": "src/main.go",
				"repository": {"full_name": "org/repo1"},
				"text_matches": [
					{
						"matches": [
							{"indices": [10, 15]},
							{"indices": [20, 25]}
						]
					}
				]
			},
			{
				"path": "docs/README.md",
				"repository": {"full_name": "org/repo2"},
				"text_matches": []
			}
		]
	}`

	input := strings.NewReader(sampleJSON)
	output := &bytes.Buffer{}

	if err := Run(input, output); err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	expected := "org/repo1/src/main.go:10:5\norg/repo1/src/main.go:20:5\norg/repo2/docs/README.md:0:0\n"
	if output.String() != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, output.String())
	}
}
