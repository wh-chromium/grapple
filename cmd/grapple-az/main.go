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
	var response CodeSearchResponse
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&response); err != nil {
		return fmt.Errorf("failed to decode Azure DevOps JSON: %w", err)
	}

	for _, result := range response.Results {
		for _, match := range result.Matches.Content {
			// Ensure path has a leading slash for consistency
			path := result.Path
			if path != "" && path[0] != '/' {
				path = "/" + path
			}
			fmt.Fprintf(w, "%s/%s%s:%d:%d\n", result.Project.Name, result.Repository.Name, path, match.CharOffset, match.Length)
		}
	}
	return nil
}
