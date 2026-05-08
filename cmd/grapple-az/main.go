package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/wh-chromium/grapple/internal/cli"
)

type Match struct {
	CharOffset int `json:"charOffset"`
	Length     int `json:"length"`
}

type Matches struct {
	Content []Match `json:"content"`
}

type Repository struct {
	Name string `json:"name"`
}

type Project struct {
	Name string `json:"name"`
}

type Result struct {
	Path       string     `json:"path"`
	Matches    Matches    `json:"matches"`
	Repository Repository `json:"repository"`
	Project    Project    `json:"project"`
}

type CodeSearchResponse struct {
	Results []Result `json:"results"`
}

func main() {
	input, err := cli.GetInput(os.Args)
	cli.ExitOnError(err)
	defer input.Close()

	err = Run(input, os.Stdout)
	cli.ExitOnError(err)
}

func Run(r io.Reader, w io.Writer) error {
	decoder := json.NewDecoder(r)

	// Expect the start of the outer object
	t, err := decoder.Token()
	if err != nil {
		return fmt.Errorf("failed to read JSON token: %w", err)
	}
	if delim, ok := t.(json.Delim); !ok || delim != '{' {
		return fmt.Errorf("expected JSON object start, got: %T %v", t, t)
	}

	foundResults := false
	for decoder.More() {
		t, err := decoder.Token()
		if err != nil {
			return fmt.Errorf("failed to read JSON token: %w", err)
		}
		
		key, ok := t.(string)
		if !ok {
			continue // Skip non-string keys
		}

		if key == "results" {
			foundResults = true
			
			// Expect the start of the array
			t, err := decoder.Token()
			if err != nil {
				return fmt.Errorf("failed to read JSON token after 'results' key: %w", err)
			}
			if delim, ok := t.(json.Delim); !ok || delim != '[' {
				return fmt.Errorf("expected JSON array start after 'results' key, got: %T %v", t, t)
			}

			for decoder.More() {
				var result Result
				if err := decoder.Decode(&result); err != nil {
					return fmt.Errorf("failed to decode Azure DevOps JSON item: %w", err)
				}

				for _, match := range result.Matches.Content {
					// Ensure path has a leading slash for consistency
					path := result.Path
					if path != "" && path[0] != '/' {
						path = "/" + path
					}
					fmt.Fprintf(w, "%s/%s%s:%d:%d\n", result.Project.Name, result.Repository.Name, path, match.CharOffset, match.Length)
				}
			}
			
			// Expect the end of the array
			t, err = decoder.Token()
			if err != nil {
				return fmt.Errorf("failed to read JSON token: %w", err)
			}
			if delim, ok := t.(json.Delim); !ok || delim != ']' {
				return fmt.Errorf("expected JSON array end, got: %T %v", t, t)
			}
		} else {
			// Skip the value for this key
			var discard interface{}
			if err := decoder.Decode(&discard); err != nil {
				return fmt.Errorf("failed to skip JSON value for key %q: %w", key, err)
			}
		}
	}

	if !foundResults {
		// Just output nothing if "results" isn't present, Azure Devops might return errors instead of items
		return nil
	}

	return nil
}
