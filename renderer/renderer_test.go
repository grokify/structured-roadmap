package renderer

import (
	"os"
	"strings"
	"testing"

	"github.com/grokify/structured-roadmap/roadmap"
)

func TestRender(t *testing.T) {
	r := &roadmap.Roadmap{
		IRVersion: "1.0",
		Project:   "Test Project",
		Areas: []roadmap.Area{
			{ID: "core", Name: "Core Features", Priority: 1},
			{ID: "api", Name: "API", Priority: 2},
		},
		Items: []roadmap.Item{
			{ID: "item-1", Title: "Feature 1", Status: roadmap.StatusCompleted, Area: "core", Type: "Added"},
			{ID: "item-2", Title: "Feature 2", Status: roadmap.StatusPlanned, Area: "core", Type: "Added"},
			{ID: "item-3", Title: "Feature 3", Status: roadmap.StatusInProgress, Area: "api", Type: "Changed"},
		},
	}

	t.Run("default options", func(t *testing.T) {
		opts := DefaultOptions()
		output := Render(r, opts)

		if !strings.Contains(output, "# Roadmap") {
			t.Error("Expected title in output")
		}
		if !strings.Contains(output, "**Project:** Test Project") {
			t.Error("Expected project name in output")
		}
		if !strings.Contains(output, "## Core Features") {
			t.Error("Expected area heading in output")
		}
		if !strings.Contains(output, "[x]") {
			t.Error("Expected completed checkbox in output")
		}
		if !strings.Contains(output, "[ ]") {
			t.Error("Expected uncompleted checkbox in output")
		}
	})

	t.Run("no checkboxes", func(t *testing.T) {
		opts := DefaultOptions().WithCheckboxes(false)
		output := Render(r, opts)

		if strings.Contains(output, "[x]") || strings.Contains(output, "[ ]") {
			t.Error("Should not contain checkboxes")
		}
	})

	t.Run("no emoji", func(t *testing.T) {
		opts := DefaultOptions().WithEmoji(false)
		output := Render(r, opts)

		if strings.Contains(output, "✅") || strings.Contains(output, "🚧") || strings.Contains(output, "📋") {
			t.Error("Should not contain emoji")
		}
	})

	t.Run("with legend", func(t *testing.T) {
		opts := DefaultOptions().WithLegend(true)
		output := Render(r, opts)

		if !strings.Contains(output, "## Legend") {
			t.Error("Expected legend section in output")
		}
	})

	t.Run("with numbered items", func(t *testing.T) {
		opts := DefaultOptions().WithNumberedItems(true)
		output := Render(r, opts)

		if !strings.Contains(output, "1.") {
			t.Error("Expected numbered items in output")
		}
	})
}

func TestRenderGroupBy(t *testing.T) {
	r := &roadmap.Roadmap{
		IRVersion: "1.0",
		Project:   "Test",
		Areas: []roadmap.Area{
			{ID: "core", Name: "Core", Priority: 1},
		},
		Phases: []roadmap.Phase{
			{ID: "phase-1", Name: "Phase 1", Status: roadmap.StatusCompleted, Order: 1},
		},
		Items: []roadmap.Item{
			{ID: "1", Title: "Item 1", Status: roadmap.StatusCompleted, Area: "core", Phase: "phase-1", Priority: roadmap.PriorityHigh, TargetQuarter: "Q1 2026", Type: "Added"},
			{ID: "2", Title: "Item 2", Status: roadmap.StatusPlanned, Area: "core", Phase: "phase-1", Priority: roadmap.PriorityMedium, TargetQuarter: "Q2 2026", Type: "Changed"},
		},
	}

	tests := []struct {
		name    string
		groupBy GroupBy
		expect  string
	}{
		{"by area", GroupByArea, "## Core"},
		{"by type", GroupByType, "## Added"},
		{"by phase", GroupByPhase, "## Phase 1"},
		{"by status", GroupByStatus, "## ✅ Completed"},
		{"by quarter", GroupByQuarter, "## Q1 2026"},
		{"by priority", GroupByPriority, "## High Priority"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := DefaultOptions().WithGroupBy(tt.groupBy)
			output := Render(r, opts)

			if !strings.Contains(output, tt.expect) {
				t.Errorf("Expected %q in output for group by %s", tt.expect, tt.groupBy)
			}
		})
	}
}

func TestRenderPhaseWithAreaSubheadings(t *testing.T) {
	r := &roadmap.Roadmap{
		IRVersion: "1.0",
		Project:   "Test",
		Areas: []roadmap.Area{
			{ID: "core", Name: "Core Package", Priority: 1},
			{ID: "api", Name: "API Layer", Priority: 2},
		},
		Phases: []roadmap.Phase{
			{ID: "phase-1", Name: "Phase 1: Foundation", Status: roadmap.StatusCompleted, Order: 1},
		},
		Items: []roadmap.Item{
			{ID: "1", Title: "interfaces.go", Status: roadmap.StatusCompleted, Area: "core", Phase: "phase-1"},
			{ID: "2", Title: "client.go", Status: roadmap.StatusCompleted, Area: "api", Phase: "phase-1"},
		},
	}

	opts := DefaultOptions().WithGroupBy(GroupByPhase)
	opts.ShowAreaSubheadings = true
	output := Render(r, opts)

	if !strings.Contains(output, "## Phase 1: Foundation") {
		t.Error("Expected phase heading")
	}
	if !strings.Contains(output, "### Core Package") {
		t.Error("Expected area subheading")
	}
	if !strings.Contains(output, "### API Layer") {
		t.Error("Expected area subheading")
	}
}

