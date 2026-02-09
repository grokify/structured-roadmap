package renderer

import (
	"os"
	"strings"
	"testing"

	"github.com/grokify/structured-tasks/tasks"
)

func TestRender(t *testing.T) {
	tl := &tasks.TaskList{
		IRVersion: "1.0",
		Project:   "Test Project",
		Areas: []tasks.Area{
			{ID: "core", Name: "Core Features"},
			{ID: "api", Name: "API"},
		},
		Tasks: []tasks.Task{
			{ID: "task-1", Title: "Feature 1", Status: tasks.StatusCompleted, Area: "core", Type: "Added"},
			{ID: "task-2", Title: "Feature 2", Status: tasks.StatusPlanned, Area: "core", Type: "Added"},
			{ID: "task-3", Title: "Feature 3", Status: tasks.StatusInProgress, Area: "api", Type: "Changed"},
		},
	}

	t.Run("default options", func(t *testing.T) {
		opts := DefaultOptions()
		output := Render(tl, opts)

		if !strings.Contains(output, "# Task List") {
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
		output := Render(tl, opts)

		if strings.Contains(output, "[x]") || strings.Contains(output, "[ ]") {
			t.Error("Should not contain checkboxes")
		}
	})

	t.Run("no emoji", func(t *testing.T) {
		opts := DefaultOptions().WithEmoji(false)
		output := Render(tl, opts)

		if strings.Contains(output, "✅") || strings.Contains(output, "🚧") || strings.Contains(output, "📋") {
			t.Error("Should not contain emoji")
		}
	})

	t.Run("with legend", func(t *testing.T) {
		opts := DefaultOptions().WithLegend(true)
		output := Render(tl, opts)

		if !strings.Contains(output, "## Legend") {
			t.Error("Expected legend section in output")
		}
	})

	t.Run("with numbered items", func(t *testing.T) {
		opts := DefaultOptions().WithNumberedItems(true)
		output := Render(tl, opts)

		if !strings.Contains(output, "1.") {
			t.Error("Expected numbered items in output")
		}
	})
}

func TestRenderGroupBy(t *testing.T) {
	tl := &tasks.TaskList{
		IRVersion: "1.0",
		Project:   "Test",
		Areas: []tasks.Area{
			{ID: "core", Name: "Core"},
		},
		Tasks: []tasks.Task{
			{ID: "1", Title: "Task 1", Status: tasks.StatusCompleted, Area: "core", Phase: 1, Type: "Added"},
			{ID: "2", Title: "Task 2", Status: tasks.StatusPlanned, Area: "core", Phase: 1, Type: "Changed"},
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := DefaultOptions().WithGroupBy(tt.groupBy)
			output := Render(tl, opts)

			if !strings.Contains(output, tt.expect) {
				t.Errorf("Expected %q in output for group by %s", tt.expect, tt.groupBy)
			}
		})
	}
}

