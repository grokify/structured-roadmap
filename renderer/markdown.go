package renderer

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/grokify/structured-changelog/changelog"
	"github.com/grokify/structured-tasks/tasks"
)

// slugify converts a heading to a GitHub-flavored markdown anchor.
func slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	re := regexp.MustCompile(`[^a-z0-9\-_]`)
	s = re.ReplaceAllString(s, "")
	re = regexp.MustCompile(`-+`)
	s = re.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}

// taskSlug returns a stable anchor slug for a task.
func taskSlug(task tasks.Task) string {
	if task.ID != "" {
		return slugify(task.ID)
	}
	return slugify(task.Title)
}

// topAnchorID returns the anchor ID for the page top.
func topAnchorID(_ string) string {
	return "task-list"
}

// renderSectionHeading writes a section heading with an optional "Top" navigation link.
func renderSectionHeading(sb *strings.Builder, title string, project string, opts Options) {
	if opts.ShowNavLinks {
		topID := topAnchorID(project)
		fmt.Fprintf(sb, "## %s <a href=\"#%s\">↑ Top</a>\n\n", title, topID)
	} else {
		fmt.Fprintf(sb, "## %s\n\n", title)
	}
}

// tocEntry represents an entry in the table of contents.
type tocEntry struct {
	Title     string
	Slug      string
	Count     int
	Completed int
	Tasks     []tocEntry
}

// Render generates Markdown from a TaskList.
func Render(tl *tasks.TaskList, opts Options) string {
	var sb strings.Builder

	// Title
	sb.WriteString("# Task List\n\n")

	// Project name
	if tl.Project != "" {
		fmt.Fprintf(&sb, "**Project:** %s\n\n", tl.Project)
	}

	// Intro text
	if opts.ShowIntro {
		intro := opts.IntroText
		if intro == "" {
			intro = DefaultIntroText
		}
		sb.WriteString(intro + "\n\n")
	}

	// Overview table
	if opts.ShowOverviewTable {
		renderOverviewTable(&sb, tl, opts)
		if opts.HorizontalRules {
			sb.WriteString("---\n\n")
		}
	}

	// Table of Contents
	if opts.ShowTOC {
		renderTOC(&sb, tl, opts)
		if opts.HorizontalRules {
			sb.WriteString("---\n\n")
		}
	}

	// Legend
	if opts.ShowLegend {
		renderLegend(&sb, tl)
		if opts.HorizontalRules {
			sb.WriteString("---\n\n")
		}
	}

	// Main content grouped by strategy
	switch opts.GroupBy {
	case GroupByPhase:
		renderByPhase(&sb, tl, opts)
	case GroupByStatus:
		renderByStatus(&sb, tl, opts)
	case GroupByType:
		renderByType(&sb, tl, opts)
	default:
		renderByArea(&sb, tl, opts)
	}

	return strings.TrimRight(sb.String(), "\n") + "\n"
}

// RenderToFile writes rendered Markdown to a file.
func RenderToFile(path string, tl *tasks.TaskList, opts Options) error {
	content := Render(tl, opts)
	return os.WriteFile(path, []byte(content), 0600)
}

func renderLegend(sb *strings.Builder, tl *tasks.TaskList) {
	sb.WriteString("## Legend\n\n")
	sb.WriteString("| Status | Description |\n")
	sb.WriteString("|--------|-------------|\n")
	legend := tl.GetLegend()
	for _, status := range tasks.StatusOrder() {
		if entry, ok := legend[status]; ok {
			fmt.Fprintf(sb, "| %s | %s |\n", entry.Emoji, entry.Description)
		}
	}
	sb.WriteString("\n")
}

// statusSortOrder returns the sort order for a status.
// In-progress first, then planned, then future, then completed at the bottom.
func statusSortOrder(s tasks.Status) int {
	switch s {
	case tasks.StatusInProgress:
		return 0
	case tasks.StatusPlanned:
		return 1
	case tasks.StatusFuture:
		return 2
	case tasks.StatusCompleted:
		return 3
	default:
		return 4
	}
}

