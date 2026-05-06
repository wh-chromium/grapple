# Azure DevOps Code Search API Reference

## Code Search Results - Fetch Code Search Results

**POST** `https://almsearch.dev.azure.com/{organization}/{project}/_apis/search/codesearchresults?api-version=7.1`

### URI Parameters
- `organization` (path, required): The name of the Azure DevOps organization.
- `project` (path): Project ID or project name.
- `api-version` (query, required): Set to `7.1`.

### Request Body
```json
{
  "searchText": "string",
  "$skip": 0,
  "$top": 10,
  "filters": {
    "Project": ["Name"],
    "Repository": ["Name"],
    "Path": ["/"],
    "Branch": ["master"]
  },
  "$orderBy": [
    {
      "field": "filename",
      "sortOrder": "ASC"
    }
  ],
  "includeFacets": false,
  "includeSnippet": false
}
```

### Response (200 OK)
Returns a `CodeSearchResponse` containing a list of `CodeResult` items with hit offsets.
