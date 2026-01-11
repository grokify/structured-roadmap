// Package roadmap provides types and utilities for structured roadmap IR.
package roadmap

import (
	"time"
)

// Status represents the status of a roadmap item or phase.
type Status string

const (
	StatusCompleted  Status = "completed"
	StatusInProgress Status = "in_progress"
	StatusPlanned    Status = "planned"
	StatusFuture     Status = "future"
)

// Priority represents the priority level of a roadmap item.
type Priority string

const (
	PriorityCritical Priority = "critical"
	PriorityHigh     Priority = "high"
	PriorityMedium   Priority = "medium"
	PriorityLow      Priority = "low"
)

// PriorityOrder returns the sort order for a priority (lower = higher priority).
func PriorityOrder(p Priority) int {
	switch p {
	case PriorityCritical:
		return 1
	case PriorityHigh:
		return 2
	case PriorityMedium:
		return 3
	case PriorityLow:
		return 4
	default:
		return 5
	}
}

// PriorityLabel returns a concise label for a priority.
// Use this in table cells where the column header provides context.
func PriorityLabel(p Priority) string {
	switch p {
	case PriorityCritical:
		return "Critical"
	case PriorityHigh:
		return "High"
	case PriorityMedium:
		return "Medium"
	case PriorityLow:
		return "Low"
	default:
		return "-"
	}
}

// PriorityLabelFull returns a full label for a priority.
// Use this in section headers and standalone contexts.
func PriorityLabelFull(p Priority) string {
	switch p {
	case PriorityCritical:
		return "Critical"
	case PriorityHigh:
		return "High Priority"
	case PriorityMedium:
		return "Medium Priority"
	case PriorityLow:
		return "Low Priority"
	default:
		return "Unspecified"
	}
}

// DefaultLegend returns the default status legend with emoji and descriptions.
func DefaultLegend() map[Status]LegendEntry {
	return map[Status]LegendEntry{
		StatusCompleted:  {Emoji: "✅", Description: "Completed"},
		StatusInProgress: {Emoji: "🚧", Description: "In Progress"},
		StatusPlanned:    {Emoji: "📋", Description: "Planned"},
		StatusFuture:     {Emoji: "💡", Description: "Under Consideration"},
	}
}

// Roadmap is the top-level IR structure for a project roadmap.
type Roadmap struct {
	IRVersion      string                 `json:"ir_version"`
	Project        string                 `json:"project"`
	Repository     string                 `json:"repository,omitempty"`
	GeneratedAt    *time.Time             `json:"generated_at,omitempty"`
	Legend         map[Status]LegendEntry `json:"legend,omitempty"`
	Areas          []Area                 `json:"areas,omitempty"`
	Phases         []Phase                `json:"phases,omitempty"`
	Items          []Item                 `json:"items,omitempty"`
	Sections       []Section              `json:"sections,omitempty"`
	VersionHistory []VersionEntry         `json:"version_history,omitempty"`
	Dependencies   *Dependencies          `json:"dependencies,omitempty"`
}

// LegendEntry defines the emoji and description for a status.
type LegendEntry struct {
	Emoji       string `json:"emoji"`
	Description string `json:"description"`
}

// Area represents a project area/component for grouping items.
type Area struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Priority int    `json:"priority,omitempty"`
}

// Phase represents a development phase.
type Phase struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Status      Status `json:"status,omitempty"`
	Order       int    `json:"order,omitempty"`
	Description string `json:"description,omitempty"`
}

// Item represents a roadmap item (feature, task, improvement).
type Item struct {
	ID            string         `json:"id"`
	Title         string         `json:"title"`
	Description   string         `json:"description,omitempty"`
	Status        Status         `json:"status"`
	Version       string         `json:"version,omitempty"`
	CompletedDate string         `json:"completed_date,omitempty"`
	TargetQuarter string         `json:"target_quarter,omitempty"`
	TargetVersion string         `json:"target_version,omitempty"`
	Area          string         `json:"area,omitempty"` // Project area/component (user-defined)
	Type          string         `json:"type,omitempty"` // Change type (aligns with structured-changelog)
	Phase         string         `json:"phase,omitempty"`
	Priority      Priority       `json:"priority,omitempty"`
	Order         int            `json:"order,omitempty"`
	DependsOn     []string       `json:"depends_on,omitempty"`
	Tasks         []Task         `json:"tasks,omitempty"`
	Content       []ContentBlock `json:"content,omitempty"`
}

