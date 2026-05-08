package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/wh-chromium/grapple/internal/cli"
)

type GHItem struct {
	Path       string `json:"path"`
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
	TextMatches []struct {
		Matches []struct {
			Indices []int `json:"indices"`
		} `json:"matches"`
	} `json:"text_matches"`
}

type GHSearchResponse struct {
	Items []GHItem `json:"items"`
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

	foundItems := false
	for decoder.More() {
		t, err := decoder.Token()
		if err != nil {
			return fmt.Errorf("failed to read JSON token: %w", err)
		}
		
		key, ok := t.(string)
		if !ok {
			continue // Skip non-string keys if they somehow occur
		}

		if key == "items" {
			foundItems = true
			
			// Expect the start of the array
			t, err := decoder.Token()
			if err != nil {
				return fmt.Errorf("failed to read JSON token after 'items' key: %w", err)
			}
			if delim, ok := t.(json.Delim); !ok || delim != '[' {
				return fmt.Errorf("expected JSON array start after 'items' key, got: %T %v", t, t)
			}

			for decoder.More() {
				var item GHItem
				if err := decoder.Decode(&item); err != nil {
					return fmt.Errorf("failed to decode GitHub API JSON item: %w", err)
				}

				repo := item.Repository.FullName
				path := item.Path

				if len(item.TextMatches) > 0 {
					for _, tm := range item.TextMatches {
						for _, m := range tm.Matches {
							if len(m.Indices) == 2 {
								fmt.Fprintf(w, "%s/%s:%d:%d\n", repo, path, m.Indices[0], m.Indices[1]-m.Indices[0])
							}
						}
					}
				} else if repo != "" && path != "" {
					fmt.Fprintf(w, "%s/%s:0:0\n", repo, path)
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
			// If it's an object or array, decoder.Decode with a discard target or json.RawMessage works,
			// but an empty interface{} is simpler to just consume and discard the value.
			var discard interface{}
			if err := decoder.Decode(&discard); err != nil {
				return fmt.Errorf("failed to skip JSON value for key %q: %w", key, err)
			}
		}
	}

	if !foundItems {
		return fmt.Errorf("could not find 'items' key in the JSON object")
	}

	return nil
}
