# Format Comparison

Copyright 2026 The Chromium Authors. All rights reserved.
Use of this source code is governed by a BSD-style license that can be
found in the LICENSE file.

*Note: The following examples use materials from the Chromium project for demonstration purposes only.*

This document compares the raw JSON output from GitHub search APIs against the "compacted" output produced by Grapple. The goal is to illustrate how Grapple simplifies verbose metadata into a clean, grep-like format for use in CLI workflows and AI contexts.

## Search Case: "Hash" in `chromium/chromium`
**Query:** `Hash repo:chromium/chromium` (Limited to top 5 results)

### 1. GitHub API (`gh api`) vs `grapple-gh-api`

The standard GitHub API response is highly verbose, repeating metadata for the repository, owner, and file for every result.

#### Raw API Snippet (Example structure)
```json
{
  "total_count": 142058,
  "incomplete_results": false,
  "items": [
    {
      "name": "hash.cc",
      "path": "crypto/hash.cc",
      "repository": {
        "full_name": "chromium/chromium",
        "owner": {
          "login": "chromium",
          "url": "https://api.github.com/users/chromium",
          ...
        },
        "html_url": "https://github.com/chromium/chromium",
        "description": "The official GitHub mirror of the Chromium source",
        ...
      },
      "text_matches": [
        {
          "fragment": "...",
          "matches": [ { "indices": [6, 10], "text": "Hash" } ]
        }
      ]
    },
    ...
  ]
}
```

#### Compacted Grapple Output
Grapple strips away the repetitive metadata, leaving only the essential location and match information.

```text
chromium/chromium/crypto/hash.cc:6:4
chromium/chromium/crypto/hash.cc:11:4
chromium/chromium/crypto/hash.cc:81:4
chromium/chromium/ui/gfx/render_text_harfbuzz.h:86:4
chromium/chromium/ui/gfx/render_text_harfbuzz.cc:95:4
chromium/chromium/ui/gfx/render_text_harfbuzz.cc:183:4
chromium/chromium/crypto/obsolete/md5.cc:47:4
chromium/chromium/crypto/obsolete/md5.cc:86:4
chromium/chromium/crypto/obsolete/md5.cc:47:4
chromium/chromium/gpu/ipc/service/context_url.h:86:4
chromium/chromium/gpu/ipc/service/context_url.h:112:4
```

---

### 2. GitHub CLI (`gh search`) vs `grapple-gh-cli`

While the GitHub CLI provides a more streamlined JSON format than the raw API, it still contains structural overhead that is unnecessary for simple search reporting.

#### Raw CLI Snippet (Example structure)
```json
[
  {
    "path": "crypto/hash.cc",
    "repository": {
      "nameWithOwner": "chromium/chromium"
    },
    "textMatches": [
      {
        "fragment": "...",
        "matches": [ { "indices": [6, 10] } ]
      }
    ]
  },
  ...
]
```

#### Compacted Grapple Output
Grapple transforms this into the same standardized format used across all supported search engines.

```text
chromium/chromium/crypto/hash.cc:6:4
...
```

---

## Conclusion

Grapple facilitates more efficient context usage by stripping away repetitive JSON metadata and verbose object structures. By transforming these responses into a standardized `:offset:length` format, Grapple allows for more surgical code analysis and easier integration with traditional Unix-style tools and AI-driven workflows.
