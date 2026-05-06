# Grapple

Grapple is a set of tools designed to transform Azure DevOps and GitHub Code Search API results into a clean, **grep-like syntax**. 

This format is optimized for AI agents and CLI workflows, providing precise `:offset:length` metadata to facilitate surgical code analysis with minimal token usage.

## Tools

- `grapple-az`: For Azure DevOps Code Search.
- `grapple-gh-api`: For GitHub Code Search via raw API responses.
- `grapple-gh-cli`: For GitHub Code Search via the `gh` CLI output.

## Installation

Ensure you have [Go](https://go.dev/) installed, then build the executables:

```powershell
go build -o grapple-az.exe ./cmd/grapple-az
go build -o grapple-gh-api.exe ./cmd/grapple-gh-api
go build -o grapple-gh-cli.exe ./cmd/grapple-gh-cli
```

## Output Format

- **Azure DevOps**: `Project/Repository/Path:Offset:Length`
- **GitHub**: `Repository/Path:Offset:Length`

---

## Azure DevOps (`grapple-az`)

Pipe the JSON output from the Azure DevOps REST API.

### Example using `az rest`

```powershell
az rest --method post `
  --uri "https://almsearch.dev.azure.com/{org}/{project}/_apis/search/codesearchresults?api-version=7.1" `
  --body '{"searchText": "MyFunction", "$top": 10}' | .\grapple-az.exe
```

---

## GitHub API (`grapple-gh-api`)

Optimized for raw API responses. This is the best method for AI context as it ensures exact character offsets are included via the `text-match` header.

```powershell
gh api -H "Accept: application/vnd.github.v3.text-match+json" `
  "/search/code?q=OmniboxEd+repo:chromium/chromium&per_page=5" | .\grapple-gh-api.exe
```

---

## GitHub CLI (`grapple-gh-cli`)

Optimized for the `gh search code` command output.

```powershell
gh search code "OmniboxEd" --repo "chromium/chromium" --limit 5 --json path,repository,textMatches | .\grapple-gh-cli.exe
```

---

## Why Grapple?

Standard API outputs are often verbose and contain heavy metadata that isn't needed for search results. Grapple strips this away, leaving only the essential location data. This allows an AI agent to:
1.  Search across an entire organization quickly.
2.  Receive a compact list of "hits".
3.  Precisely target only the relevant lines/blocks for subsequent "read" operations.

### Technical Note: Character Offsets

- **Azure DevOps (`grapple-az`)**: Character offsets are **absolute** to the start of the file.
- **GitHub (`grapple-gh-api` / `grapple-gh-cli`)**: Character offsets are **relative** to the start of the match fragment (snippet). GitHub's search API does not provide absolute file offsets in its search results. AI agents must be aware that these offsets refer to the provided snippet context.