func TestRenderPhaseWithAreaSubheadings(t *testing.T) {
	tl := &tasks.TaskList{
		IRVersion: "1.0",
		Project:   "Test",
		Areas: []tasks.Area{
			{ID: "core", Name: "Core Package"},
			{ID: "api", Name: "API Layer"},
		},
		Tasks: []tasks.Task{
			{ID: "1", Title: "interfaces.go", Status: tasks.StatusCompleted, Area: "core", Phase: 1},
			{ID: "2", Title: "client.go", Status: tasks.StatusCompleted, Area: "api", Phase: 1},
		},
	}

	opts := DefaultOptions().WithGroupBy(GroupByPhase)
	opts.ShowAreaSubheadings = true
	output := Render(tl, opts)

	if !strings.Contains(output, "## Phase 1") {
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
	tl := &tasks.TaskList{
		IRVersion: "1.0",
		Project:   "Test",
		Areas: []tasks.Area{
			{ID: "core", Name: "Core"},
		},
		Tasks: []tasks.Task{
			{ID: "1", Title: "Task 1", Status: tasks.StatusCompleted, Area: "core"},
			{ID: "2", Title: "Task 2", Status: tasks.StatusPlanned, Area: "core"},
		},
	}

	opts := DefaultOptions()
	opts.ShowTOC = true
	output := Render(tl, opts)

	if !strings.Contains(output, "## Table of Contents") {
		t.Error("Expected TOC section")
	}
	if !strings.Contains(output, "(1/2)") {
		t.Error("Expected progress count in TOC")
	}
}

func TestRenderOverviewTable(t *testing.T) {
	tl := &tasks.TaskList{
		IRVersion: "1.0",
		Project:   "Test",
		Areas: []tasks.Area{
			{ID: "core", Name: "Core"},
		},
		Tasks: []tasks.Task{
			{ID: "1", Title: "Task 1", Status: tasks.StatusCompleted, Area: "core"},
			{ID: "2", Title: "Task 2", Status: tasks.StatusPlanned, Area: "core"},
		},
	}

	opts := DefaultOptions()
	opts.ShowOverviewTable = true
	output := Render(tl, opts)

	if !strings.Contains(output, "## Status") {
		t.Error("Expected status table section")
	}
	if !strings.Contains(output, "| Task | Status | Phase | Area |") {
		t.Error("Expected table headers")
	}
}

func TestRenderSubtasks(t *testing.T) {
	tl := &tasks.TaskList{
		IRVersion: "1.0",
		Project:   "Test",
		Tasks: []tasks.Task{
			{
				ID:     "1",
				Title:  "Feature with subtasks",
				Status: tasks.StatusInProgress,
				Subtasks: []tasks.Subtask{
					{Description: "Completed subtask", Completed: true},
					{Description: "Pending subtask", Completed: false},
				},
			},
		},
	}

	opts := DefaultOptions()
	output := Render(tl, opts)

	if !strings.Contains(output, "[x] Completed subtask") {
		t.Error("Expected completed subtask")
	}
	if !strings.Contains(output, "[ ] Pending subtask") {
		t.Error("Expected pending subtask")
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
	tl := &tasks.TaskList{
		IRVersion: "1.0",
		Project:   "Test",
		Tasks: []tasks.Task{
			{ID: "1", Title: "Task 1", Status: tasks.StatusCompleted},
		},
	}

	opts := DefaultOptions()
	tmpFile := t.TempDir() + "/tasks.md"

	err := RenderToFile(tmpFile, tl, opts)
	if err != nil {
		t.Fatalf("RenderToFile() error = %v", err)
	}

	content, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	if !strings.Contains(string(content), "# Task List") {
		t.Error("Expected task list title in output file")
	}
}

func TestRenderTOCDepth(t *testing.T) {
	tl := &tasks.TaskList{
		IRVersion: "1.0",
		Project:   "Test",
		Areas: []tasks.Area{
			{ID: "core", Name: "Core"},
		},
		Tasks: []tasks.Task{
			{ID: "1", Title: "First Task", Status: tasks.StatusCompleted, Area: "core"},
			{ID: "2", Title: "Second Task", Status: tasks.StatusPlanned, Area: "core"},
		},
	}

	t.Run("TOC depth 1", func(t *testing.T) {
		opts := DefaultOptions()
		opts.ShowTOC = true
		opts.TOCDepth = 1
		output := Render(tl, opts)

		if !strings.Contains(output, "## Table of Contents") {
			t.Error("Expected TOC section")
		}
		if !strings.Contains(output, "- [Core (") {
			t.Error("Expected area in TOC")
		}
	})

	t.Run("TOC depth 2", func(t *testing.T) {
		opts := DefaultOptions()
		opts.ShowTOC = true
		opts.TOCDepth = 2
		output := Render(tl, opts)

		if !strings.Contains(output, "## Table of Contents") {
			t.Error("Expected TOC section")
		}
	})
}

func TestRenderNoIntro(t *testing.T) {
	tl := &tasks.TaskList{
		IRVersion: "1.0",
		Project:   "Test",
		Tasks: []tasks.Task{
			{ID: "1", Title: "Task 1", Status: tasks.StatusCompleted},
		},
	}

	opts := DefaultOptions()
	opts.ShowIntro = false
	output := Render(tl, opts)

	if !strings.Contains(output, "# Task List") {
		t.Error("Expected title")
	}
}

func TestRenderNoHorizontalRules(t *testing.T) {
	tl := &tasks.TaskList{
		IRVersion: "1.0",
		Project:   "Test",
		Areas: []tasks.Area{
			{ID: "core", Name: "Core"},
			{ID: "api", Name: "API"},
		},
		Tasks: []tasks.Task{
			{ID: "1", Title: "Task 1", Status: tasks.StatusCompleted, Area: "core"},
			{ID: "2", Title: "Task 2", Status: tasks.StatusCompleted, Area: "api"},
		},
	}

	opts := DefaultOptions()
	opts.HorizontalRules = false
	output := Render(tl, opts)

	if strings.Contains(output, "\n---\n") {
		t.Error("Should not contain horizontal rules")
	}
}

func TestRenderEmptyTaskList(t *testing.T) {
	tl := &tasks.TaskList{
		IRVersion: "1.0",
		Project:   "Empty Project",
	}

	opts := DefaultOptions()
	output := Render(tl, opts)

	if !strings.Contains(output, "# Task List") {
		t.Error("Expected title even for empty task list")
	}
	if !strings.Contains(output, "**Project:** Empty Project") {
		t.Error("Expected project name even for empty task list")
	}
}

func TestRenderTasksWithUnspecifiedArea(t *testing.T) {
	tl := &tasks.TaskList{
		IRVersion: "1.0",
		Project:   "Test",
		Tasks: []tasks.Task{
			{ID: "1", Title: "Task without area", Status: tasks.StatusCompleted},
		},
	}

	opts := DefaultOptions()
	output := Render(tl, opts)

	if !strings.Contains(output, "Task without area") {
		t.Error("Expected task to be rendered")
	}
}

func TestRenderFilterCompleted(t *testing.T) {
	tl := &tasks.TaskList{
		IRVersion: "1.0",
		Project:   "Test",
		Tasks: []tasks.Task{
			{ID: "1", Title: "Completed Task", Status: tasks.StatusCompleted},
			{ID: "2", Title: "Planned Task", Status: tasks.StatusPlanned},
		},
	}

	opts := DefaultOptions()
	opts.ShowCompleted = false
	output := Render(tl, opts)

	if strings.Contains(output, "Completed Task") {
		t.Error("Should not show completed tasks when ShowCompleted is false")
	}
	if !strings.Contains(output, "Planned Task") {
		t.Error("Should still show planned tasks")
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
		{"  Spaces  Around  ", "spaces-around"},
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

func TestRenderByPhaseUnphased(t *testing.T) {
	tl := &tasks.TaskList{
		IRVersion: "1.0",
		Project:   "Test",
		Areas: []tasks.Area{
			{ID: "core", Name: "Core"},
		},
		Tasks: []tasks.Task{
			{ID: "1", Title: "Phased Task", Status: tasks.StatusPlanned, Phase: 1, Area: "core"},
			{ID: "2", Title: "Unphased Task", Status: tasks.StatusPlanned, Area: "core"},
		},
	}

	opts := DefaultOptions()
	opts.GroupBy = GroupByPhase
	output := Render(tl, opts)

	if !strings.Contains(output, "## Phase 1") {
		t.Error("Expected Phase 1 section")
	}
	if !strings.Contains(output, "## Unphased") {
		t.Error("Expected Unphased section for tasks without phase")
	}
}

func TestRenderByStatusHideCompleted(t *testing.T) {
	tl := &tasks.TaskList{
		IRVersion: "1.0",
		Project:   "Test",
		Tasks: []tasks.Task{
			{ID: "1", Title: "Done", Status: tasks.StatusCompleted},
			{ID: "2", Title: "Active", Status: tasks.StatusInProgress},
		},
	}

	opts := DefaultOptions()
	opts.GroupBy = GroupByStatus
	opts.ShowTOC = true
	opts.ShowCompleted = false
	output := Render(tl, opts)

	if strings.Contains(output, "Done") {
		t.Error("Completed tasks should be hidden when ShowCompleted is false")
	}
}

func TestIsTaskCompleteAllSubtasksDone(t *testing.T) {
	task := tasks.Task{
		ID:     "1",
		Title:  "Test",
		Status: tasks.StatusInProgress,
		Subtasks: []tasks.Subtask{
			{Description: "Subtask 1", Completed: true},
			{Description: "Subtask 2", Completed: true},
		},
	}

	if !isTaskComplete(task) {
		t.Error("Task with all subtasks complete should be complete")
	}
}

func TestIsTaskCompleteNoSubtasks(t *testing.T) {
	task := tasks.Task{
		ID:     "1",
		Title:  "Test",
		Status: tasks.StatusInProgress,
	}

	if isTaskComplete(task) {
		t.Error("Task without subtasks and not completed status should not be complete")
	}
}
