// Package tasks provides types and utilities for structured task list IR.
// The Type field uses category names from structured-changelog for consistency.
package tasks

// Status represents the status of a task.
type Status string

const (
	StatusInProgress Status = "in_progress"
	StatusPlanned    Status = "planned"
	StatusFuture     Status = "future"
	StatusCompleted  Status = "completed"
)

// DefaultLegend returns the default status legend with emoji and descriptions.
func DefaultLegend() map[Status]LegendEntry {
	return map[Status]LegendEntry{
		StatusInProgress: {Emoji: "🚧", Description: "In Progress"},
		StatusPlanned:    {Emoji: "📋", Description: "Planned"},
		StatusFuture:     {Emoji: "💡", Description: "Under Consideration"},
		StatusCompleted:  {Emoji: "✅", Description: "Completed"},
	}
}

// TaskList is the top-level IR structure for a project task list.
type TaskList struct {
	IRVersion string                 `json:"irVersion"`
	Project   string                 `json:"project"`
	Legend    map[Status]LegendEntry `json:"legend,omitempty"`
	Areas     []Area                 `json:"areas,omitempty"`
	Tasks     []Task                 `json:"tasks,omitempty"`
}

// LegendEntry defines the emoji and description for a status.
type LegendEntry struct {
	Emoji       string `json:"emoji"`
	Description string `json:"description"`
}

// Area represents a project area/component for grouping tasks.
type Area struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Task represents a work item (feature, task, improvement).
// Order is determined by position in the Tasks array.
// Type should be a valid category name from structured-changelog (e.g., "Added", "Fixed").
type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Status      Status    `json:"status"`
	Phase       int       `json:"phase,omitempty"`
	Area        string    `json:"area,omitempty"`
	Type        string    `json:"type,omitempty"`
	DependsOn   []string  `json:"dependsOn,omitempty"`
	Blocks      []string  `json:"blocks,omitempty"`
	Subtasks    []Subtask `json:"subtasks,omitempty"`
}

// Subtask represents a checkbox item within a task.
type Subtask struct {
	ID          string `json:"id,omitempty"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

// GetLegend returns the task list's legend, falling back to defaults.
func (tl *TaskList) GetLegend() map[Status]LegendEntry {
	if len(tl.Legend) > 0 {
		legend := DefaultLegend()
		for k, v := range tl.Legend {
			legend[k] = v
		}
		return legend
	}
	return DefaultLegend()
}

// GetStatusEmoji returns the emoji for a status.
func (tl *TaskList) GetStatusEmoji(status Status) string {
	legend := tl.GetLegend()
	if entry, ok := legend[status]; ok {
		return entry.Emoji
	}
	return ""
}

// TasksByArea returns tasks grouped by area.
func (tl *TaskList) TasksByArea() map[string][]Task {
	result := make(map[string][]Task)
	for _, task := range tl.Tasks {
		area := task.Area
		if area == "" {
			area = "_unspecified"
		}
		result[area] = append(result[area], task)
	}
	return result
}

// TasksByType returns tasks grouped by change type.
func (tl *TaskList) TasksByType() map[string][]Task {
	result := make(map[string][]Task)
	for _, task := range tl.Tasks {
		t := task.Type
		if t == "" {
			t = "_unspecified"
		}
		result[t] = append(result[t], task)
	}
	return result
}

// TasksByPhase returns tasks grouped by phase number.
func (tl *TaskList) TasksByPhase() map[int][]Task {
	result := make(map[int][]Task)
	for _, task := range tl.Tasks {
		result[task.Phase] = append(result[task.Phase], task)
	}
	return result
}

// TasksByStatus returns tasks grouped by status.
func (tl *TaskList) TasksByStatus() map[Status][]Task {
	result := make(map[Status][]Task)
	for _, task := range tl.Tasks {
		result[task.Status] = append(result[task.Status], task)
	}
	return result
}

// Stats returns statistics about the task list.
func (tl *TaskList) Stats() Stats {
	stats := Stats{
		ByStatus: make(map[Status]int),
		ByArea:   make(map[string]int),
		ByType:   make(map[string]int),
		ByPhase:  make(map[int]int),
	}
	stats.Total = len(tl.Tasks)
	for _, task := range tl.Tasks {
		stats.ByStatus[task.Status]++
		if task.Area != "" {
			stats.ByArea[task.Area]++
		}
		if task.Type != "" {
			stats.ByType[task.Type]++
		}
		stats.ByPhase[task.Phase]++
	}
	return stats
}

// Stats holds task list statistics.
type Stats struct {
	Total    int
	ByStatus map[Status]int
	ByArea   map[string]int
	ByType   map[string]int
	ByPhase  map[int]int
}

// InProgressCount returns the number of in-progress tasks.
func (s Stats) InProgressCount() int {
	return s.ByStatus[StatusInProgress]
}

// PlannedCount returns the number of planned tasks.
func (s Stats) PlannedCount() int {
	return s.ByStatus[StatusPlanned]
}

// CompletedCount returns the number of completed tasks.
func (s Stats) CompletedCount() int {
	return s.ByStatus[StatusCompleted]
}

// StatusOrder returns the canonical order of statuses for display.
func StatusOrder() []Status {
	return []Status{StatusInProgress, StatusPlanned, StatusFuture, StatusCompleted}
}

// PhaseNumbers returns sorted phase numbers from the task list.
func (tl *TaskList) PhaseNumbers() []int {
	phases := make(map[int]bool)
	for _, task := range tl.Tasks {
		if task.Phase > 0 {
			phases[task.Phase] = true
		}
	}
	result := make([]int, 0, len(phases))
	for p := range phases {
		result = append(result, p)
	}
	// Sort phases
	for i := 0; i < len(result)-1; i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i] > result[j] {
				result[i], result[j] = result[j], result[i]
			}
		}
	}
	return result
}
