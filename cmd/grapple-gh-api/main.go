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
	var response GHSearchResponse
	if err := json.NewDecoder(r).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode GitHub API JSON object: %w", err)
	}

	for _, item := range response.Items {
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
	return nil
}
