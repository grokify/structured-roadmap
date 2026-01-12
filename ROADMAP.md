# Roadmap

**Project:** Structured Roadmap

## Summary

| Item | Status | Priority | Area |
|------|--------|----------|------|
| [CLI with validate, generate, stats, deps](#cli-commands) | ✅ | High | CLI |
| [GoReleaser configuration](#goreleaser) | ✅ | High | Distribution |
| [Homebrew tap distribution](#homebrew) | ✅ | High | Distribution |
| [JSON IR schema (v1.0)](#json-ir) | ✅ | High | Core |
| [Multiple grouping strategies](#grouping-strategies) | ✅ | High | Renderer |
| [Renderer package](#renderer-pkg) | ✅ | High | Renderer |
| [Roadmap package](#roadmap-pkg) | ✅ | High | Core |
| [Dependency graph generation](#deps-graph) | ✅ | Medium | CLI |
| [JSON Schema for validation](#json-schema) | ✅ | Medium | Core |
| [Overview table](#overview-table) | ✅ | Medium | Renderer |
| [Overview table sorting](#overview-sorting) | ✅ | Medium | Renderer |
| [Phased roadmaps with area sub-sections](#phased-roadmaps) | ✅ | Medium | Renderer |
| [Rich content blocks](#rich-content) | ✅ | Medium | Renderer |
| [Table of contents with progress counts](#toc-progress) | ✅ | Medium | Renderer |
| [Two-dimensional categorization](#two-dim-categorization) | ✅ | Medium | Core |
| [Concise priority labels](#concise-priority) | ✅ | Low | Renderer |
| [Version command](#version-cmd) | ✅ | Low | CLI |
| [Claude Code plugin](#claude-plugin) | 📋 | High | Integrations |
| [Embed dependency graph in Markdown](#embed-mermaid) | 📋 | High | Renderer |
| [GitHub Issues/Projects sync](#github-sync) | 📋 | High | Integrations |
| [`sroadmap init` command](#init-cmd) | 📋 | High | CLI |
| [HTML output format](#html-output) | 📋 | Medium | Renderer |
| [Jira import](#jira-import) | 💡 | Medium | Integrations |
| [Progress visualization](#progress-viz) | 📋 | Medium | Renderer |
| [Structured Changelog sync](#schangelog-sync) | 📋 | Medium | Integrations |
| [Velocity tracking](#velocity-tracking) | 💡 | Medium | Core |
| [Watch mode for auto-regeneration](#watch-mode) | 📋 | Medium | CLI |
| [`sroadmap diff` command](#diff-cmd) | 📋 | Medium | CLI |
| [Linear sync](#linear-sync) | 💡 | Low | Integrations |
| [Multi-project aggregation](#multi-project) | 💡 | Low | Core |
| [Notion export](#notion-export) | 💡 | Low | Integrations |
| [Overdue item alerts](#overdue-alerts) | 💡 | Low | CLI |
| [Stakeholder view filtering](#stakeholder-filter) | 📋 | Low | Renderer |
| [Timeline/Gantt view](#timeline-view) | 📋 | Low | Renderer |
| [`sroadmap migrate` command](#migrate-cmd) | 📋 | Low | CLI |

---

## Overview <a href="#roadmap">↑ Top</a>

Structured Roadmap provides a machine-readable JSON intermediate representation for project roadmaps with deterministic Markdown generation. It is designed to complement [Structured Changelog](https://github.com/grokify/structured-changelog).

---

## Table of Contents

- [v0.1.0 - Initial Release (12/12)](#v010-initial-release)
- [v0.2.0 - Distribution (5/5)](#v020-distribution)
- [v0.3.0 - Enhanced Output (0/4)](#v030-enhanced-output)
- [v0.4.0 - Workflow Improvements (0/5)](#v040-workflow-improvements)
- [v0.5.0 - Integrations (0/3)](#v050-integrations)
- [Future (0/6)](#future)
- [References (3)](#references)

---

## Legend

| Status | Description |
|--------|-------------|
| ✅ | Completed |
| 🚧 | In Progress |
| 📋 | Planned |
| 💡 | Under Consideration |

---

## v0.1.0 - Initial Release ✅ <a href="#roadmap">↑ Top</a>

Machine-readable JSON roadmaps with deterministic Markdown generation

<a id="cli-commands"></a>

### [x] CLI with validate, generate, stats, deps

Core CLI subcommands for roadmap management

**Version:** 0.1.0

<a id="json-ir"></a>

### [x] JSON IR schema (v1.0)

Machine-readable roadmap format with rich metadata

**Version:** 0.1.0

<a id="grouping-strategies"></a>

### [x] Multiple grouping strategies

Group by area, type, phase, status, quarter, priority

**Version:** 0.1.0

<a id="renderer-pkg"></a>

### [x] Renderer package

Deterministic Markdown generation

**Version:** 0.1.0

<a id="roadmap-pkg"></a>

### [x] Roadmap package

IR types, parsing, and validation

**Version:** 0.1.0

<a id="deps-graph"></a>

### [x] Dependency graph generation

Mermaid and DOT format graph output

**Version:** 0.1.0

<a id="json-schema"></a>

### [x] JSON Schema for validation

Schema-based IR validation

**Version:** 0.1.0

<a id="overview-table"></a>

### [x] Overview table

Summary table of all items

**Version:** 0.1.0

<a id="phased-roadmaps"></a>

### [x] Phased roadmaps with area sub-sections

Support for large projects with hierarchical structure

**Version:** 0.1.0

<a id="rich-content"></a>

### [x] Rich content blocks

Text, code, diagram, table, list, blockquote

**Version:** 0.1.0

<a id="toc-progress"></a>

### [x] Table of contents with progress counts

TOC showing completion status per section

**Version:** 0.1.0

<a id="two-dim-categorization"></a>

### [x] Two-dimensional categorization

Area (project component) + Type (change type)

**Version:** 0.1.0

---

## v0.2.0 - Distribution ✅ <a href="#roadmap">↑ Top</a>

GoReleaser, Homebrew, and enhanced table sorting

<a id="goreleaser"></a>

### [x] GoReleaser configuration

Multi-platform binary releases (Linux, macOS, Windows)

**Version:** 0.2.0

<a id="homebrew"></a>

### [x] Homebrew tap distribution

Install via `brew install grokify/tap/sroadmap`

**Version:** 0.2.0

<a id="overview-sorting"></a>

### [x] Overview table sorting

Sort by completion percentage then priority

**Version:** 0.2.0

<a id="concise-priority"></a>

### [x] Concise priority labels

P0, P1, P2, P3 format in table cells

**Version:** 0.2.0

<a id="version-cmd"></a>

### [x] Version command

Show build info (version, commit, date)

**Version:** 0.2.0

---

## v0.3.0 - Enhanced Output 📋 <a href="#roadmap">↑ Top</a>

Additional output formats and embedded visualizations

<a id="embed-mermaid"></a>

### [ ] Embed dependency graph in Markdown

Option to include Mermaid diagram in generated ROADMAP.md

**Target:** 0.3.0

<a id="html-output"></a>

### [ ] HTML output format

Generate standalone HTML with styling

**Target:** 0.3.0

<a id="progress-viz"></a>

### [ ] Progress visualization

Progress bars or burndown charts in output

**Target:** 0.3.0

<a id="timeline-view"></a>

### [ ] Timeline/Gantt view

Generate Gantt-style timeline from target dates

**Target:** 0.3.0

---

## v0.4.0 - Workflow Improvements 📋 <a href="#roadmap">↑ Top</a>

Better authoring and maintenance workflows

<a id="init-cmd"></a>

### [ ] `sroadmap init` command

Create starter ROADMAP.json interactively

**Target:** 0.4.0

<a id="watch-mode"></a>

### [ ] Watch mode for auto-regeneration

`sroadmap generate --watch` to auto-regenerate on changes

**Target:** 0.4.0

<a id="diff-cmd"></a>

### [ ] `sroadmap diff` command

Compare two roadmap versions and show changes

**Target:** 0.4.0

<a id="stakeholder-filter"></a>

### [ ] Stakeholder view filtering

Filter output by audience (dev, product, exec)

**Target:** 0.4.0

<a id="migrate-cmd"></a>

### [ ] `sroadmap migrate` command

Migrate ROADMAP.json between schema versions

**Target:** 0.4.0

---

## v0.5.0 - Integrations 📋 <a href="#roadmap">↑ Top</a>

External tool integrations

<a id="claude-plugin"></a>

### [ ] Claude Code plugin

Plugin for AI-assisted roadmap management

**Target:** 0.5.0

<a id="github-sync"></a>

### [ ] GitHub Issues/Projects sync

Import from and export to GitHub Issues or Projects

**Target:** 0.5.0

<a id="schangelog-sync"></a>

### [ ] Structured Changelog sync

Auto-move completed items to CHANGELOG.json

**Target:** 0.5.0

---

## Future 💡 <a href="#roadmap">↑ Top</a>

Ideas under consideration

<a id="jira-import"></a>

### [ ] Jira import

Import roadmap items from Jira epics/stories

<a id="velocity-tracking"></a>

### [ ] Velocity tracking

Track completion velocity and predict timelines

<a id="linear-sync"></a>

### [ ] Linear sync

Bidirectional sync with Linear projects

<a id="multi-project"></a>

### [ ] Multi-project aggregation

Aggregate roadmaps from multiple projects

<a id="notion-export"></a>

### [ ] Notion export

Export roadmap to Notion database

<a id="overdue-alerts"></a>

### [ ] Overdue item alerts

Highlight items past their target date

---

## References <a href="#roadmap">↑ Top</a>

- [PRD.md](PRD.md) - Product requirements
- [CHANGELOG.md](CHANGELOG.md) - Release history
- [Structured Changelog](https://github.com/grokify/structured-changelog) - Companion project

---

## Version History <a href="#roadmap">↑ Top</a>

| Version | Date | Status | Summary |
|---------|------|--------|--------|
| 0.1.0 | 2026-01-10 | ✅ | Initial release with JSON IR and Markdown generation |
| 0.2.0 | 2026-01-11 | ✅ | GoReleaser, Homebrew, enhanced sorting |
| 0.3.0 | TBD | 📋 | Enhanced output formats |
| 0.4.0 | TBD | 📋 | Workflow improvements |
| 0.5.0 | TBD | 📋 | External integrations |