func renderOverviewTable(sb *strings.Builder, tl *tasks.TaskList, opts Options) {
	sb.WriteString("## Status\n\n")
	sb.WriteString("| Task | Status | Phase | Area |\n")
	sb.WriteString("|------|--------|-------|------|\n")

	legend := tl.GetLegend()

	// Build area name lookup
	areaNames := make(map[string]string)
	for _, area := range tl.Areas {
		areaNames[area.ID] = area.Name
	}

	// Sort tasks: by phase first, then by status (completed at bottom), then by title
	sorted := make([]tasks.Task, len(tl.Tasks))
	copy(sorted, tl.Tasks)
	sort.Slice(sorted, func(i, j int) bool {
		// Phase first (0 = unphased goes last)
		iPhase := sorted[i].Phase
		jPhase := sorted[j].Phase
		if iPhase == 0 {
			iPhase = 9999
		}
		if jPhase == 0 {
			jPhase = 9999
		}
		if iPhase != jPhase {
			return iPhase < jPhase
		}
		// Within same phase, sort by status (completed at bottom)
		iOrder := statusSortOrder(sorted[i].Status)
		jOrder := statusSortOrder(sorted[j].Status)
		if iOrder != jOrder {
			return iOrder < jOrder
		}
		// Finally by title
		return sorted[i].Title < sorted[j].Title
	})

	for _, task := range sorted {
		if task.Status == tasks.StatusCompleted && !opts.ShowCompleted {
			continue
		}

		// Status emoji
		status := ""
		if opts.UseEmoji {
			if entry, ok := legend[task.Status]; ok {
				status = entry.Emoji
			}
		} else {
			status = string(task.Status)
		}

		// Phase
		phase := "-"
		if task.Phase > 0 {
			phase = fmt.Sprintf("%d", task.Phase)
		}

		// Area name
		areaName := areaNames[task.Area]
		if areaName == "" {
			areaName = "-"
		}

		// Task title with anchor link
		titleLink := fmt.Sprintf("[%s](#%s)", task.Title, taskSlug(task))

		fmt.Fprintf(sb, "| %s | %s | %s | %s |\n", titleLink, status, phase, areaName)
	}
	sb.WriteString("\n")
}

func renderTOC(sb *strings.Builder, tl *tasks.TaskList, opts Options) {
	sb.WriteString("## Table of Contents\n\n")

	entries := buildTOCEntries(tl, opts)

	for _, entry := range entries {
		fmt.Fprintf(sb, "- [%s (%d/%d)](#%s)\n", entry.Title, entry.Completed, entry.Count, entry.Slug)

		if opts.TOCDepth >= 2 {
			for _, task := range entry.Tasks {
				fmt.Fprintf(sb, "  - [%s](#%s)\n", task.Title, task.Slug)
			}
		}
	}

	sb.WriteString("\n")
}

// isTaskComplete returns true if a task is considered complete.
func isTaskComplete(task tasks.Task) bool {
	if task.Status == tasks.StatusCompleted {
		return true
	}
	if len(task.Subtasks) > 0 {
		for _, subtask := range task.Subtasks {
			if !subtask.Completed {
				return false
			}
		}
		return true
	}
	return false
}

// countCompleted counts how many tasks in the slice are complete.
func countCompleted(taskList []tasks.Task) int {
	count := 0
	for _, task := range taskList {
		if isTaskComplete(task) {
			count++
		}
	}
	return count
}