func TestRenderTOC(t *testing.T) {
	r := &roadmap.Roadmap{
		IRVersion: "1.0",
		Project:   "Test",
		Areas: []roadmap.Area{
			{ID: "core", Name: "Core", Priority: 1},
		},
		Items: []roadmap.Item{
			{ID: "1", Title: "Item 1", Status: roadmap.StatusCompleted, Area: "core"},
			{ID: "2", Title: "Item 2", Status: roadmap.StatusPlanned, Area: "core"},
		},
	}

	opts := DefaultOptions()
	opts.ShowTOC = true
	output := Render(r, opts)

	if !strings.Contains(output, "## Table of Contents") {
		t.Error("Expected TOC section")
	}
	if !strings.Contains(output, "(1/2)") {
		t.Error("Expected progress count in TOC")
	}
}

func TestRenderOverviewTable(t *testing.T) {
	r := &roadmap.Roadmap{
		IRVersion: "1.0",
		Project:   "Test",
		Areas: []roadmap.Area{
			{ID: "core", Name: "Core", Priority: 1},
		},
		Items: []roadmap.Item{
			{ID: "1", Title: "Item 1", Status: roadmap.StatusCompleted, Area: "core", Priority: roadmap.PriorityHigh},
			{ID: "2", Title: "Item 2", Status: roadmap.StatusPlanned, Area: "core", Priority: roadmap.PriorityMedium},
		},
	}

	opts := DefaultOptions()
	opts.ShowOverviewTable = true
	output := Render(r, opts)

	if !strings.Contains(output, "## Summary") {
		t.Error("Expected summary table section")
	}
	if !strings.Contains(output, "| Item | Status | Priority | Area |") {
		t.Error("Expected table headers")
	}
}

func TestRenderContentBlocks(t *testing.T) {
	r := &roadmap.Roadmap{
		IRVersion: "1.0",
		Project:   "Test",
		Items: []roadmap.Item{
			{
				ID:     "1",
				Title:  "Item with content",
				Status: roadmap.StatusCompleted,
				Content: []roadmap.ContentBlock{
					{Type: roadmap.ContentTypeText, Value: "Some text content"},
					{Type: roadmap.ContentTypeCode, Value: "fmt.Println(\"hello\")", Language: "go"},
					{Type: roadmap.ContentTypeBlockquote, Value: "A blockquote"},
					{Type: roadmap.ContentTypeList, Items: []string{"Item A", "Item B"}},
					{Type: roadmap.ContentTypeTable, Headers: []string{"Col1", "Col2"}, Rows: [][]string{{"A", "B"}}},
				},
			},
		},
	}

	opts := DefaultOptions()
	output := Render(r, opts)

	if !strings.Contains(output, "Some text content") {
		t.Error("Expected text content")
	}
	if !strings.Contains(output, "```go") {
		t.Error("Expected code block with language")
	}
	if !strings.Contains(output, "> A blockquote") {
		t.Error("Expected blockquote")
	}
	if !strings.Contains(output, "- Item A") {
		t.Error("Expected list items")
	}
	if !strings.Contains(output, "| Col1 | Col2 |") {
		t.Error("Expected table")
	}
}

func TestRenderSections(t *testing.T) {
	r := &roadmap.Roadmap{
		IRVersion: "1.0",
		Project:   "Test",
		Sections: []roadmap.Section{
			{
				ID:    "overview",
				Title: "Overview",
				Order: 1,
				Content: []roadmap.ContentBlock{
					{Type: roadmap.ContentTypeText, Value: "Project overview text."},
				},
			},
		},
	}

	opts := DefaultOptions()
	output := Render(r, opts)

	if !strings.Contains(output, "## Overview") {
		t.Error("Expected section heading")
	}
	if !strings.Contains(output, "Project overview text.") {
		t.Error("Expected section content")
	}
}

func TestRenderVersionHistory(t *testing.T) {
	r := &roadmap.Roadmap{
		IRVersion: "1.0",
		Project:   "Test",
		VersionHistory: []roadmap.VersionEntry{
			{Version: "1.0.0", Date: "2026-01-01", Status: roadmap.StatusCompleted, Summary: "Initial release"},
			{Version: "2.0.0", Status: roadmap.StatusPlanned, Summary: "Major update"},
		},
	}

	opts := DefaultOptions()
	output := Render(r, opts)

	if !strings.Contains(output, "## Version History") {
		t.Error("Expected version history section")
	}
	if !strings.Contains(output, "1.0.0") {
		t.Error("Expected version in output")
	}
}

