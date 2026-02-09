# Structured Tasks

[![Build Status][build-status-svg]][build-status-url]
[![Lint Status][lint-status-svg]][lint-status-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![Visualization][viz-svg]][viz-url]
[![License][license-svg]][license-url]

Structured Tasks provides a machine-readable JSON intermediate representation (IR) for project task lists, with deterministic Markdown generation. It is modeled after [Structured Changelog](https://github.com/grokify/structured-changelog).

## Features

- **JSON IR** - Machine-readable task list format with rich metadata
- **Deterministic output** - Same JSON always produces identical Markdown
- **Two-dimensional categorization** - Area (project component) + Type (change type)
- **Multiple grouping strategies** - Group by area, type, phase, status, quarter, or priority
- **Phased task lists** - Support for large projects with phases and area sub-sections
- **Rich content support** - Code blocks, tables, diagrams, lists, and blockquotes
- **Type validation** - Integrates with [structured-changelog](https://github.com/grokify/structured-changelog) for type consistency
- **Dependency tracking** - Item dependencies with graph generation
- **Validation** - Schema validation with detailed error messages
- **Statistics** - Track progress and completion rates

## Installation

### Homebrew (macOS/Linux)

```bash
brew tap grokify/tap
brew install structured-tasks
```

This installs the `stasks` CLI (also available as `structured-tasks`).

### Go Install

```bash
go install github.com/grokify/structured-tasks/cmd/stasks@latest
```

### Go Library

```bash
go get github.com/grokify/structured-tasks
```

## Quick Start

### Create a TASKS.json

```json
{
  "irVersion": "1.0",
  "project": "my-project",
  "areas": [
    {"id": "core", "name": "Core Features", "priority": 1}
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
      "description": "Add configurable rate limits",
      "status": "planned",
      "targetQuarter": "Q2 2026",
      "area": "core",
      "type": "Added",
      "priority": "medium",
      "dependsOn": ["feature-1"]
    }
  ]
}
```

### Generate Markdown

```bash
stasks generate -i TASKS.json -o TASKS.md
```

### Validate

```bash
stasks validate TASKS.json
```

### Show Statistics

```bash
stasks stats TASKS.json
```

### Generate Dependency Graph

```bash
stasks deps TASKS.json --format mermaid
```

## CLI Commands

### validate

Validate a TASKS.json file against the schema.

```bash
stasks validate TASKS.json
```

### generate

Generate TASKS.md from TASKS.json.

```bash
stasks generate -i TASKS.json -o TASKS.md
```

Options:

| Flag | Default | Description |
|------|---------|-------------|
| `-i, --input` | TASKS.json | Input JSON file |
| `-o, --output` | stdout | Output Markdown file |
| `--group-by` | area | Grouping: area, type, phase, status, quarter, priority |
| `--checkboxes` | true | Use [x]/[ ] checkbox syntax |
| `--emoji` | true | Include emoji status indicators |
| `--legend` | false | Show legend table |
| `--no-intro` | false | Omit introductory paragraph (intro shown by default) |
| `--toc` | false | Show table of contents with progress counts |
| `--toc-depth` | 1 | TOC depth: 1 = sections only, 2 = sections + items |
| `--overview` | true | Show summary table with all items |
| `--area-subheadings` | false | Show area sub-sections within phases |
| `--numbered` | false | Number items |
| `--no-rules` | false | Omit horizontal rules between sections |

### stats

Show task list statistics.

```bash
stasks stats TASKS.json
```

Output:

```
Task List: my-project
Total items: 10

By Status:
  ✅ Completed: 4 (40%)
  🚧 In Progress: 2 (20%)
  📋 Planned: 3 (30%)
  💡 Under Consideration: 1 (10%)

By Priority:
  High Priority: 3 (30%)
  Medium Priority: 5 (50%)
  Low Priority: 2 (20%)

By Area:
  Core Features: 6
  Improvements: 4

By Type:
  Added: 5
  Changed: 3
  Fixed: 2

Progress: 40% complete
```

### deps

Generate dependency graph in Mermaid or DOT format.

```bash
stasks deps TASKS.json --format mermaid
```

## JSON IR Schema

### Top-Level Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `irVersion` | string | Yes | Schema version ("1.0") |
| `project` | string | Yes | Project name |
| `repository` | string | No | Repository URL |
| `generatedAt` | datetime | No | Generation timestamp |
| `legend` | object | No | Custom status legend |
| `areas` | array | No | Project areas/components |
| `phases` | array | No | Development phases |
| `items` | array | No | Roadmap items |
| `sections` | array | No | Freeform content sections |
| `versionHistory` | array | No | Version milestones |
| `dependencies` | object | No | External/internal dependencies |

### Item Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Yes | Unique identifier |
| `title` | string | Yes | Item title |
| `description` | string | No | Item description |
| `status` | enum | Yes | completed, inProgress, planned, future |
| `version` | string | No | Version where completed |
| `completedDate` | date | No | Completion date |
| `targetQuarter` | string | No | Target quarter (e.g., "Q2 2026") |
| `targetVersion` | string | No | Target version |
| `area` | string | No | Area ID (project component) |
| `type` | string | No | Change type (aligns with structured-changelog) |
| `phase` | string | No | Phase ID |
| `priority` | enum | No | critical, high, medium, low |
| `order` | int | No | Explicit sort order within groups |
| `dependsOn` | array | No | IDs of dependencies |
| `tasks` | array | No | Sub-tasks with completion status |
| `content` | array | No | Rich content blocks |

### Two-Dimensional Categorization

Items can be categorized along two orthogonal dimensions:

| Dimension | Field | Purpose | Example Values |
|-----------|-------|---------|----------------|
| **Area** | `area` | WHERE in the project | "core", "api", "ui", "backend" |
| **Type** | `type` | WHAT kind of change | "Added", "Changed", "Fixed", "Security" |

- **Area** is user-defined and groups items by project component
- **Type** aligns with [structured-changelog](https://github.com/grokify/structured-changelog) change types and is validated against the registry

This allows grouping by area for task lists (`--group-by area`) while preserving type information for changelog integration when items are completed.

### Content Block Types

| Type | Fields | Description |
|------|--------|-------------|
| `text` | value | Markdown text |
| `code` | value, language | Code block |
| `diagram` | value, format | ASCII or Mermaid diagram |
| `table` | headers, rows | Markdown table |
| `list` | items | Bullet list |
| `blockquote` | value | Blockquote/callout (renders with `>` prefix) |

## Phased Task Lists (Large Projects)

For large projects with multiple development phases (like [omnistorage](https://github.com/grokify/omnistorage)), use the combination of `phases` and `areas` to create hierarchical task lists.

### Structure

```
Phase 1: Foundation ✅
├── Core Package
│   ├── [x] interfaces.go
│   └── [x] options.go
├── Format Layer
│   └── [x] ndjson/writer.go
└── Backend Layer
    └── [x] file/backend.go

Phase 2: Extended Interfaces ✅
├── Core Interfaces
│   └── [x] extended.go
└── Utilities
    └── [x] copy.go
```

### JSON Structure

```json
{
  "irVersion": "1.0",
  "project": "my-large-project",
  "phases": [
    {"id": "phase-1", "name": "Phase 1: Foundation", "status": "completed", "order": 1},
    {"id": "phase-2", "name": "Phase 2: Extended Interfaces", "status": "completed", "order": 2},
    {"id": "phase-3", "name": "Phase 3: Cloud Storage", "status": "inProgress", "order": 3}
  ],
  "areas": [
    {"id": "core", "name": "Core Package", "priority": 1},
    {"id": "format", "name": "Format Layer", "priority": 2},
    {"id": "backend", "name": "Backend Layer", "priority": 3},
    {"id": "utils", "name": "Utilities", "priority": 4}
  ],
  "items": [
    {
      "id": "interfaces",
      "title": "`interfaces.go` - Backend, RecordWriter, RecordReader interfaces",
      "status": "completed",
      "phase": "phase-1",
      "area": "core",
      "type": "Added"
    },
    {
      "id": "ndjson-writer",
      "title": "`format/ndjson/writer.go` - NDJSON RecordWriter",
      "status": "completed",
      "phase": "phase-1",
      "area": "format",
      "type": "Added"
    }
  ],
  "sections": [
    {
      "id": "design-philosophy",
      "title": "Design Philosophy",
      "order": 0,
      "content": [
        {"type": "text", "value": "This project uses **interface composition** to support both simple and advanced use cases."},
        {"type": "blockquote", "value": "**Note:** Cloud provider backends with large SDK dependencies are in separate repos."}
      ]
    }
  ]
}
```

### Generate with Area Sub-headings

When using phases, enable area sub-headings to show logical groupings within each phase:

```bash
stasks generate -i TASKS.json -o TASKS.md \
  --group-by phase \
  --area-subheadings \
  --toc
```

This produces output like:

```markdown
## Phase 1: Foundation ✅

Core interfaces and essential implementations.

### Core Package

- [x] `interfaces.go` - Backend, RecordWriter, RecordReader interfaces
- [x] `options.go` - WriterOption, ReaderOption, configs

### Format Layer

- [x] `format/ndjson/writer.go` - NDJSON RecordWriter
- [x] `format/ndjson/reader.go` - NDJSON RecordReader

---

## Phase 2: Extended Interfaces ✅
...
```

### Key Points

1. **Phases** define sequential development stages with their own status
2. **Areas** define logical groupings within phases (components, layers)
3. Items have both `phase` and `area` fields for two-dimensional organization
4. Use `--group-by phase --area-subheadings` for hierarchical output
5. Items render as compact task lists under area sub-headings
6. Use `blockquote` content type for notes and callouts

## Library Usage

```go
package main

import (
    "fmt"
    "github.com/grokify/structured-tasks/tasks"
    "github.com/grokify/structured-tasks/renderer"
)

func main() {
    // Parse task list
    tl, err := tasks.ParseFile("TASKS.json")
    if err != nil {
        panic(err)
    }

    // Validate
    result := tasks.Validate(tl)
    if !result.Valid {
        for _, e := range result.Errors {
            fmt.Printf("Error: %s: %s\n", e.Field, e.Message)
        }
        return
    }

    // Get statistics
    stats := tl.Stats()
    fmt.Printf("Progress: %.0f%% complete\n", stats.CompletedPercent())

    // Render to Markdown
    opts := renderer.DefaultOptions()
    opts.GroupBy = renderer.GroupByPriority
    opts.ShowTOC = true
    output := renderer.Render(tl, opts)
    fmt.Println(output)
}
```

## Related Projects

- [Structured Changelog](https://github.com/grokify/structured-changelog) - Machine-readable changelogs
- [Keep a Changelog](https://keepachangelog.com/) - Changelog format inspiration

## License

MIT License - see [LICENSE](LICENSE) for details.

 [build-status-svg]: https://github.com/grokify/structured-tasks/actions/workflows/ci.yaml/badge.svg?branch=main
 [build-status-url]: https://github.com/grokify/structured-tasks/actions/workflows/ci.yaml
 [lint-status-svg]: https://github.com/grokify/structured-tasks/actions/workflows/lint.yaml/badge.svg?branch=main
 [lint-status-url]: https://github.com/grokify/structured-tasks/actions/workflows/lint.yaml
 [coverage-svg]: https://img.shields.io/badge/coverage-94.9%25-brightgreen
 [coverage-url]: https://github.com/grokify/structured-tasks
 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/structured-tasks
 [goreport-url]: https://goreportcard.com/report/github.com/grokify/structured-tasks
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/grokify/structured-tasks
 [docs-godoc-url]: https://pkg.go.dev/github.com/grokify/structured-tasks
 [viz-svg]: https://img.shields.io/badge/visualizaton-Go-blue.svg
 [viz-url]: https://mango-dune-07a8b7110.1.azurestaticapps.net/?repo=grokify%2Fstructured-tasks
 [loc-svg]: https://tokei.rs/b1/github/grokify/structured-tasks
 [repo-url]: https://github.com/grokify/structured-tasks
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/grokify/structured-tasks/blob/master/LICENSE