func buildTOCEntries(tl *tasks.TaskList, opts Options) []tocEntry {
	var entries []tocEntry

	switch opts.GroupBy {
	case GroupByArea:
		tasksByArea := tl.TasksByArea()
		for _, area := range tl.Areas {
			areaTasks := tasksByArea[area.ID]
			if len(areaTasks) == 0 {
				continue
			}
			entry := tocEntry{
				Title:     area.Name,
				Slug:      slugify(area.Name),
				Count:     len(areaTasks),
				Completed: countCompleted(areaTasks),
			}
			for i, task := range sortTasks(areaTasks, opts) {
				title := task.Title
				if opts.NumberItems {
					title = fmt.Sprintf("%d. %s", i+1, task.Title)
				}
				entry.Tasks = append(entry.Tasks, tocEntry{
					Title: title,
					Slug:  taskSlug(task),
				})
			}
			entries = append(entries, entry)
		}

	case GroupByStatus:
		tasksByStatus := tl.TasksByStatus()
		legend := tl.GetLegend()
		for _, status := range tasks.StatusOrder() {
			statusTasks := tasksByStatus[status]
			if len(statusTasks) == 0 {
				continue
			}
			if status == tasks.StatusCompleted && !opts.ShowCompleted {
				continue
			}
			title := legend[status].Description
			entry := tocEntry{
				Title:     title,
				Slug:      slugify(title),
				Count:     len(statusTasks),
				Completed: countCompleted(statusTasks),
			}
			for i, task := range sortTasks(statusTasks, opts) {
				taskTitle := task.Title
				if opts.NumberItems {
					taskTitle = fmt.Sprintf("%d. %s", i+1, task.Title)
				}
				entry.Tasks = append(entry.Tasks, tocEntry{
					Title: taskTitle,
					Slug:  taskSlug(task),
				})
			}
			entries = append(entries, entry)
		}

	case GroupByPhase:
		tasksByPhase := tl.TasksByPhase()
		phases := tl.PhaseNumbers()
		for _, phase := range phases {
			phaseTasks := tasksByPhase[phase]
			if len(phaseTasks) == 0 {
				continue
			}
			title := fmt.Sprintf("Phase %d", phase)
			entry := tocEntry{
				Title:     title,
				Slug:      slugify(title),
				Count:     len(phaseTasks),
				Completed: countCompleted(phaseTasks),
			}
			for i, task := range sortTasks(phaseTasks, opts) {
				taskTitle := task.Title
				if opts.NumberItems {
					taskTitle = fmt.Sprintf("%d. %s", i+1, task.Title)
				}
				entry.Tasks = append(entry.Tasks, tocEntry{
					Title: taskTitle,
					Slug:  taskSlug(task),
				})
			}
			entries = append(entries, entry)
		}
		// Unphased tasks (phase 0)
		if phaseTasks := tasksByPhase[0]; len(phaseTasks) > 0 {
			entry := tocEntry{
				Title:     "Unphased",
				Slug:      "unphased",
				Count:     len(phaseTasks),
				Completed: countCompleted(phaseTasks),
			}
			entries = append(entries, entry)
		}

	case GroupByType:
		tasksByType := tl.TasksByType()
		registry := changelog.DefaultRegistry
		allTypes := registry.All()
		for _, ct := range allTypes {
			typeTasks := tasksByType[ct.Name]
			if len(typeTasks) == 0 {
				continue
			}
			entry := tocEntry{
				Title:     ct.Name,
				Slug:      slugify(ct.Name),
				Count:     len(typeTasks),
				Completed: countCompleted(typeTasks),
			}
			for i, task := range sortTasks(typeTasks, opts) {
				taskTitle := task.Title
				if opts.NumberItems {
					taskTitle = fmt.Sprintf("%d. %s", i+1, task.Title)
				}
				entry.Tasks = append(entry.Tasks, tocEntry{
					Title: taskTitle,
					Slug:  taskSlug(task),
				})
			}
			entries = append(entries, entry)
		}
	}

	return entries
}

// sortTasks returns a sorted copy of tasks for consistent ordering.
// Completed tasks are placed at the bottom within their group.
func sortTasks(taskList []tasks.Task, _ Options) []tasks.Task {
	sorted := make([]tasks.Task, len(taskList))
	copy(sorted, taskList)
	sort.Slice(sorted, func(i, j int) bool {
		// Within same phase, sort by status (completed at bottom)
		iOrder := statusSortOrder(sorted[i].Status)
		jOrder := statusSortOrder(sorted[j].Status)
		if iOrder != jOrder {
			return iOrder < jOrder
		}
		return sorted[i].Title < sorted[j].Title
	})
	return sorted
}

