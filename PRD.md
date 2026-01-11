# Structured Roadmap - Product Requirements Document

## Overview

Structured Roadmap provides a machine-readable JSON intermediate representation (IR) for project roadmaps, with deterministic Markdown generation following common roadmap conventions. It is modeled after [Structured Changelog](https://github.com/grokify/structured-changelog).

## Problem Statement

Project roadmaps across repositories use inconsistent formats, making it difficult to:

- Maintain consistency across multiple projects
- Track status programmatically
- Generate standardized documentation
- Visualize dependencies and progress

## Research

The following repositories were reviewed to understand existing roadmap patterns:

| Repository | File |
|------------|------|
| `github.com/agentplexus/agentkit` | `ROADMAP.md` |
| `github.com/agentplexus/omnillm` | `ROADMAP.md` |
| `github.com/agentplexus/omnivault` | `ROADMAP.md` |
| `github.com/agentplexus/omnivoice` | `ROADMAP.md` |
| `github.com/agentplexus/stats-agent-team` | `ROADMAP.md` |
| `github.com/grokify/omnistorage` | `ROADMAP.md` |
| `github.com/grokify/traffic2openapi` | `ROADMAP.md` |

## Features Inventory

Features identified across all reviewed roadmap files:

| Feature | Files Using It |
|---------|----------------|
| **Status Indicators** | |
| ✅ completed | omnillm, stats-agent, agentkit, omnivault, traffic2openapi, omnistorage |
| 🚧 in progress | traffic2openapi |
| 📋 planned | traffic2openapi |
| 💡 under consideration | traffic2openapi |
| ⬚ TODO open | stats-agent |
| `[x]` / `[ ]` checkboxes | agentkit, stats-agent, omnivault, traffic2openapi, omnistorage |
| ~~strikethrough~~ completed | omnistorage (README) |
| Status column in tables | stats-agent, omnivoice, traffic2openapi |
| **Timing** | |
| Calendar quarters (Q1 2026) | stats-agent |
| Version targets (v0.3.0) | agentkit, omnillm, traffic2openapi, omnistorage |
| Version-based sections (v0.2.0) | traffic2openapi |
| **Priority** | |
| Numbered items (1-11) | omnillm |
| Priority column (P0, P1) | stats-agent |
| "Priority 1, 2, 3" ordering | omnivoice |
| Tiers (High/Medium/Nice) | omnillm |
| **Categories** | |
| Phases with status (Phase 1 ✅) | omnivoice, omnivault, omnistorage |
| Status sections (Completed/Planned) | agentkit, omnivault |
| Feature areas / subsections | agentkit, stats-agent, traffic2openapi |
| **Rich Content** | |
| Code blocks (go, yaml, bash) | omnillm, stats-agent, omnivoice, omnivault, traffic2openapi, omnistorage |
| ASCII diagrams | stats-agent, omnivoice, omnivault |
| Tables (multi-column) | all |
| Emoji prefixes (✨📊🔮🌐) | stats-agent |
| File path references | omnivault, omnistorage |
| Before/After tables | traffic2openapi |
| API endpoint tables | omnivault |
| Platform-specific tables | omnivault |
| **Dependencies** | |
| External SDK dependencies | omnivoice |
| Internal dependencies | omnivoice |
| `depends_on` relationships | omnivoice |
| **Special Sections** | |
| Legend table | traffic2openapi |
| Design Philosophy | omnistorage |
| Security Model | omnivault |
| Data Format / Structure | omnivault |
| Version History table | omnistorage |
| Contributing | agentkit, stats-agent, traffic2openapi, omnistorage |
| Resources/links | omnivoice |
| Future Considerations | stats-agent, traffic2openapi, omnistorage |
| Decision points with checkboxes | stats-agent |
| Horizontal rule separators | traffic2openapi, omnistorage |
| Blockquote notes (>) | omnistorage |
| **Nesting** | |
| Nested checkboxes with file paths | omnivault, omnistorage |
| Subsections within phases | omnivault, omnistorage |
| Grouped items by feature area | traffic2openapi |

## Proposed JSON IR (v1.0)

```json
{
  "ir_version": "1.0",
  "project": "omnillm",
  "repository": "https://github.com/agentplexus/omnillm",
  "generated_at": "2026-01-10T00:00:00Z",

  "legend": {
    "completed": {"emoji": "✅", "description": "Completed"},
    "in_progress": {"emoji": "🚧", "description": "In Progress"},
    "planned": {"emoji": "📋", "description": "Planned"},
    "future": {"emoji": "💡", "description": "Under Consideration"}
  },

  "categories": [
    {"id": "high_value", "name": "High Value", "priority": 1},
    {"id": "medium_value", "name": "Medium Value", "priority": 2},
    {"id": "improvement", "name": "Areas for Improvement", "priority": 3}
  ],

  "phases": [
    {
      "id": "phase1",
      "name": "Phase 1: Foundation",
      "status": "completed",
      "order": 1,
      "description": "Core interfaces and essential implementations"
    },
    {
      "id": "phase2",
      "name": "Phase 2: Extensions",
      "status": "in_progress",
      "order": 2
    }
  ],

  "items": [
    {
      "id": "timeout",
      "title": "Request Timeouts",
      "description": "Per-request timeout configuration",
      "status": "completed",
      "version": "0.10.0",
      "completed_date": "2026-01-04",
      "category": "high_value",
      "phase": "phase1",
      "priority": 2,
      "content": [
        {"type": "text", "value": "**Status:** Implemented via `ClientConfig.Timeout`."},
        {"type": "code", "language": "go", "value": "ClientConfig{\n    Timeout: 300 * time.Second,\n}"}
      ]
    },
    {
      "id": "fallback",
      "title": "Fallback Providers",
      "description": "Automatic failover when primary provider fails",
      "status": "planned",
      "target_quarter": "Q2 2026",
      "target_version": "0.12.0",
      "category": "high_value",
      "phase": "phase2",
      "priority": 4,
      "depends_on": ["timeout"],
      "tasks": [
        {"id": "fallback-api", "description": "Design fallback API", "completed": false},
        {"id": "health-checks", "description": "Implement provider health checks", "completed": false, "file_path": "internal/health/checker.go"}
      ]
    },
    {
      "id": "daemon-api",
      "title": "Daemon API",
      "description": "Background service for secure secret access",
      "status": "planned",
      "phase": "phase2",
      "content": [
        {
          "type": "table",
          "headers": ["Endpoint", "Method", "Description"],
          "rows": [
            ["/status", "GET", "Daemon status"],
            ["/secrets", "GET", "List all secrets"],
            ["/secret/:path", "GET", "Get secret value"],
            ["/lock", "POST", "Lock the vault"]
          ]
        }
      ]
    },
    {
      "id": "testing",
      "title": "Expand test coverage",
      "description": "Integration tests require API keys; add mock-based unit tests",
      "status": "planned",
      "category": "improvement",
      "subcategory": "Testing"
    }
  ],

  "sections": [
    {
      "id": "overview",
      "title": "Overview",
      "order": 1,
      "content": [
        {"type": "text", "value": "OmniLLM is a unified Go SDK for LLM providers."}
      ]
    },
    {
      "id": "design-philosophy",
      "title": "Design Philosophy",
      "order": 2,
      "content": [
        {"type": "text", "value": "Uses **interface composition** to support both simple and advanced use cases."},
        {"type": "code", "language": "go", "value": "// Simple apps use Backend\nfunc SaveData(backend Backend) error {\n    // ...\n}"}
      ]
    },
    {
      "id": "architecture",
      "title": "Recommended Architecture",
      "order": 3,
      "content": [
        {"type": "diagram", "format": "ascii", "value": "┌───────────┐\n│  Client   │\n└─────┬─────┘\n      │\n      ▼\n┌───────────┐\n│ Provider  │\n└───────────┘"}
      ]
    },
    {
      "id": "security",
      "title": "Security Model",
      "order": 4,
      "content": [
        {"type": "text", "value": "### Encryption\n- **Algorithm**: AES-256-GCM\n- **Key Derivation**: Argon2id"}
      ]
    },
    {
      "id": "contributing",
      "title": "Contributing",
      "order": 99,
      "content": [
        {"type": "text", "value": "We welcome contributions! Priority areas:"},
        {"type": "list", "items": ["New backends", "Tests", "Documentation", "Bug fixes"]}
      ]
    },
    {
      "id": "future",
      "title": "Future Considerations",
      "order": 100,
      "content": [
        {"type": "list", "items": ["GraphQL inference", "AsyncAPI support", "VS Code extension"]}
      ]
    }
  ],

  "version_history": [
    {"version": "0.5.0", "date": "2026-01-10", "status": "completed", "summary": "Sync engine complete"},
    {"version": "0.4.0", "date": null, "status": "planned", "summary": "Cloud providers"},
    {"version": "0.3.0", "date": "2025-12-15", "status": "completed", "summary": "S3 backend"}
  ],

  "dependencies": {
    "external": [
      {"name": "github.com/aws/aws-sdk-go-v2", "status": "available"},
      {"name": "github.com/deepgram/deepgram-go-sdk", "status": "available"},
      {"name": "Recall.ai", "status": "build_client", "note": "None (REST API)"}
    ],
    "internal": [
      {"package": "tts/elevenlabs", "depends_on": "go-elevenlabs"},
      {"package": "agent/custom", "depends_on": ["OmniLLM", "tts/", "stt/"]}
    ]
  }
}
```

## Schema Design Decisions

### 1. Items with Flexible Content

Items support an array of content blocks for rich formatting:

```json
"content": [
  {"type": "text", "value": "Markdown text here"},
  {"type": "code", "language": "go", "value": "code here"},
  {"type": "diagram", "format": "ascii", "value": "diagram here"},
  {"type": "table", "headers": ["Col1", "Col2"], "rows": [["a", "b"], ["c", "d"]]},
  {"type": "list", "items": ["Item 1", "Item 2", "Item 3"]}
]
```

### 2. Tasks as Sub-Items

Checkboxes become structured `tasks[]` with `completed` boolean and optional file paths:

```json
"tasks": [
  {"id": "task-1", "description": "Design API", "completed": true},
  {"id": "task-2", "description": "Write handler", "completed": false, "file_path": "internal/handler.go"}
]
```

### 3. Multiple Grouping Options

Items can be grouped by multiple dimensions:

| Field | Purpose | Example |
|-------|---------|---------|
| `category` | Primary grouping | `"high_value"`, `"improvement"` |
| `phase` | Development phase | `"phase1"`, `"phase2"` |
| `subcategory` | Secondary grouping | `"Testing"`, `"Documentation"` |
| `target_quarter` | Time-based planning | `"Q1 2026"`, `"Q2 2026"` |
| `target_version` | Version-based planning | `"v0.12.0"` |

### 4. Dependencies

Items can declare dependencies on other items:

```json
"depends_on": ["timeout", "retry"]
```

Project-level dependencies (external SDKs, internal packages) are also supported:

```json
"dependencies": {
  "external": [
    {"name": "github.com/aws/aws-sdk-go-v2", "status": "available"},
    {"name": "Recall.ai", "status": "build_client", "note": "REST API only"}
  ],
  "internal": [
    {"package": "agent/custom", "depends_on": ["tts/", "stt/"]}
  ]
}
```

### 5. Status Enum with Custom Legend

Standardized status values with customizable emoji/description:

| Status | Default Emoji | Description |
|--------|---------------|-------------|
| `completed` | ✅ | Done, optionally with `version` and `completed_date` |
| `in_progress` | 🚧 | Currently being worked on |
| `planned` | 📋 | Scheduled for future work |
| `future` | 💡 | Long-term consideration, not yet scheduled |

Custom legends can override defaults:

```json
"legend": {
  "completed": {"emoji": "✓", "description": "Done"},
  "in_progress": {"emoji": "⏳", "description": "Working"}
}
```

### 6. Phases with Status

Phases can have their own status to show overall phase progress:

```json
"phases": [
  {
    "id": "phase1",
    "name": "Phase 1: Foundation",
    "status": "completed",
    "order": 1,
    "description": "Core interfaces and essential implementations"
  }
]
```

### 7. Sections for Freeform Content

Non-item content with ordering support:

```json
"sections": [
  {
    "id": "design-philosophy",
    "title": "Design Philosophy",
    "order": 1,
    "content": [...]
  },
  {
    "id": "contributing",
    "title": "Contributing",
    "order": 99,
    "content": [...]
  }
]
```

Standard section types: `overview`, `design-philosophy`, `architecture`, `security`, `contributing`, `future`, `resources`.

### 8. Version History

Track milestone versions with dates and status:

```json
"version_history": [
  {"version": "0.5.0", "date": "2026-01-10", "status": "completed", "summary": "Sync engine"},
  {"version": "0.6.0", "date": null, "status": "planned", "summary": "Cloud backends"}
]
```

## CLI Commands

### `sroadmap validate`

Validate ROADMAP.json against the schema.

```bash
sroadmap validate ROADMAP.json
```

### `sroadmap generate`

Generate ROADMAP.md from ROADMAP.json.

```bash
sroadmap generate -i ROADMAP.json -o ROADMAP.md
```

Options:

| Flag | Description |
|------|-------------|
| `--group-by` | Grouping strategy: `category`, `status`, `phase`, `quarter` |
| `--show-completed` | Include completed items (default: true) |
| `--checkboxes` | Use `[x]`/`[ ]` syntax for tasks |
| `--emoji` | Include emoji prefixes for status |

### `sroadmap stats`

Show roadmap statistics.

```bash
sroadmap stats ROADMAP.json
```

Output:

```
Total items: 15
  Completed: 5 (33%)
  In Progress: 2 (13%)
  Planned: 6 (40%)
  Future: 2 (13%)

By Category:
  High Value: 8
  Medium Value: 4
  Improvements: 3
```

### `sroadmap deps`

Generate dependency graph (Mermaid format).

```bash
sroadmap deps ROADMAP.json --format mermaid
```

## Rendering Options

### Group by Category (default)

```markdown
## High Value

### 1. Request Timeouts ✅

Per-request timeout configuration.

**Status:** Implemented in v0.10.0

### 2. Fallback Providers

Automatic failover when primary provider fails.

- [ ] Design fallback API
- [ ] Implement provider health checks
```

### Group by Status

```markdown
## Completed (v0.10.0)

- [x] Request Timeouts

## In Progress

- [ ] Extended sampling parameters

## Planned

- [ ] Fallback Providers
- [ ] Rate Limiting
```

### Group by Quarter

```markdown
## Q1 2026

- Request Timeouts ✅
- Extended Sampling Parameters

## Q2 2026

- Fallback Providers
- Rate Limiting
```

## File Structure

```
structured-roadmap/
├── cmd/
│   └── sroadmap/
│       ├── main.go
│       ├── root.go
│       ├── validate.go
│       ├── generate.go
│       ├── stats.go
│       └── deps.go
├── roadmap/
│   ├── types.go          # IR types
│   ├── parse.go          # JSON parsing
│   ├── validate.go       # Schema validation
│   └── roadmap_test.go
├── renderer/
│   ├── markdown.go       # Markdown generation
│   ├── options.go        # Rendering options
│   └── renderer_test.go
├── schema/
│   └── roadmap.v1.schema.json
├── examples/
│   ├── minimal.json
│   ├── full.json
│   └── omnillm.json
├── PRD.md
├── README.md
├── go.mod
└── go.sum
```

## Success Criteria

1. **Deterministic output** - Same JSON always produces identical Markdown
2. **Round-trip fidelity** - All features from reviewed roadmaps can be represented
3. **CLI usability** - Simple commands for common operations
4. **Extensibility** - Schema supports future additions without breaking changes

## Future Enhancements

- GitHub Projects export
- Jira/Linear integration
- Progress tracking over time
- Multi-repo aggregation
- Web dashboard
