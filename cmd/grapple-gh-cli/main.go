package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/wh-chromium/grapple/internal/cli"
)

type GHCLIItem struct {
	Path       string `json:"path"`
	Repository struct {
		NameWithOwner string `json:"nameWithOwner"`
	} `json:"repository"`
	TextMatches []struct {
		Matches []struct {
			Indices []int `json:"indices"`
		} `json:"matches"`
	} `json:"textMatches"`
}

func main() {
	input, err := cli.GetInput(os.Args)
	cli.ExitOnError(err)
	defer input.Close()

	err = Run(input, os.Stdout)
	cli.ExitOnError(err)
}

func Run(r io.Reader, w io.Writer) error {
	var items []GHCLIItem
	if err := json.NewDecoder(r).Decode(&items); err != nil {
		return fmt.Errorf("failed to decode GitHub CLI JSON array: %w", err)
	}

	for _, item := range items {
		repo := item.Repository.NameWithOwner
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
	return nil
}