func renderByArea(sb *strings.Builder, tl *tasks.TaskList, opts Options) {
	tasksByArea := tl.TasksByArea()

	for _, area := range tl.Areas {
		areaTasks := tasksByArea[area.ID]
		if len(areaTasks) == 0 {
			continue
		}

		renderSectionHeading(sb, area.Name, tl.Project, opts)
		renderTasks(sb, areaTasks, tl, opts)

		if opts.HorizontalRules {
			sb.WriteString("---\n\n")
		}
	}

	// Unspecified area tasks
	if areaTasks, ok := tasksByArea["_unspecified"]; ok && len(areaTasks) > 0 {
		renderSectionHeading(sb, "Other", tl.Project, opts)
		renderTasks(sb, areaTasks, tl, opts)
	}
}

func renderByType(sb *strings.Builder, tl *tasks.TaskList, opts Options) {
	tasksByType := tl.TasksByType()

	registry := changelog.DefaultRegistry
	allTypes := registry.All()

	for _, ct := range allTypes {
		typeTasks := tasksByType[ct.Name]
		if len(typeTasks) == 0 {
			continue
		}

		renderSectionHeading(sb, ct.Name, tl.Project, opts)
		renderTasks(sb, typeTasks, tl, opts)

		if opts.HorizontalRules {
			sb.WriteString("---\n\n")
		}
	}

	// Unspecified type tasks
	if typeTasks, ok := tasksByType["_unspecified"]; ok && len(typeTasks) > 0 {
		renderSectionHeading(sb, "Other", tl.Project, opts)
		renderTasks(sb, typeTasks, tl, opts)
	}
}

func renderByPhase(sb *strings.Builder, tl *tasks.TaskList, opts Options) {
	tasksByPhase := tl.TasksByPhase()
	phases := tl.PhaseNumbers()

	// Build area name lookup
	areaNames := make(map[string]string)
	for _, area := range tl.Areas {
		areaNames[area.ID] = area.Name
	}

	for _, phase := range phases {
		phaseTasks := tasksByPhase[phase]
		if len(phaseTasks) == 0 {
			continue
		}

		header := fmt.Sprintf("Phase %d", phase)
		renderSectionHeading(sb, header, tl.Project, opts)

		if opts.ShowAreaSubheadings && len(tl.Areas) > 0 {
			renderTasksByAreaWithinPhase(sb, phaseTasks, tl, opts, areaNames)
		} else {
			renderTasks(sb, phaseTasks, tl, opts)
		}

		if opts.HorizontalRules {
			sb.WriteString("---\n\n")
		}
	}

	// Unphased tasks (phase 0)
	if phaseTasks := tasksByPhase[0]; len(phaseTasks) > 0 {
		renderSectionHeading(sb, "Unphased", tl.Project, opts)
		if opts.ShowAreaSubheadings && len(tl.Areas) > 0 {
			renderTasksByAreaWithinPhase(sb, phaseTasks, tl, opts, areaNames)
		} else {
			renderTasks(sb, phaseTasks, tl, opts)
		}
	}
}

// renderTasksByAreaWithinPhase renders tasks grouped by area as sub-sections.
func renderTasksByAreaWithinPhase(sb *strings.Builder, taskList []tasks.Task, tl *tasks.TaskList, opts Options, areaNames map[string]string) {
	// Group tasks by area
	tasksByArea := make(map[string][]tasks.Task)
	for _, task := range taskList {
		areaID := task.Area
		if areaID == "" {
			areaID = "_unspecified"
		}
		tasksByArea[areaID] = append(tasksByArea[areaID], task)
	}

	// Get sorted area IDs
	areaIDs := make([]string, 0, len(tasksByArea))
	for areaID := range tasksByArea {
		areaIDs = append(areaIDs, areaID)
	}
	sort.Strings(areaIDs)

	// Render each area as a sub-section
	for _, areaID := range areaIDs {
		areaTasks := tasksByArea[areaID]
		if len(areaTasks) == 0 {
			continue
		}

		areaName := areaNames[areaID]
		if areaName == "" {
			if areaID == "_unspecified" {
				areaName = "Other"
			} else {
				areaName = areaID
			}
		}
		fmt.Fprintf(sb, "### %s\n\n", areaName)

		renderTasksAsList(sb, areaTasks, tl, opts)
	}
}