// Task represents a sub-task with completion status.
type Task struct {
	ID          string `json:"id,omitempty"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
	FilePath    string `json:"file_path,omitempty"`
}

// ContentBlock represents a rich content block within an item or section.
type ContentBlock struct {
	Type     ContentType `json:"type"`
	Value    string      `json:"value,omitempty"`
	Language string      `json:"language,omitempty"`
	Format   string      `json:"format,omitempty"`
	Headers  []string    `json:"headers,omitempty"`
	Rows     [][]string  `json:"rows,omitempty"`
	Items    []string    `json:"items,omitempty"`
}

// ContentType represents the type of a content block.
type ContentType string

const (
	ContentTypeText       ContentType = "text"
	ContentTypeCode       ContentType = "code"
	ContentTypeDiagram    ContentType = "diagram"
	ContentTypeTable      ContentType = "table"
	ContentTypeList       ContentType = "list"
	ContentTypeBlockquote ContentType = "blockquote"
)

// Section represents a freeform content section.
type Section struct {
	ID      string         `json:"id"`
	Title   string         `json:"title"`
	Order   int            `json:"order,omitempty"`
	Content []ContentBlock `json:"content,omitempty"`
}

// VersionEntry represents a version milestone.
type VersionEntry struct {
	Version string `json:"version"`
	Date    string `json:"date,omitempty"`
	Status  Status `json:"status,omitempty"`
	Summary string `json:"summary,omitempty"`
}

// Dependencies contains external and internal dependencies.
type Dependencies struct {
	External []ExternalDependency `json:"external,omitempty"`
	Internal []InternalDependency `json:"internal,omitempty"`
}

// ExternalDependency represents an external SDK dependency.
type ExternalDependency struct {
	Name   string `json:"name"`
	Status string `json:"status,omitempty"`
	Note   string `json:"note,omitempty"`
}

// InternalDependency represents an internal package dependency.
type InternalDependency struct {
	Package   string   `json:"package"`
	DependsOn []string `json:"depends_on,omitempty"`
}

// GetLegend returns the roadmap's legend, falling back to defaults.
func (r *Roadmap) GetLegend() map[Status]LegendEntry {
	if len(r.Legend) > 0 {
		legend := DefaultLegend()
		for k, v := range r.Legend {
			legend[k] = v
		}
		return legend
	}
	return DefaultLegend()
}

// GetStatusEmoji returns the emoji for a status.
func (r *Roadmap) GetStatusEmoji(status Status) string {
	legend := r.GetLegend()
	if entry, ok := legend[status]; ok {
		return entry.Emoji
	}
	return ""
}

// ItemsByArea returns items grouped by area.
func (r *Roadmap) ItemsByArea() map[string][]Item {
	result := make(map[string][]Item)
	for _, item := range r.Items {
		area := item.Area
		if area == "" {
			area = "_unspecified"
		}
		result[area] = append(result[area], item)
	}
	return result
}

// ItemsByType returns items grouped by change type.
func (r *Roadmap) ItemsByType() map[string][]Item {
	result := make(map[string][]Item)
	for _, item := range r.Items {
		t := item.Type
		if t == "" {
			t = "_unspecified"
		}
		result[t] = append(result[t], item)
	}
	return result
}

// ItemsByPhase returns items grouped by phase.
func (r *Roadmap) ItemsByPhase() map[string][]Item {
	result := make(map[string][]Item)
	for _, item := range r.Items {
		phase := item.Phase
		if phase == "" {
			phase = "_unphased"
		}
		result[phase] = append(result[phase], item)
	}
	return result
}

// ItemsByStatus returns items grouped by status.
func (r *Roadmap) ItemsByStatus() map[Status][]Item {
	result := make(map[Status][]Item)
	for _, item := range r.Items {
		result[item.Status] = append(result[item.Status], item)
	}
	return result
}

// ItemsByQuarter returns items grouped by target quarter.
func (r *Roadmap) ItemsByQuarter() map[string][]Item {
	result := make(map[string][]Item)
	for _, item := range r.Items {
		quarter := item.TargetQuarter
		if quarter == "" {
			quarter = "_unscheduled"
		}
		result[quarter] = append(result[quarter], item)
	}
	return result
}

// ItemsByPriority returns items grouped by priority level.
func (r *Roadmap) ItemsByPriority() map[Priority][]Item {
	result := make(map[Priority][]Item)
	for _, item := range r.Items {
		priority := item.Priority
		if priority == "" {
			priority = "_unspecified"
		}
		result[priority] = append(result[priority], item)
	}
	return result
}

// Stats returns statistics about the roadmap.
func (r *Roadmap) Stats() Stats {
	stats := Stats{
		ByStatus:   make(map[Status]int),
		ByArea:     make(map[string]int),
		ByType:     make(map[string]int),
		ByPriority: make(map[Priority]int),
	}
	stats.Total = len(r.Items)
	for _, item := range r.Items {
		stats.ByStatus[item.Status]++
		if item.Area != "" {
			stats.ByArea[item.Area]++
		}
		if item.Type != "" {
			stats.ByType[item.Type]++
		}
		if item.Priority != "" {
			stats.ByPriority[item.Priority]++
		}
	}
	return stats
}

// Stats holds roadmap statistics.
type Stats struct {
	Total      int
	ByStatus   map[Status]int
	ByArea     map[string]int
	ByType     map[string]int
	ByPriority map[Priority]int
}

// CompletedCount returns the number of completed items.
func (s Stats) CompletedCount() int {
	return s.ByStatus[StatusCompleted]
}

// CompletedPercent returns the percentage of completed items.
func (s Stats) CompletedPercent() float64 {
	if s.Total == 0 {
		return 0
	}
	return float64(s.CompletedCount()) / float64(s.Total) * 100
}

// StatusOrder returns the canonical order of statuses for display.
func StatusOrder() []Status {
	return []Status{StatusCompleted, StatusInProgress, StatusPlanned, StatusFuture}
}

// PriorityOrderList returns the canonical order of priorities for display.
func PriorityOrderList() []Priority {
	return []Priority{PriorityCritical, PriorityHigh, PriorityMedium, PriorityLow}
}
