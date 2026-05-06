# Token Reduction Comparison

Copyright 2026 The Chromium Authors. All rights reserved.
Use of this source code is governed by a BSD-style license that can be
found in the LICENSE file.

*Note: The following examples use materials from the Chromium project for demonstration purposes only.*

This document compares the raw JSON output from GitHub search APIs against the "compacted" output produced by Grapple. The goal is to show the significant reduction in character count (and thus token usage) when providing context to AI agents.

## Search Case: "Hash" in `chromium/chromium`
**Query:** `Hash repo:chromium/chromium` (Limited to top 5 results)

### 1. GitHub API (`gh api`) vs `grapple-gh-api`

| Format | Character Count | Reduction |
| :--- | :--- | :--- |
| **Raw JSON (API)** | 52,464 | - |
| **Grapple (Compacted)** | 1,218 | **97.7%** |

#### Raw API Snippet (First ~500 chars)
```json
{
  "total_count": 142058,
  "incomplete_results": false,
  "items": [
    {
      "name": "hash.h",
      "path": "base/hash/hash.h",
      "sha": "...",
      "url": "...",
      "git_url": "...",
      "html_url": "...",
      "repository": {
        "id": 120360765,
        "node_id": "...",
        "name": "chromium",
        "full_name": "chromium/chromium",
        ...
```

#### Compacted Grapple Output
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

| Format | Character Count | Reduction |
| :--- | :--- | :--- |
| **Raw JSON (CLI)** | 6,082 | - |
| **Grapple (Compacted)** | 1,218 | **80.0%** |

#### Raw CLI Snippet
```json
[
  {
    "path": "base/hash/hash.h",
    "repository": {
      "nameWithOwner": "chromium/chromium"
    },
    "textMatches": [
      {
        "fragment": "...",
        "matches": [ ... ]
      }
    ]
  },
  ...
]
```

#### Compacted Grapple Output
*(Identical to the API compacted output above)*

---

## Conclusion

By stripping away verbose metadata and repetitive repository objects, Grapple allows you to fit **~43x more search results** into the same AI context window compared to raw GitHub API responses.