// renderTasksAsList renders tasks as a simple list with checkboxes.
func renderTasksAsList(sb *strings.Builder, taskList []tasks.Task, _ *tasks.TaskList, opts Options) {
	sorted := sortTasks(taskList, opts)

	for _, task := range sorted {
		if task.Status == tasks.StatusCompleted && !opts.ShowCompleted {
			continue
		}

		isComplete := isTaskComplete(task)

		var line string
		if opts.UseCheckboxes {
			checkbox := "[ ]"
			if isComplete {
				checkbox = "[x]"
			}
			line = fmt.Sprintf("- %s %s", checkbox, task.Title)
		} else {
			line = fmt.Sprintf("- %s", task.Title)
		}

		if task.Description != "" {
			line += " - " + task.Description
		}

		sb.WriteString(line + "\n")

		// Render subtasks
		for _, subtask := range task.Subtasks {
			subtaskCheckbox := "[ ]"
			if subtask.Completed {
				subtaskCheckbox = "[x]"
			}
			fmt.Fprintf(sb, "  - %s %s\n", subtaskCheckbox, subtask.Description)
		}
	}
	sb.WriteString("\n")
}

func renderByStatus(sb *strings.Builder, tl *tasks.TaskList, opts Options) {
	tasksByStatus := tl.TasksByStatus()

	for _, status := range tasks.StatusOrder() {
		statusTasks := tasksByStatus[status]
		if len(statusTasks) == 0 {
			continue
		}
		if status == tasks.StatusCompleted && !opts.ShowCompleted {
			continue
		}

		legend := tl.GetLegend()
		header := legend[status].Description
		if opts.UseEmoji {
			header = legend[status].Emoji + " " + header
		}
		renderSectionHeading(sb, header, tl.Project, opts)
		renderTasks(sb, statusTasks, tl, opts)

		if opts.HorizontalRules {
			sb.WriteString("---\n\n")
		}
	}
}

func renderTasks(sb *strings.Builder, taskList []tasks.Task, tl *tasks.TaskList, opts Options) {
	sorted := sortTasks(taskList, opts)

	for i, task := range sorted {
		if task.Status == tasks.StatusCompleted && !opts.ShowCompleted {
			continue
		}
		renderTask(sb, task, i+1, tl, opts)
	}
}

func renderTask(sb *strings.Builder, task tasks.Task, num int, tl *tasks.TaskList, opts Options) {
	isComplete := isTaskComplete(task)

	// Task header with checkbox
	var title string
	if opts.UseCheckboxes {
		checkbox := "[ ]"
		if isComplete {
			checkbox = "[x]"
		}
		if opts.NumberItems {
			title = fmt.Sprintf("%d. %s %s", num, checkbox, task.Title)
		} else {
			title = fmt.Sprintf("%s %s", checkbox, task.Title)
		}
	} else {
		if opts.NumberItems {
			title = fmt.Sprintf("%d. %s", num, task.Title)
		} else {
			title = task.Title
		}
	}

	// Add emoji suffix if not using checkboxes
	if opts.UseEmoji && !opts.UseCheckboxes {
		title += " " + tl.GetStatusEmoji(task.Status)
	}

	// Add stable anchor for navigation
	fmt.Fprintf(sb, "<a id=\"%s\"></a>\n\n", taskSlug(task))
	fmt.Fprintf(sb, "### %s\n\n", title)

	// Description
	if task.Description != "" {
		sb.WriteString(task.Description + "\n\n")
	}

	// Subtasks
	if len(task.Subtasks) > 0 {
		for _, subtask := range task.Subtasks {
			checkbox := "[ ]"
			if subtask.Completed {
				checkbox = "[x]"
			}
			if opts.UseCheckboxes {
				fmt.Fprintf(sb, "- %s %s\n", checkbox, subtask.Description)
			} else {
				prefix := "- "
				if subtask.Completed {
					prefix = "- ✅ "
				}
				sb.WriteString(prefix + subtask.Description + "\n")
			}
		}
		sb.WriteString("\n")
	}
}
