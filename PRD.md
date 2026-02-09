# Structured Tasks - Product Requirements Document

## Overview

Structured Tasks provides a machine-readable JSON intermediate representation (IR) for project task lists, with deterministic Markdown generation. It is designed for both human review and AI agent manipulation, following the pattern established by [Structured Changelog](https://github.com/grokify/structured-changelog).

## Problem Statement

Project task lists across repositories use inconsistent formats, making it difficult to:

- Maintain consistency across multiple projects
- Track status programmatically
- Generate standardized documentation
- Visualize dependencies and progress
- Allow AI agents to read and update task status

## Design Principles

1. **Simple schema** - Minimal fields, no deep nesting
2. **Array ordering** - Position in array determines priority (no explicit order/priority fields)
3. **Integer phases** - Phases are simple integers (1, 2, 3) not string references
4. **Bidirectional dependencies** - Both `depends_on` and `blocks` fields
5. **Subtask IDs** - Enable precise programmatic updates
6. **Deterministic output** - Same JSON always produces identical Markdown
7. **Status table first** - Overview ordered by phase with completed items at bottom

## JSON IR Schema (v1.0)

```json
{
  "ir_version": "1.0",
  "project": "My Project",

  "legend": {
    "completed": {"emoji": "✅", "description": "Completed"},
    "in_progress": {"emoji": "🚧", "description": "In Progress"},
    "planned": {"emoji": "📋", "description": "Planned"},
    "future": {"emoji": "💡", "description": "Under Consideration"}
  },

  "areas": [
    {"id": "core", "name": "Core"},
    {"id": "api", "name": "API"}
  ],

  "tasks": [
    {
      "id": "task-1",
      "title": "Implement feature X",
      "description": "Detailed description here",
      "status": "in_progress",
      "phase": 1,
      "area": "core",
      "type": "Added",
      "depends_on": [],
      "blocks": ["task-2"],
      "subtasks": [
        {"id": "sub-1", "description": "Design API", "completed": true},
        {"id": "sub-2", "description": "Write tests", "completed": false}
      ]
    },
    {
      "id": "task-2",
      "title": "Build on feature X",
      "status": "planned",
      "phase": 2,
      "area": "api",
      "type": "Added",
      "depends_on": ["task-1"]
    }
  ]
}
```

## Type Definitions

### TaskList (root)

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `ir_version` | string | Yes | Schema version (currently "1.0") |
| `project` | string | Yes | Project name |
| `legend` | map | No | Custom status emoji/descriptions |
| `areas` | Area[] | No | Project areas for grouping |
| `tasks` | Task[] | No | List of tasks |

### Task

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Yes | Unique identifier |
| `title` | string | Yes | Short title |
| `description` | string | No | Detailed description |
| `status` | Status | Yes | Current status |
| `phase` | int | No | Phase number (1, 2, 3...) |
| `area` | string | No | Area ID reference |
| `type` | string | No | Change type (Added, Fixed, etc.) |
| `depends_on` | string[] | No | Task IDs this depends on |
| `blocks` | string[] | No | Task IDs this blocks |
| `subtasks` | Subtask[] | No | Checklist items |

### Subtask

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | No | Unique identifier for updates |
| `description` | string | Yes | Subtask description |
| `completed` | bool | Yes | Completion status |

### Status Enum

| Value | Description |
|-------|-------------|
| `in_progress` | Currently being worked on |
| `planned` | Scheduled for future work |
| `future` | Under consideration, not scheduled |
| `completed` | Done |

### Area

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Yes | Unique identifier |
| `name` | string | Yes | Display name |

## CLI Commands

### `stasks validate`

Validate TASKS.json against the schema.

```bash
stasks validate TASKS.json
```

### `stasks generate`

Generate TASKS.md from TASKS.json.

```bash
stasks generate -i TASKS.json -o TASKS.md
```

Options:

| Flag | Description |
|------|-------------|
| `--group-by` | Grouping: `area`, `status`, `phase`, `type` |
| `--toc` | Include table of contents |
| `--legend` | Include status legend |
| `--no-completed` | Hide completed tasks |

### `stasks stats`

Show task list statistics.

```bash
stasks stats TASKS.json
```

### `stasks deps`

Generate dependency graph.

```bash
stasks deps TASKS.json --format mermaid
stasks deps TASKS.json --format dot
```

## Rendered Output

The generated Markdown includes:

1. **Title and project name**
2. **Status table** - Ordered by phase, completed at bottom of each phase
3. **Table of contents** - With progress counts (e.g., "Core (4/7)")
4. **Legend** - Status emoji reference
5. **Task sections** - Grouped by area, phase, status, or type

### Status Table Example

```markdown
## Status

| Task | Status | Phase | Area |
|------|--------|-------|------|
| [Implement auth](#auth) | 🚧 | 1 | Core |
| [Add API endpoints](#api) | 📋 | 1 | API |
| [Setup CI](#ci) | ✅ | 1 | Infra |
| [Implement caching](#cache) | 📋 | 2 | Core |
```

## File Structure

```
structured-tasks/
├── cmd/stasks/          # CLI application
├── tasks/               # IR types, parsing, validation
├── renderer/            # Markdown generation
├── schema/              # JSON Schema
├── examples/            # Example TASKS.json files
├── TASKS.json           # This project's task list
├── TASKS.md             # Generated from TASKS.json
└── PRD.md               # This document
```

## Integration with Structured Changelog

- Task `type` field uses category names from structured-changelog (Added, Fixed, Changed, etc.)
- Completed tasks can be archived to CHANGELOG.json when phases are done
- Same deterministic Markdown philosophy

## AI Agent Usage

The simplified schema enables AI agents to:

1. **Parse** - Read TASKS.json to understand project state
2. **Query** - Find tasks by status, phase, or area
3. **Update** - Modify task status, add subtasks
4. **Regenerate** - Run `stasks generate` to update Markdown

Example agent workflow:

```bash
# Read current state
cat TASKS.json | jq '.tasks[] | select(.status == "in_progress")'

# Update via JSON manipulation
# ... agent modifies TASKS.json ...

# Regenerate Markdown
stasks generate -i TASKS.json -o TASKS.md
```

## Success Criteria

1. **Deterministic output** - Same JSON always produces identical Markdown
2. **Simple schema** - Easy for both humans and AI to understand
3. **CLI usability** - Simple commands for common operations
4. **Extensibility** - Schema supports future additions without breaking changes
