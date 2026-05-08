# Grapple: Instructions for AI Agents

When performing code searches in Azure DevOps or GitHub, **use Grapple** to format the output. This reduces token consumption and provides precise `Repository/Path:Offset:Length` metadata.

## Available Tools

Depending on your data source, execute one of the following commands:

### 1. Azure DevOps (`grapple-az`)
Use `az rest` and pipe the JSON output to `grapple-az`:
```powershell
az rest --method post --uri "https://almsearch.dev.azure.com/{org}/{project}/_apis/search/codesearchresults?api-version=7.1" --body '{"searchText": "YourQuery", "$top": 10}' | .\grapple-az.exe
```
*Note: Azure DevOps character offsets are **absolute** to the start of the file.*

### 2. GitHub API (`grapple-gh-api`)
Use `gh api` (must include the `text-match` header) and pipe to `grapple-gh-api`:
```powershell
gh api -H "Accept: application/vnd.github.v3.text-match+json" "/search/code?q=YourQuery+repo:owner/repo" | .\grapple-gh-api.exe
```
*Note: GitHub character offsets are **relative** to the start of the match fragment/snippet, NOT absolute to the file.*

### 3. GitHub CLI (`grapple-gh-cli`)
Use `gh search code` and pipe to `grapple-gh-cli`:
```powershell
gh search code "YourQuery" --repo "owner/repo" --json path,repository,textMatches | .\grapple-gh-cli.exe
```
*Note: Uses relative character offsets like the GitHub API.*