func TestRenderTasks(t *testing.T) {
	r := &roadmap.Roadmap{
		IRVersion: "1.0",
		Project:   "Test",
		Items: []roadmap.Item{
			{
				ID:     "1",
				Title:  "Feature with tasks",
				Status: roadmap.StatusInProgress,
				Tasks: []roadmap.Task{
					{Description: "Completed task", Completed: true},
					{Description: "Pending task", Completed: false},
				},
			},
		},
	}

	opts := DefaultOptions()
	output := Render(r, opts)

	if !strings.Contains(output, "[x] Completed task") {
		t.Error("Expected completed task")
	}
	if !strings.Contains(output, "[ ] Pending task") {
		t.Error("Expected pending task")
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := DefaultOptions()

	if opts.GroupBy != GroupByArea {
		t.Errorf("Default GroupBy = %s, want %s", opts.GroupBy, GroupByArea)
	}
	if !opts.ShowCompleted {
		t.Error("ShowCompleted should default to true")
	}
	if !opts.UseCheckboxes {
		t.Error("UseCheckboxes should default to true")
	}
	if !opts.UseEmoji {
		t.Error("UseEmoji should default to true")
	}
	if !opts.ShowIntro {
		t.Error("ShowIntro should default to true")
	}
	if opts.ShowLegend {
		t.Error("ShowLegend should default to false")
	}
	if opts.ShowTOC {
		t.Error("ShowTOC should default to false")
	}
	if !opts.HorizontalRules {
		t.Error("HorizontalRules should default to true")
	}
}

func TestRenderToFile(t *testing.T) {
	r := &roadmap.Roadmap{
		IRVersion: "1.0",
		Project:   "Test",
		Items: []roadmap.Item{
			{ID: "1", Title: "Item 1", Status: roadmap.StatusCompleted},
		},
	}

	opts := DefaultOptions()
	tmpFile := t.TempDir() + "/roadmap.md"

	err := RenderToFile(tmpFile, r, opts)
	if err != nil {
		t.Fatalf("RenderToFile() error = %v", err)
	}

	// Verify file was created and has content
	content, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	if !strings.Contains(string(content), "# Roadmap") {
		t.Error("Expected roadmap title in output file")
	}
}

func TestRenderDependencies(t *testing.T) {
	r := &roadmap.Roadmap{
		IRVersion: "1.0",
		Project:   "Test",
		Dependencies: &roadmap.Dependencies{
			External: []roadmap.ExternalDependency{
				{Name: "github.com/example/pkg", Status: "available"},
				{Name: "github.com/another/lib", Status: "pending", Note: "Waiting for v2"},
			},
			Internal: []roadmap.InternalDependency{
				{Package: "pkg/core", DependsOn: []string{"pkg/utils"}},
			},
		},
	}

	opts := DefaultOptions()
	opts.ShowDependencies = true
	output := Render(r, opts)

	if !strings.Contains(output, "## Dependencies") {
		t.Error("Expected dependencies section")
	}
	if !strings.Contains(output, "### External") {
		t.Error("Expected external dependencies section")
	}
	if !strings.Contains(output, "github.com/example/pkg") {
		t.Error("Expected external dependency name")
	}
	if !strings.Contains(output, "### Internal") {
		t.Error("Expected internal dependencies section")
	}
}

func TestRenderTOCDepth(t *testing.T) {
	r := &roadmap.Roadmap{
		IRVersion: "1.0",
		Project:   "Test",
		Areas: []roadmap.Area{
			{ID: "core", Name: "Core", Priority: 1},
		},
		Items: []roadmap.Item{
			{ID: "1", Title: "First Item", Status: roadmap.StatusCompleted, Area: "core"},
			{ID: "2", Title: "Second Item", Status: roadmap.StatusPlanned, Area: "core"},
		},
	}

	t.Run("TOC depth 1", func(t *testing.T) {
		opts := DefaultOptions()
		opts.ShowTOC = true
		opts.TOCDepth = 1
		output := Render(r, opts)

		if !strings.Contains(output, "## Table of Contents") {
			t.Error("Expected TOC section")
		}
		// TOC format includes count like: - [Core (0/1)](#core)
		if !strings.Contains(output, "- [Core (") {
			t.Error("Expected area in TOC")
		}
	})

	t.Run("TOC depth 2", func(t *testing.T) {
		opts := DefaultOptions()
		opts.ShowTOC = true
		opts.TOCDepth = 2
		output := Render(r, opts)

		if !strings.Contains(output, "## Table of Contents") {
			t.Error("Expected TOC section")
		}
	})
}

func TestRenderNoIntro(t *testing.T) {
	r := &roadmap.Roadmap{
		IRVersion: "1.0",
		Project:   "Test",
		Items: []roadmap.Item{
			{ID: "1", Title: "Item 1", Status: roadmap.StatusCompleted},
		},
	}

	opts := DefaultOptions()
	opts.ShowIntro = false
	output := Render(r, opts)

	// Should still have title but no intro paragraph
	if !strings.Contains(output, "# Roadmap") {
		t.Error("Expected title")
	}
}

func TestRenderNoHorizontalRules(t *testing.T) {
	r := &roadmap.Roadmap{
		IRVersion: "1.0",
		Project:   "Test",
		Areas: []roadmap.Area{
			{ID: "core", Name: "Core", Priority: 1},
			{ID: "api", Name: "API", Priority: 2},
		},
		Items: []roadmap.Item{
			{ID: "1", Title: "Item 1", Status: roadmap.StatusCompleted, Area: "core"},
			{ID: "2", Title: "Item 2", Status: roadmap.StatusCompleted, Area: "api"},
		},
	}

	opts := DefaultOptions()
	opts.HorizontalRules = false
	output := Render(r, opts)

	if strings.Contains(output, "\n---\n") {
		t.Error("Should not contain horizontal rules")
	}
}

func TestRenderItemWithAllFields(t *testing.T) {
	r := &roadmap.Roadmap{
		IRVersion: "1.0",
		Project:   "Test",
		Items: []roadmap.Item{
			{
				ID:            "1",
				Title:         "Complete Feature",
				Description:   "A detailed description",
				Status:        roadmap.StatusCompleted,
				Version:       "1.0.0",
				CompletedDate: "2026-01-01",
				TargetQuarter: "Q1 2026",
				TargetVersion: "1.0.0",
				Priority:      roadmap.PriorityHigh,
				Type:          "Added",
				DependsOn:     []string{"other-item"},
			},
		},
	}

	opts := DefaultOptions()
	output := Render(r, opts)

	if !strings.Contains(output, "Complete Feature") {
		t.Error("Expected item title")
	}
	if !strings.Contains(output, "A detailed description") {
		t.Error("Expected item description")
	}
}

func TestRenderItemsAsTasks(t *testing.T) {
	r := &roadmap.Roadmap{
		IRVersion: "1.0",
		Project:   "Test",
		Areas: []roadmap.Area{
			{ID: "core", Name: "Core", Priority: 1},
		},
		Phases: []roadmap.Phase{
			{ID: "phase-1", Name: "Phase 1", Status: roadmap.StatusInProgress, Order: 1},
		},
		Items: []roadmap.Item{
			{ID: "1", Title: "Task 1", Status: roadmap.StatusCompleted, Area: "core", Phase: "phase-1"},
			{ID: "2", Title: "Task 2", Status: roadmap.StatusInProgress, Area: "core", Phase: "phase-1"},
			{ID: "3", Title: "Task 3", Status: roadmap.StatusPlanned, Area: "core", Phase: "phase-1"},
		},
	}

	opts := DefaultOptions().WithGroupBy(GroupByPhase)
	opts.ShowAreaSubheadings = true
	output := Render(r, opts)

	// Items should be rendered as task list within area subheadings
	if !strings.Contains(output, "### Core") {
		t.Error("Expected area subheading")
	}
	if !strings.Contains(output, "[x] Task 1") {
		t.Error("Expected completed task")
	}
	if !strings.Contains(output, "[ ] Task 2") || !strings.Contains(output, "[ ] Task 3") {
		t.Error("Expected incomplete tasks")
	}
}

func TestRenderMultipleSections(t *testing.T) {
	r := &roadmap.Roadmap{
		IRVersion: "1.0",
		Project:   "Test",
		Sections: []roadmap.Section{
			{
				ID:    "intro",
				Title: "Introduction",
				Order: 1,
				Content: []roadmap.ContentBlock{
					{Type: roadmap.ContentTypeText, Value: "Welcome to the project."},
				},
			},
			{
				ID:    "philosophy",
				Title: "Design Philosophy",
				Order: 2,
				Content: []roadmap.ContentBlock{
					{Type: roadmap.ContentTypeText, Value: "We believe in simplicity."},
					{Type: roadmap.ContentTypeBlockquote, Value: "Keep it simple."},
				},
			},
			{
				ID:    "future",
				Title: "Future Plans",
				Order: 3,
				Content: []roadmap.ContentBlock{
					{Type: roadmap.ContentTypeList, Items: []string{"Feature A", "Feature B"}},
				},
			},
		},
	}

	opts := DefaultOptions()
	opts.ShowSections = true
	output := Render(r, opts)

	if !strings.Contains(output, "## Introduction") {
		t.Error("Expected Introduction section")
	}
	if !strings.Contains(output, "## Design Philosophy") {
		t.Error("Expected Design Philosophy section")
	}
	if !strings.Contains(output, "## Future Plans") {
		t.Error("Expected Future Plans section")
	}
	if !strings.Contains(output, "> Keep it simple.") {
		t.Error("Expected blockquote in section")
	}
}

func TestRenderDiagramContentBlock(t *testing.T) {
	r := &roadmap.Roadmap{
		IRVersion: "1.0",
		Project:   "Test",
		Items: []roadmap.Item{
			{
				ID:     "1",
				Title:  "Architecture",
				Status: roadmap.StatusCompleted,
				Content: []roadmap.ContentBlock{
					{Type: roadmap.ContentTypeDiagram, Value: "graph TD\n  A --> B", Format: "mermaid"},
				},
			},
		},
	}

	opts := DefaultOptions()
	output := Render(r, opts)

	// Diagrams use plain code fence (not language-specific)
	if !strings.Contains(output, "```\n") {
		t.Error("Expected code fence for diagram")
	}
	if !strings.Contains(output, "A --> B") {
		t.Error("Expected diagram content")
	}
}

func TestRenderCodeBlockNoLanguage(t *testing.T) {
	r := &roadmap.Roadmap{
		IRVersion: "1.0",
		Project:   "Test",
		Items: []roadmap.Item{
			{
				ID:     "1",
				Title:  "Example",
				Status: roadmap.StatusCompleted,
				Content: []roadmap.ContentBlock{
					{Type: roadmap.ContentTypeCode, Value: "some code"},
				},
			},
		},
	}

	opts := DefaultOptions()
	output := Render(r, opts)

	if !strings.Contains(output, "```\nsome code\n```") {
		t.Error("Expected code block without language")
	}
}

func TestRenderEmptyRoadmap(t *testing.T) {
	r := &roadmap.Roadmap{
		IRVersion: "1.0",
		Project:   "Empty Project",
	}

	opts := DefaultOptions()
	output := Render(r, opts)

	if !strings.Contains(output, "# Roadmap") {
		t.Error("Expected title even for empty roadmap")
	}
	if !strings.Contains(output, "**Project:** Empty Project") {
		t.Error("Expected project name even for empty roadmap")
	}
}

func TestRenderItemsWithUnspecifiedArea(t *testing.T) {
	r := &roadmap.Roadmap{
		IRVersion: "1.0",
		Project:   "Test",
		Items: []roadmap.Item{
			{ID: "1", Title: "Item without area", Status: roadmap.StatusCompleted},
		},
	}

	opts := DefaultOptions()
	output := Render(r, opts)

	// Should still render without crashing
	if !strings.Contains(output, "Item without area") {
		t.Error("Expected item to be rendered")
	}
}

func TestRenderPhaseWithDescription(t *testing.T) {
	r := &roadmap.Roadmap{
		IRVersion: "1.0",
		Project:   "Test",
		Phases: []roadmap.Phase{
			{ID: "phase-1", Name: "Phase 1", Status: roadmap.StatusCompleted, Order: 1, Description: "Foundation phase."},
		},
		Items: []roadmap.Item{
			{ID: "1", Title: "Item 1", Status: roadmap.StatusCompleted, Phase: "phase-1"},
		},
	}

	opts := DefaultOptions().WithGroupBy(GroupByPhase)
	output := Render(r, opts)

	if !strings.Contains(output, "Foundation phase.") {
		t.Error("Expected phase description")
	}
}

func TestRenderVersionHistoryNoDate(t *testing.T) {
	r := &roadmap.Roadmap{
		IRVersion: "1.0",
		Project:   "Test",
		VersionHistory: []roadmap.VersionEntry{
			{Version: "2.0.0", Status: roadmap.StatusPlanned, Summary: "Future version"},
		},
	}

	opts := DefaultOptions()
	output := Render(r, opts)

	if !strings.Contains(output, "2.0.0") {
		t.Error("Expected version without date")
	}
	if !strings.Contains(output, "Future version") {
		t.Error("Expected version summary")
	}
}

func TestRenderSortItems(t *testing.T) {
	r := &roadmap.Roadmap{
		IRVersion: "1.0",
		Project:   "Test",
		Items: []roadmap.Item{
			{ID: "3", Title: "Third", Status: roadmap.StatusCompleted, Order: 3},
			{ID: "1", Title: "First", Status: roadmap.StatusCompleted, Order: 1},
			{ID: "2", Title: "Second", Status: roadmap.StatusCompleted, Order: 2},
		},
	}

	opts := DefaultOptions()
	output := Render(r, opts)

	// Items should be sorted by Order field
	firstIdx := strings.Index(output, "First")
	secondIdx := strings.Index(output, "Second")
	thirdIdx := strings.Index(output, "Third")

	if firstIdx > secondIdx || secondIdx > thirdIdx {
		t.Error("Expected items to be sorted by order")
	}
}

func TestRenderFilterCompleted(t *testing.T) {
	r := &roadmap.Roadmap{
		IRVersion: "1.0",
		Project:   "Test",
		Items: []roadmap.Item{
			{ID: "1", Title: "Completed Item", Status: roadmap.StatusCompleted},
			{ID: "2", Title: "Planned Item", Status: roadmap.StatusPlanned},
		},
	}

	opts := DefaultOptions()
	opts.ShowCompleted = false
	output := Render(r, opts)

	if strings.Contains(output, "Completed Item") {
		t.Error("Should not show completed items when ShowCompleted is false")
	}
	if !strings.Contains(output, "Planned Item") {
		t.Error("Should still show planned items")
	}
}

func TestRenderTaskWithFilePath(t *testing.T) {
	r := &roadmap.Roadmap{
		IRVersion: "1.0",
		Project:   "Test",
		Items: []roadmap.Item{
			{
				ID:     "1",
				Title:  "Feature",
				Status: roadmap.StatusInProgress,
				Tasks: []roadmap.Task{
					{Description: "Implement handler", Completed: false, FilePath: "internal/handler.go"},
				},
			},
		},
	}

	opts := DefaultOptions()
	output := Render(r, opts)

	if !strings.Contains(output, "internal/handler.go") {
		t.Error("Expected file path in task")
	}
}

func TestSlugify(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Simple Text", "simple-text"},
		{"Text With Numbers 123", "text-with-numbers-123"},
		{"Special!@#Characters", "specialcharacters"},
		{"  Spaces  Around  ", "spaces-around"}, // multiple hyphens are collapsed
		{"UPPERCASE", "uppercase"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := slugify(tt.input)
			if result != tt.expected {
				t.Errorf("slugify(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestBuildTOCEntriesByPhase(t *testing.T) {
	r := &roadmap.Roadmap{
		IRVersion: "1.0",
		Project:   "Test",
		Phases: []roadmap.Phase{
			{ID: "phase-1", Name: "Phase 1", Order: 1},
			{ID: "phase-2", Name: "Phase 2", Order: 2},
		},
		Items: []roadmap.Item{
			{ID: "1", Title: "Item 1", Status: roadmap.StatusCompleted, Phase: "phase-1"},
			{ID: "2", Title: "Item 2", Status: roadmap.StatusPlanned, Phase: "phase-2"},
		},
	}

	opts := DefaultOptions()
	opts.ShowTOC = true
	opts.GroupBy = GroupByPhase
	output := Render(r, opts)

	if !strings.Contains(output, "## Table of Contents") {
		t.Error("Expected TOC section")
	}
	if !strings.Contains(output, "Phase 1") {
		t.Error("Expected Phase 1 in TOC")
	}
}

func TestBuildTOCEntriesByStatus(t *testing.T) {
	r := &roadmap.Roadmap{
		IRVersion: "1.0",
		Project:   "Test",
		Items: []roadmap.Item{
			{ID: "1", Title: "Done Item", Status: roadmap.StatusCompleted},
			{ID: "2", Title: "Planned Item", Status: roadmap.StatusPlanned},
			{ID: "3", Title: "In Progress Item", Status: roadmap.StatusInProgress},
		},
	}

	opts := DefaultOptions()
	opts.ShowTOC = true
	opts.GroupBy = GroupByStatus
	opts.ShowCompleted = true
	output := Render(r, opts)

	if !strings.Contains(output, "## Table of Contents") {
		t.Error("Expected TOC section")
	}
	if !strings.Contains(output, "Completed") {
		t.Error("Expected Completed in output")
	}
	if !strings.Contains(output, "Planned") {
		t.Error("Expected Planned in output")
	}
}

func TestBuildTOCEntriesByQuarter(t *testing.T) {
	r := &roadmap.Roadmap{
		IRVersion: "1.0",
		Project:   "Test",
		Items: []roadmap.Item{
			{ID: "1", Title: "Q1 Item", Status: roadmap.StatusPlanned, TargetQuarter: "Q1 2026"},
			{ID: "2", Title: "Q2 Item", Status: roadmap.StatusPlanned, TargetQuarter: "Q2 2026"},
			{ID: "3", Title: "Unscheduled Item", Status: roadmap.StatusFuture},
		},
	}

	opts := DefaultOptions()
	opts.ShowTOC = true
	opts.GroupBy = GroupByQuarter
	output := Render(r, opts)

	if !strings.Contains(output, "## Table of Contents") {
		t.Error("Expected TOC section")
	}
	if !strings.Contains(output, "Q1 2026") {
		t.Error("Expected Q1 2026 in output")
	}
	if !strings.Contains(output, "Unscheduled") {
		t.Error("Expected Unscheduled in output")
	}
}

func TestBuildTOCEntriesByPriority(t *testing.T) {
	r := &roadmap.Roadmap{
		IRVersion: "1.0",
		Project:   "Test",
		Items: []roadmap.Item{
			{ID: "1", Title: "Critical Item", Status: roadmap.StatusPlanned, Priority: roadmap.PriorityCritical},
			{ID: "2", Title: "High Item", Status: roadmap.StatusPlanned, Priority: roadmap.PriorityHigh},
			{ID: "3", Title: "Low Item", Status: roadmap.StatusPlanned, Priority: roadmap.PriorityLow},
		},
	}

	opts := DefaultOptions()
	opts.ShowTOC = true
	opts.GroupBy = GroupByPriority
	output := Render(r, opts)

	if !strings.Contains(output, "## Table of Contents") {
		t.Error("Expected TOC section")
	}
	if !strings.Contains(output, "Critical") {
		t.Error("Expected Critical in output")
	}
	if !strings.Contains(output, "High Priority") {
		t.Error("Expected High Priority in output")
	}
}

func TestBuildTOCEntriesByType(t *testing.T) {
	r := &roadmap.Roadmap{
		IRVersion: "1.0",
		Project:   "Test",
		Items: []roadmap.Item{
			{ID: "1", Title: "New Feature", Status: roadmap.StatusPlanned, Type: "Added"},
			{ID: "2", Title: "Bug Fix", Status: roadmap.StatusPlanned, Type: "Fixed"},
		},
	}

	opts := DefaultOptions()
	opts.ShowTOC = true
	opts.GroupBy = GroupByType
	output := Render(r, opts)

	if !strings.Contains(output, "## Table of Contents") {
		t.Error("Expected TOC section")
	}
	if !strings.Contains(output, "Added") {
		t.Error("Expected Added type in output")
	}
}

func TestRenderWithSectionInTOC(t *testing.T) {
	r := &roadmap.Roadmap{
		IRVersion: "1.0",
		Project:   "Test",
		Areas: []roadmap.Area{
			{ID: "core", Name: "Core", Priority: 1},
		},
		Items: []roadmap.Item{
			{ID: "1", Title: "Item 1", Status: roadmap.StatusPlanned, Area: "core"},
		},
		Sections: []roadmap.Section{
			{
				ID:    "notes",
				Title: "Notes",
				Order: 1,
				Content: []roadmap.ContentBlock{
					{Type: roadmap.ContentTypeList, Items: []string{"Note 1", "Note 2"}},
				},
			},
		},
	}

	opts := DefaultOptions()
	opts.ShowTOC = true
	opts.ShowSections = true
	output := Render(r, opts)

	if !strings.Contains(output, "## Table of Contents") {
		t.Error("Expected TOC section")
	}
	if !strings.Contains(output, "Notes") {
		t.Error("Expected Notes section in TOC")
	}
	if !strings.Contains(output, "(2)") {
		t.Error("Expected list item count in section TOC entry")
	}
}

func TestRenderItemsAsTasksWithSubheadings(t *testing.T) {
	r := &roadmap.Roadmap{
		IRVersion: "1.0",
		Project:   "Test",
		Phases: []roadmap.Phase{
			{ID: "phase-1", Name: "Phase 1", Order: 1, Description: "First phase"},
		},
		Areas: []roadmap.Area{
			{ID: "core", Name: "Core", Priority: 1},
			{ID: "api", Name: "API", Priority: 2},
		},
		Items: []roadmap.Item{
			{
				ID:          "1",
				Title:       "Core Feature",
				Description: "A core feature",
				Status:      roadmap.StatusPlanned,
				Phase:       "phase-1",
				Area:        "core",
				Tasks: []roadmap.Task{
					{Description: "Task 1", Completed: true},
					{Description: "Task 2", Completed: false, FilePath: "file.go"},
				},
			},
			{
				ID:     "2",
				Title:  "API Feature",
				Status: roadmap.StatusInProgress,
				Phase:  "phase-1",
				Area:   "api",
			},
		},
	}

	opts := DefaultOptions()
	opts.GroupBy = GroupByPhase
	opts.ShowAreaSubheadings = true
	opts.UseCheckboxes = true
	output := Render(r, opts)

	if !strings.Contains(output, "### Core") {
		t.Error("Expected Core subheading")
	}
	if !strings.Contains(output, "### API") {
		t.Error("Expected API subheading")
	}
	if !strings.Contains(output, "[x] Task 1") {
		t.Error("Expected completed task checkbox")
	}
	if !strings.Contains(output, "[ ] Task 2") {
		t.Error("Expected incomplete task checkbox")
	}
	if !strings.Contains(output, "(`file.go`)") {
		t.Error("Expected file path in task")
	}
}

func TestRenderItemsAsTasksWithoutCheckboxes(t *testing.T) {
	r := &roadmap.Roadmap{
		IRVersion: "1.0",
		Project:   "Test",
		Phases: []roadmap.Phase{
			{ID: "phase-1", Name: "Phase 1", Order: 1},
		},
		Areas: []roadmap.Area{
			{ID: "core", Name: "Core", Priority: 1},
		},
		Items: []roadmap.Item{
			{
				ID:          "1",
				Title:       "Feature",
				Description: "A feature",
				Status:      roadmap.StatusCompleted,
				Phase:       "phase-1",
				Area:        "core",
			},
		},
	}

	opts := DefaultOptions()
	opts.GroupBy = GroupByPhase
	opts.ShowAreaSubheadings = true
	opts.UseCheckboxes = false
	opts.ShowCompleted = true
	output := Render(r, opts)

	if !strings.Contains(output, "### Core") {
		t.Error("Expected Core subheading")
	}
	if !strings.Contains(output, "- Feature - A feature") {
		t.Error("Expected item with description")
	}
}

func TestRenderItemAllFields(t *testing.T) {
	r := &roadmap.Roadmap{
		IRVersion: "1.0",
		Project:   "Test",
		Areas: []roadmap.Area{
			{ID: "core", Name: "Core", Priority: 1},
		},
		Items: []roadmap.Item{
			{
				ID:            "1",
				Title:         "Full Item",
				Description:   "This item has all fields",
				Status:        roadmap.StatusInProgress,
				Version:       "1.0.0",
				CompletedDate: "2024-01-15",
				TargetVersion: "2.0.0",
				TargetQuarter: "Q1 2026",
				Area:          "core",
				Type:          "added",
				Priority:      roadmap.PriorityHigh,
				Order:         1,
				Tasks: []roadmap.Task{
					{Description: "Subtask", Completed: false},
				},
				Content: []roadmap.ContentBlock{
					{Type: roadmap.ContentTypeText, Value: "Additional info"},
				},
			},
		},
	}

	opts := DefaultOptions()
	opts.UseEmoji = true
	opts.UseCheckboxes = false
	output := Render(r, opts)

	if !strings.Contains(output, "Full Item") {
		t.Error("Expected item title")
	}
	if !strings.Contains(output, "**Version:** 1.0.0") {
		t.Error("Expected version info")
	}
	if !strings.Contains(output, "(2024-01-15)") {
		t.Error("Expected completed date")
	}
}

func TestRenderItemTargetOnly(t *testing.T) {
	r := &roadmap.Roadmap{
		IRVersion: "1.0",
		Project:   "Test",
		Areas: []roadmap.Area{
			{ID: "core", Name: "Core", Priority: 1},
		},
		Items: []roadmap.Item{
			{
				ID:            "1",
				Title:         "Planned Item",
				Status:        roadmap.StatusPlanned,
				TargetVersion: "2.0.0",
				TargetQuarter: "Q2 2026",
				Area:          "core",
			},
		},
	}

	opts := DefaultOptions()
	output := Render(r, opts)

	if !strings.Contains(output, "**Target:** 2.0.0") {
		t.Error("Expected target version")
	}
	if !strings.Contains(output, "(Q2 2026)") {
		t.Error("Expected target quarter")
	}
}

func TestRenderItemQuarterOnly(t *testing.T) {
	r := &roadmap.Roadmap{
		IRVersion: "1.0",
		Project:   "Test",
		Areas: []roadmap.Area{
			{ID: "core", Name: "Core", Priority: 1},
		},
		Items: []roadmap.Item{
			{
				ID:            "1",
				Title:         "Quarterly Item",
				Status:        roadmap.StatusPlanned,
				TargetQuarter: "Q3 2026",
				Area:          "core",
			},
		},
	}

	opts := DefaultOptions()
	output := Render(r, opts)

	if !strings.Contains(output, "**Target:** Q3 2026") {
		t.Error("Expected target quarter")
	}
}

func TestRenderByUnphasedItems(t *testing.T) {
	r := &roadmap.Roadmap{
		IRVersion: "1.0",
		Project:   "Test",
		Phases: []roadmap.Phase{
			{ID: "phase-1", Name: "Phase 1", Order: 1},
		},
		Areas: []roadmap.Area{
			{ID: "core", Name: "Core", Priority: 1},
		},
		Items: []roadmap.Item{
			{ID: "1", Title: "Phased Item", Status: roadmap.StatusPlanned, Phase: "phase-1", Area: "core"},
			{ID: "2", Title: "Unphased Item", Status: roadmap.StatusPlanned, Area: "core"},
		},
	}

	opts := DefaultOptions()
	opts.GroupBy = GroupByPhase
	opts.ShowAreaSubheadings = true
	output := Render(r, opts)

	if !strings.Contains(output, "## Other") {
		t.Error("Expected Other section for unphased items")
	}
}

func TestRenderVersionHistoryEmoji(t *testing.T) {
	r := &roadmap.Roadmap{
		IRVersion: "1.0",
		Project:   "Test",
		VersionHistory: []roadmap.VersionEntry{
			{Version: "1.0.0", Date: "2024-01-15", Status: roadmap.StatusCompleted, Summary: "Initial release"},
			{Version: "2.0.0", Status: roadmap.StatusPlanned, Summary: "Major update"},
		},
	}

	opts := DefaultOptions()
	opts.ShowVersionHistory = true
	opts.UseEmoji = true
	output := Render(r, opts)

	if !strings.Contains(output, "## Version History") {
		t.Error("Expected version history section")
	}
	if !strings.Contains(output, "TBD") {
		t.Error("Expected TBD for missing date")
	}
}

func TestRenderOverviewTableVariations(t *testing.T) {
	r := &roadmap.Roadmap{
		IRVersion: "1.0",
		Project:   "Test",
		Areas: []roadmap.Area{
			{ID: "core", Name: "Core", Priority: 1},
		},
		Items: []roadmap.Item{
			{ID: "1", Title: "Item 1", Status: roadmap.StatusCompleted, Area: "core", Priority: roadmap.PriorityHigh, Order: 1},
			{ID: "2", Title: "Item 2", Status: roadmap.StatusPlanned, Order: 2},
		},
	}

	t.Run("without emoji", func(t *testing.T) {
		opts := DefaultOptions()
		opts.ShowOverviewTable = true
		opts.UseEmoji = false
		output := Render(r, opts)

		if !strings.Contains(output, "| completed |") {
			t.Error("Expected status text in overview table")
		}
	})

	t.Run("with emoji", func(t *testing.T) {
		opts := DefaultOptions()
		opts.ShowOverviewTable = true
		opts.UseEmoji = true
		output := Render(r, opts)

		if !strings.Contains(output, "| ✅ |") {
			t.Error("Expected emoji in overview table")
		}
	})

	t.Run("hide completed", func(t *testing.T) {
		opts := DefaultOptions()
		opts.ShowOverviewTable = true
		opts.ShowCompleted = false
		output := Render(r, opts)

		if strings.Contains(output, "Item 1") {
			t.Error("Completed item should be hidden")
		}
	})
}

func TestRenderByStatusHideCompleted(t *testing.T) {
	r := &roadmap.Roadmap{
		IRVersion: "1.0",
		Project:   "Test",
		Items: []roadmap.Item{
			{ID: "1", Title: "Done", Status: roadmap.StatusCompleted},
			{ID: "2", Title: "Active", Status: roadmap.StatusInProgress},
		},
	}

	opts := DefaultOptions()
	opts.GroupBy = GroupByStatus
	opts.ShowTOC = true
	opts.ShowCompleted = false
	output := Render(r, opts)

	if strings.Contains(output, "## ") && strings.Contains(output, "Completed") && strings.Contains(output, "Done") {
		t.Error("Completed section should be hidden when ShowCompleted is false")
	}
}

func TestRenderTOCWithNumberedItems(t *testing.T) {
	r := &roadmap.Roadmap{
		IRVersion: "1.0",
		Project:   "Test",
		Areas: []roadmap.Area{
			{ID: "core", Name: "Core", Priority: 1},
		},
		Items: []roadmap.Item{
			{ID: "1", Title: "Item 1", Status: roadmap.StatusPlanned, Area: "core"},
			{ID: "2", Title: "Item 2", Status: roadmap.StatusPlanned, Area: "core"},
		},
	}

	opts := DefaultOptions()
	opts.ShowTOC = true
	opts.TOCDepth = 2
	opts.NumberItems = true
	output := Render(r, opts)

	if !strings.Contains(output, "1. Item 1") {
		t.Error("Expected numbered item in TOC")
	}
}

func TestIsItemCompleteAllTasksDone(t *testing.T) {
	item := roadmap.Item{
		ID:     "1",
		Title:  "Test",
		Status: roadmap.StatusInProgress,
		Tasks: []roadmap.Task{
			{Description: "Task 1", Completed: true},
			{Description: "Task 2", Completed: true},
		},
	}

	if !isItemComplete(item) {
		t.Error("Item with all tasks complete should be complete")
	}
}

func TestIsItemCompleteNoTasks(t *testing.T) {
	item := roadmap.Item{
		ID:     "1",
		Title:  "Test",
		Status: roadmap.StatusInProgress,
	}

	if isItemComplete(item) {
		t.Error("Item without tasks and not completed status should not be complete")
	}
}

func TestRenderSectionWithNoListContent(t *testing.T) {
	r := &roadmap.Roadmap{
		IRVersion: "1.0",
		Project:   "Test",
		Sections: []roadmap.Section{
			{
				ID:    "notes",
				Title: "Notes",
				Order: 1,
				Content: []roadmap.ContentBlock{
					{Type: roadmap.ContentTypeText, Value: "Some text"},
				},
			},
		},
	}

	opts := DefaultOptions()
	opts.ShowTOC = true
	opts.ShowSections = true
	output := Render(r, opts)

	// Section without list items should not show count
	if strings.Contains(output, "Notes (") && strings.Contains(output, ")") && !strings.Contains(output, "Notes](#notes)") {
		t.Error("Section with no list content should not show count in TOC")
	}
}
