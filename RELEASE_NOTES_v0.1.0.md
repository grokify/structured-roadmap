# Release Notes v0.1.0

**Release Date:** January 10, 2026

## Overview

This is the initial release of Structured Roadmap, a machine-readable JSON intermediate representation (IR) for project roadmaps with deterministic Markdown generation.

Structured Roadmap is modeled after [Structured Changelog](https://github.com/grokify/structured-changelog) and provides a canonical JSON format as the source of truth for your project roadmap, enabling consistent documentation across repositories.

## Installation

### CLI Tool

```bash
go install github.com/grokify/structured-roadmap/cmd/sroadmap@v0.1.0
```

### Library

```bash
go get github.com/grokify/structured-roadmap@v0.1.0
```

## Features

### JSON IR Schema (v1.0)

The JSON Intermediate Representation is the canonical source of truth. Key features:

- **Two-dimensional categorization**: Items grouped by `area` (project component) and `type` (change type)
- **Multiple grouping strategies**: area, type, phase, status, quarter, priority
- **Phased roadmaps**: Support for large projects with development phases
- **Rich content blocks**: text, code, diagram, table, list, blockquote
- **Dependency tracking**: Item-to-item dependencies with graph visualization
- **Progress tracking**: Statistics by status, area, type, priority

Example `ROADMAP.json`:

```json
{
  "ir_version": "1.0",
  "project": "my-project",
  "repository": "https://github.com/org/my-project",
  "areas": [
    { "id": "core", "name": "Core Features", "priority": 1 }
  ],
  "items": [
    {
      "id": "feature-1",
      "title": "User Authentication",
      "description": "Add OAuth2 login support",
      "status": "completed",
      "version": "1.0.0",
      "area": "core",
      "type": "Added",
      "priority": "high"
    },
    {
      "id": "feature-2",
      "title": "API Rate Limiting",
      "status": "planned",
      "target_quarter": "Q2 2026",
      "area": "core",
      "depends_on": ["feature-1"]
    }
  ]
}
```

### Status Values

| Status | Emoji | Description |
|--------|-------|-------------|
| `completed` | ✅ | Done, optionally with `version` and `completed_date` |
| `in_progress` | 🚧 | Currently being worked on |
| `planned` | 📋 | Scheduled for future work |
| `future` | 💡 | Under consideration, not yet scheduled |

### Priority Levels

| Priority | Description |
|----------|-------------|
| `critical` | Must have, blocking |
| `high` | Important, near-term |
| `medium` | Standard priority |
| `low` | Nice to have |

### Content Blocks

Rich content support for detailed item descriptions:

| Type | Fields | Description |
|------|--------|-------------|
| `text` | `value` | Markdown text |
| `code` | `value`, `language` | Code block with syntax highlighting |
| `diagram` | `value`, `format` | ASCII or Mermaid diagram |
| `table` | `headers`, `rows` | Multi-column table |
| `list` | `items` | Bullet list |
| `blockquote` | `value` | Quoted text / callout |

### CLI Tool (`sroadmap`)

Cobra-based CLI with four subcommands:

**Validate a roadmap:**

```bash
sroadmap validate ROADMAP.json
```

Output:

```
✓ ROADMAP.json is valid

Summary:
  Project: my-project
  Items: 10
  Completed: 4 (40%)
  In Progress: 2 (20%)
  Planned: 4 (40%)
```

**Generate Markdown:**

```bash
# To stdout
sroadmap generate -i ROADMAP.json

# To file
sroadmap generate -i ROADMAP.json -o ROADMAP.md

# With options
sroadmap generate -i ROADMAP.json -o ROADMAP.md \
  --group-by phase \
  --area-subheadings \
  --toc
```

**Show statistics:**

```bash
sroadmap stats ROADMAP.json
```

Output:

```
Roadmap: my-project
Total items: 10

By Status:
  ✅ Completed: 4 (40%)
  🚧 In Progress: 2 (20%)
  📋 Planned: 3 (30%)
  💡 Under Consideration: 1 (10%)

By Priority:
  Critical: 1 (10%)
  High: 3 (30%)
  Medium: 4 (40%)
  Low: 2 (20%)

Progress: 40% complete
```

**Generate dependency graph:**

```bash
# Mermaid format (default)
sroadmap deps ROADMAP.json --format mermaid

# Graphviz DOT format
sroadmap deps ROADMAP.json --format dot
```

### Generation Options

| Flag | Default | Description |
|------|---------|-------------|
| `--group-by` | `area` | Grouping: area, type, phase, status, quarter, priority |
| `--checkboxes` | `true` | Use `[x]`/`[ ]` checkbox syntax |
| `--emoji` | `true` | Include emoji status indicators |
| `--legend` | `false` | Show legend table |
| `--toc` | `false` | Show table of contents with progress counts |
| `--toc-depth` | `1` | TOC depth: 1 = sections only, 2 = sections + items |
| `--overview` | `false` | Show overview table with all items |
| `--area-subheadings` | `false` | Show area sub-sections within phases |
| `--numbered` | `false` | Number items sequentially |
| `--no-rules` | `false` | Omit horizontal rules between sections |
| `--no-intro` | `false` | Omit introductory paragraph |

### Go Library

**Load and validate:**

```go
import "github.com/grokify/structured-roadmap/roadmap"

r, err := roadmap.ParseFile("ROADMAP.json")
if err != nil {
    log.Fatal(err)
}

result := roadmap.Validate(r)
if !result.Valid {
    for _, e := range result.Errors {
        log.Printf("Error: %s: %s", e.Field, e.Message)
    }
}
```

**Get statistics:**

```go
stats := r.Stats()
fmt.Printf("Progress: %.0f%% complete\n", stats.CompletedPercent())
fmt.Printf("By Status: %v\n", stats.ByStatus)
fmt.Printf("By Area: %v\n", stats.ByArea)
```

**Render to Markdown:**

```go
import "github.com/grokify/structured-roadmap/renderer"

opts := renderer.DefaultOptions()
opts.GroupBy = renderer.GroupByPhase
opts.ShowTOC = true
opts.ShowAreaSubheadings = true

md := renderer.Render(r, opts)
```

**Grouping methods:**

```go
// Group items by different dimensions
byArea := r.ItemsByArea()
byType := r.ItemsByType()
byPhase := r.ItemsByPhase()
byStatus := r.ItemsByStatus()
byQuarter := r.ItemsByQuarter()
byPriority := r.ItemsByPriority()
```

### Phased Roadmaps

For large projects with multiple development phases:

```json
{
  "phases": [
    { "id": "phase-1", "name": "Phase 1: Foundation", "status": "completed", "order": 1 },
    { "id": "phase-2", "name": "Phase 2: Extensions", "status": "in_progress", "order": 2 }
  ],
  "areas": [
    { "id": "core", "name": "Core Package", "priority": 1 },
    { "id": "format", "name": "Format Layer", "priority": 2 }
  ],
  "items": [
    {
      "id": "interfaces",
      "title": "`interfaces.go` - Core interfaces",
      "status": "completed",
      "phase": "phase-1",
      "area": "core"
    }
  ]
}
```

Generate with area sub-headings:

```bash
sroadmap generate -i ROADMAP.json --group-by phase --area-subheadings --toc
```

Output structure:

```markdown
## Phase 1: Foundation ✅

### Core Package

- [x] `interfaces.go` - Core interfaces
- [x] `options.go` - Configuration options

### Format Layer

- [x] `ndjson/writer.go` - NDJSON writer
```

### Type Validation

Item `type` field is validated against [Structured Changelog](https://github.com/grokify/structured-changelog) change types:

- Added, Changed, Deprecated, Removed, Fixed, Security
- Highlights, Breaking, Performance, Dependencies
- Documentation, Build, Infrastructure, Internal

This enables seamless integration between roadmap items and changelog entries when features are completed.

## Project Structure

```
structured-roadmap/
├── roadmap/           # IR types, parsing, validation
├── renderer/          # Markdown generation
├── cmd/sroadmap/      # CLI tool
├── schema/            # JSON Schema v1
└── examples/          # Example roadmaps
```

## Examples

Two example roadmaps are included:

- `examples/minimal.json` - Simple two-item roadmap
- `examples/full.json` - Complex example with phases, areas, rich content, dependencies

## Dependencies

- `github.com/grokify/structured-changelog` v0.5.0 - Change type validation
- `github.com/spf13/cobra` v1.10.2 - CLI framework

## Related Projects

- [Structured Changelog](https://github.com/grokify/structured-changelog) - Machine-readable changelogs

## What's Next

Planned for future releases:

- `sroadmap init` - Create empty ROADMAP.json
- `sroadmap add` - Add items interactively
- GoReleaser configuration for binary releases
- Homebrew tap distribution
- MkDocs documentation site
