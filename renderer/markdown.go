package renderer

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/grokify/structured-changelog/changelog"
	"github.com/grokify/structured-roadmap/roadmap"
)

// slugify converts a heading to a GitHub-flavored markdown anchor.
func slugify(s string) string {
	// Remove emoji and special characters, lowercase, replace spaces with hyphens
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	// Remove characters that aren't alphanumeric, hyphens, or underscores
	re := regexp.MustCompile(`[^a-z0-9\-_]`)
	s = re.ReplaceAllString(s, "")
	// Collapse multiple hyphens
	re = regexp.MustCompile(`-+`)
	s = re.ReplaceAllString(s, "-")
	// Trim leading/trailing hyphens
	s = strings.Trim(s, "-")
	return s
}

// tocEntry represents an entry in the table of contents.
type tocEntry struct {
	Title     string
	Slug      string
	Count     int
	Completed int
	Items     []tocEntry
}

// Render generates Markdown from a Roadmap.
func Render(r *roadmap.Roadmap, opts Options) string {
	var sb strings.Builder

	// Title
	sb.WriteString(fmt.Sprintf("# %s Roadmap\n\n", r.Project))

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
		renderOverviewTable(&sb, r, opts)
		if opts.HorizontalRules {
			sb.WriteString("---\n\n")
		}
	}

	// Overview section (if exists)
	for _, section := range r.Sections {
		if section.ID == "overview" {
			renderSection(&sb, section, r, opts)
			if opts.HorizontalRules {
				sb.WriteString("---\n\n")
			}
			break
		}
	}

	// Table of Contents
	if opts.ShowTOC {
		renderTOC(&sb, r, opts)
		if opts.HorizontalRules {
			sb.WriteString("---\n\n")
		}
	}

	// Legend
	if opts.ShowLegend {
		renderLegend(&sb, r)
		if opts.HorizontalRules {
			sb.WriteString("---\n\n")
		}
	}

	// Main content grouped by strategy
	switch opts.GroupBy {
	case GroupByPhase:
		renderByPhase(&sb, r, opts)
	case GroupByStatus:
		renderByStatus(&sb, r, opts)
	case GroupByQuarter:
		renderByQuarter(&sb, r, opts)
	case GroupByPriority:
		renderByPriority(&sb, r, opts)
	case GroupByType:
		renderByType(&sb, r, opts)
	default:
		renderByArea(&sb, r, opts)
	}

	// Additional sections
	if opts.ShowSections {
		renderSections(&sb, r, opts)
	}

	// Version history
	if opts.ShowVersionHistory && len(r.VersionHistory) > 0 {
		if opts.HorizontalRules {
			sb.WriteString("---\n\n")
		}
		renderVersionHistory(&sb, r, opts)
	}

	// Dependencies
	if opts.ShowDependencies && r.Dependencies != nil {
		if opts.HorizontalRules {
			sb.WriteString("---\n\n")
		}
		renderDependencies(&sb, r, opts)
	}

	return strings.TrimRight(sb.String(), "\n") + "\n"
}

// RenderToFile writes rendered Markdown to a file.
func RenderToFile(path string, r *roadmap.Roadmap, opts Options) error {
	content := Render(r, opts)
	return os.WriteFile(path, []byte(content), 0600)
}

func renderLegend(sb *strings.Builder, r *roadmap.Roadmap) {
	sb.WriteString("## Legend\n\n")
	sb.WriteString("| Status | Description |\n")
	sb.WriteString("|--------|-------------|\n")
	legend := r.GetLegend()
	for _, status := range []roadmap.Status{roadmap.StatusCompleted, roadmap.StatusInProgress, roadmap.StatusPlanned, roadmap.StatusFuture} {
		if entry, ok := legend[status]; ok {
			fmt.Fprintf(sb, "| %s | %s |\n", entry.Emoji, entry.Description)
		}
	}
	sb.WriteString("\n")
}

func renderOverviewTable(sb *strings.Builder, r *roadmap.Roadmap, opts Options) {
	sb.WriteString("## Overview\n\n")
	sb.WriteString("| Item | Status | Priority | Area |\n")
	sb.WriteString("|------|--------|----------|------|\n")

	legend := r.GetLegend()

	// Build area name lookup
	areaNames := make(map[string]string)
	for _, area := range r.Areas {
		areaNames[area.ID] = area.Name
	}

	// Sort items by order, then priority
	sorted := make([]roadmap.Item, len(r.Items))
	copy(sorted, r.Items)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Order != sorted[j].Order {
			return sorted[i].Order < sorted[j].Order
		}
		pi := roadmap.PriorityOrder(sorted[i].Priority)
		pj := roadmap.PriorityOrder(sorted[j].Priority)
		return pi < pj
	})

	for _, item := range sorted {
		if item.Status == roadmap.StatusCompleted && !opts.ShowCompleted {
			continue
		}

		// Status emoji
		status := ""
		if opts.UseEmoji {
			if entry, ok := legend[item.Status]; ok {
				status = entry.Emoji
			}
		} else {
			status = string(item.Status)
		}

		// Priority label
		priority := roadmap.PriorityLabel(item.Priority)

		// Area name
		areaName := areaNames[item.Area]
		if areaName == "" {
			areaName = "-"
		}

		fmt.Fprintf(sb, "| %s | %s | %s | %s |\n", item.Title, status, priority, areaName)
	}
	sb.WriteString("\n")
}

func renderTOC(sb *strings.Builder, r *roadmap.Roadmap, opts Options) {
	sb.WriteString("## Table of Contents\n\n")

	entries := buildTOCEntries(r, opts)

	for _, entry := range entries {
		// Section level with progress (completed/total)
		fmt.Fprintf(sb, "- [%s (%d/%d)](#%s)\n", entry.Title, entry.Completed, entry.Count, entry.Slug)

		// Item level (depth 2)
		if opts.TOCDepth >= 2 {
			for _, item := range entry.Items {
				fmt.Fprintf(sb, "  - [%s](#%s)\n", item.Title, item.Slug)
			}
		}
	}

	// Add freeform sections to TOC
	if opts.ShowSections {
		for _, section := range r.Sections {
			if section.ID != "overview" {
				count := countSectionItems(section)
				if count > 0 {
					fmt.Fprintf(sb, "- [%s (%d)](#%s)\n", section.Title, count, slugify(section.Title))
				} else {
					fmt.Fprintf(sb, "- [%s](#%s)\n", section.Title, slugify(section.Title))
				}
			}
		}
	}

	sb.WriteString("\n")
}

// isItemComplete returns true if an item is considered complete.
// An item is complete if its status is "completed" or all its tasks are done.
func isItemComplete(item roadmap.Item) bool {
	if item.Status == roadmap.StatusCompleted {
		return true
	}
	if len(item.Tasks) > 0 {
		for _, task := range item.Tasks {
			if !task.Completed {
				return false
			}
		}
		return true
	}
	return false
}

// countCompleted counts how many items in the slice are complete.
func countCompleted(items []roadmap.Item) int {
	count := 0
	for _, item := range items {
		if isItemComplete(item) {
			count++
		}
	}
	return count
}

// countSectionItems counts list items in a section's content blocks.
func countSectionItems(section roadmap.Section) int {
	count := 0
	for _, block := range section.Content {
		if block.Type == roadmap.ContentTypeList {
			count += len(block.Items)
		}
	}
	return count
}

func buildTOCEntries(r *roadmap.Roadmap, opts Options) []tocEntry {
	var entries []tocEntry

	switch opts.GroupBy {
	case GroupByArea: //nolint:dupl // similar structure to GroupByPhase but different types
		// Sort areas by priority
		areas := make([]roadmap.Area, len(r.Areas))
		copy(areas, r.Areas)
		sort.Slice(areas, func(i, j int) bool {
			return areas[i].Priority < areas[j].Priority
		})

		itemsByArea := r.ItemsByArea()
		for _, area := range areas {
			items := itemsByArea[area.ID]
			if len(items) == 0 {
				continue
			}
			entry := tocEntry{
				Title:     area.Name,
				Slug:      slugify(area.Name),
				Count:     len(items),
				Completed: countCompleted(items),
			}
			for i, item := range sortItems(items, opts) {
				title := item.Title
				if opts.NumberItems {
					title = fmt.Sprintf("%d. %s", i+1, item.Title)
				}
				entry.Items = append(entry.Items, tocEntry{
					Title: title,
					Slug:  slugify(title),
				})
			}
			entries = append(entries, entry)
		}

	case GroupByPriority:
		priorityOrder := []roadmap.Priority{
			roadmap.PriorityCritical,
			roadmap.PriorityHigh,
			roadmap.PriorityMedium,
			roadmap.PriorityLow,
		}
		itemsByPriority := r.ItemsByPriority()
		for _, priority := range priorityOrder {
			items := itemsByPriority[priority]
			if len(items) == 0 {
				continue
			}
			title := roadmap.PriorityLabelFull(priority)
			entry := tocEntry{
				Title:     title,
				Slug:      slugify(title),
				Count:     len(items),
				Completed: countCompleted(items),
			}
			for i, item := range sortItems(items, opts) {
				itemTitle := item.Title
				if opts.NumberItems {
					itemTitle = fmt.Sprintf("%d. %s", i+1, item.Title)
				}
				entry.Items = append(entry.Items, tocEntry{
					Title: itemTitle,
					Slug:  slugify(itemTitle),
				})
			}
			entries = append(entries, entry)
		}

	case GroupByStatus:
		statusOrder := []roadmap.Status{roadmap.StatusCompleted, roadmap.StatusInProgress, roadmap.StatusPlanned, roadmap.StatusFuture}
		itemsByStatus := r.ItemsByStatus()
		legend := r.GetLegend()
		for _, status := range statusOrder {
			items := itemsByStatus[status]
			if len(items) == 0 {
				continue
			}
			if status == roadmap.StatusCompleted && !opts.ShowCompleted {
				continue
			}
			title := legend[status].Description
			entry := tocEntry{
				Title:     title,
				Slug:      slugify(title),
				Count:     len(items),
				Completed: countCompleted(items),
			}
			for i, item := range sortItems(items, opts) {
				itemTitle := item.Title
				if opts.NumberItems {
					itemTitle = fmt.Sprintf("%d. %s", i+1, item.Title)
				}
				entry.Items = append(entry.Items, tocEntry{
					Title: itemTitle,
					Slug:  slugify(itemTitle),
				})
			}
			entries = append(entries, entry)
		}

	case GroupByPhase: //nolint:dupl // similar structure to GroupByArea but different types
		// Sort phases by order
		phases := make([]roadmap.Phase, len(r.Phases))
		copy(phases, r.Phases)
		sort.Slice(phases, func(i, j int) bool {
			return phases[i].Order < phases[j].Order
		})

		itemsByPhase := r.ItemsByPhase()
		for _, phase := range phases {
			items := itemsByPhase[phase.ID]
			if len(items) == 0 {
				continue
			}
			entry := tocEntry{
				Title:     phase.Name,
				Slug:      slugify(phase.Name),
				Count:     len(items),
				Completed: countCompleted(items),
			}
			for i, item := range sortItems(items, opts) {
				itemTitle := item.Title
				if opts.NumberItems {
					itemTitle = fmt.Sprintf("%d. %s", i+1, item.Title)
				}
				entry.Items = append(entry.Items, tocEntry{
					Title: itemTitle,
					Slug:  slugify(itemTitle),
				})
			}
			entries = append(entries, entry)
		}

	case GroupByQuarter:
		itemsByQuarter := r.ItemsByQuarter()
		quarters := make([]string, 0, len(itemsByQuarter))
		for q := range itemsByQuarter {
			quarters = append(quarters, q)
		}
		sort.Strings(quarters)

		for _, quarter := range quarters {
			items := itemsByQuarter[quarter]
			if len(items) == 0 {
				continue
			}
			title := quarter
			if quarter == "_unscheduled" {
				title = "Unscheduled"
			}
			entry := tocEntry{
				Title:     title,
				Slug:      slugify(title),
				Count:     len(items),
				Completed: countCompleted(items),
			}
			for i, item := range sortItems(items, opts) {
				itemTitle := item.Title
				if opts.NumberItems {
					itemTitle = fmt.Sprintf("%d. %s", i+1, item.Title)
				}
				entry.Items = append(entry.Items, tocEntry{
					Title: itemTitle,
					Slug:  slugify(itemTitle),
				})
			}
			entries = append(entries, entry)
		}

	case GroupByType:
		itemsByType := r.ItemsByType()
		registry := changelog.DefaultRegistry
		allTypes := registry.All()

		for _, ct := range allTypes {
			items := itemsByType[ct.Name]
			if len(items) == 0 {
				continue
			}
			entry := tocEntry{
				Title:     ct.Name,
				Slug:      slugify(ct.Name),
				Count:     len(items),
				Completed: countCompleted(items),
			}
			for i, item := range sortItems(items, opts) {
				itemTitle := item.Title
				if opts.NumberItems {
					itemTitle = fmt.Sprintf("%d. %s", i+1, item.Title)
				}
				entry.Items = append(entry.Items, tocEntry{
					Title: itemTitle,
					Slug:  slugify(itemTitle),
				})
			}
			entries = append(entries, entry)
		}
	}

	return entries
}

// sortItems returns a sorted copy of items for consistent ordering.
func sortItems(items []roadmap.Item, _ Options) []roadmap.Item {
	sorted := make([]roadmap.Item, len(items))
	copy(sorted, items)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Order != sorted[j].Order {
			return sorted[i].Order < sorted[j].Order
		}
		pi := roadmap.PriorityOrder(sorted[i].Priority)
		pj := roadmap.PriorityOrder(sorted[j].Priority)
		if pi != pj {
			return pi < pj
		}
		return sorted[i].Title < sorted[j].Title
	})
	return sorted
}

func renderByArea(sb *strings.Builder, r *roadmap.Roadmap, opts Options) {
	// Sort areas by priority
	areas := make([]roadmap.Area, len(r.Areas))
	copy(areas, r.Areas)
	sort.Slice(areas, func(i, j int) bool {
		return areas[i].Priority < areas[j].Priority
	})

	itemsByArea := r.ItemsByArea()

	for _, area := range areas {
		items := itemsByArea[area.ID]
		if len(items) == 0 {
			continue
		}

		fmt.Fprintf(sb, "## %s\n\n", area.Name)
		renderItems(sb, items, r, opts)

		if opts.HorizontalRules {
			sb.WriteString("---\n\n")
		}
	}

	// Unspecified area items
	if items, ok := itemsByArea["_unspecified"]; ok && len(items) > 0 {
		sb.WriteString("## Other\n\n")
		renderItems(sb, items, r, opts)
	}
}

func renderByType(sb *strings.Builder, r *roadmap.Roadmap, opts Options) {
	itemsByType := r.ItemsByType()

	// Get types in canonical order from structured-changelog registry
	registry := changelog.DefaultRegistry
	allTypes := registry.All()

	for _, ct := range allTypes {
		items := itemsByType[ct.Name]
		if len(items) == 0 {
			continue
		}

		fmt.Fprintf(sb, "## %s\n\n", ct.Name)
		renderItems(sb, items, r, opts)

		if opts.HorizontalRules {
			sb.WriteString("---\n\n")
		}
	}

	// Unspecified type items
	if items, ok := itemsByType["_unspecified"]; ok && len(items) > 0 {
		sb.WriteString("## Other\n\n")
		renderItems(sb, items, r, opts)
	}
}

func renderByPhase(sb *strings.Builder, r *roadmap.Roadmap, opts Options) {
	// Sort phases by order
	phases := make([]roadmap.Phase, len(r.Phases))
	copy(phases, r.Phases)
	sort.Slice(phases, func(i, j int) bool {
		return phases[i].Order < phases[j].Order
	})

	// Build area name lookup and sorted area list
	areaNames := make(map[string]string)
	areaOrder := make(map[string]int)
	for _, area := range r.Areas {
		areaNames[area.ID] = area.Name
		areaOrder[area.ID] = area.Priority
	}

	itemsByPhase := r.ItemsByPhase()

	for _, phase := range phases {
		items := itemsByPhase[phase.ID]

		// Phase header with status
		header := phase.Name
		if opts.UseEmoji && phase.Status != "" {
			header += " " + r.GetStatusEmoji(phase.Status)
		}
		fmt.Fprintf(sb, "## %s\n\n", header)

		if phase.Description != "" {
			sb.WriteString(phase.Description + "\n\n")
		}

		if len(items) > 0 {
			if opts.ShowAreaSubheadings && len(r.Areas) > 0 {
				// Group items by area within this phase
				renderItemsByAreaWithinPhase(sb, items, r, opts, areaNames, areaOrder)
			} else {
				renderItems(sb, items, r, opts)
			}
		}

		if opts.HorizontalRules {
			sb.WriteString("---\n\n")
		}
	}

	// Unphased items
	if items, ok := itemsByPhase["_unphased"]; ok && len(items) > 0 {
		sb.WriteString("## Other\n\n")
		if opts.ShowAreaSubheadings && len(r.Areas) > 0 {
			renderItemsByAreaWithinPhase(sb, items, r, opts, areaNames, areaOrder)
		} else {
			renderItems(sb, items, r, opts)
		}
	}
}

// renderItemsByAreaWithinPhase renders items grouped by area as sub-sections within a phase.
func renderItemsByAreaWithinPhase(sb *strings.Builder, items []roadmap.Item, r *roadmap.Roadmap, opts Options, areaNames map[string]string, areaOrder map[string]int) {
	// Group items by area
	itemsByArea := make(map[string][]roadmap.Item)
	for _, item := range items {
		areaID := item.Area
		if areaID == "" {
			areaID = "_unspecified"
		}
		itemsByArea[areaID] = append(itemsByArea[areaID], item)
	}

	// Get sorted area IDs (by priority)
	areaIDs := make([]string, 0, len(itemsByArea))
	for areaID := range itemsByArea {
		areaIDs = append(areaIDs, areaID)
	}
	sort.Slice(areaIDs, func(i, j int) bool {
		return areaOrder[areaIDs[i]] < areaOrder[areaIDs[j]]
	})

	// Render each area as a sub-section
	for _, areaID := range areaIDs {
		areaItems := itemsByArea[areaID]
		if len(areaItems) == 0 {
			continue
		}

		// Area sub-heading
		areaName := areaNames[areaID]
		if areaName == "" {
			if areaID == "_unspecified" {
				areaName = "Other"
			} else {
				areaName = areaID
			}
		}
		fmt.Fprintf(sb, "### %s\n\n", areaName)

		// Render items as task list (simpler format for sub-sections)
		renderItemsAsTasks(sb, areaItems, r, opts)
	}
}

// renderItemsAsTasks renders items as a simple task list with checkboxes.
// This is used for sub-sections where full item headers would be too verbose.
func renderItemsAsTasks(sb *strings.Builder, items []roadmap.Item, _ *roadmap.Roadmap, opts Options) {
	// Sort items
	sorted := sortItems(items, opts)

	for _, item := range sorted {
		if item.Status == roadmap.StatusCompleted && !opts.ShowCompleted {
			continue
		}

		isComplete := isItemComplete(item)

		// Build task line
		var line string
		if opts.UseCheckboxes {
			checkbox := "[ ]"
			if isComplete {
				checkbox = "[x]"
			}
			line = fmt.Sprintf("- %s %s", checkbox, item.Title)
		} else {
			line = fmt.Sprintf("- %s", item.Title)
		}

		// Add description if present
		if item.Description != "" {
			line += " - " + item.Description
		}

		sb.WriteString(line + "\n")

		// Render sub-tasks if present
		for _, task := range item.Tasks {
			taskCheckbox := "[ ]"
			if task.Completed {
				taskCheckbox = "[x]"
			}
			taskLine := fmt.Sprintf("  - %s %s", taskCheckbox, task.Description)
			if task.FilePath != "" {
				taskLine += fmt.Sprintf(" (`%s`)", task.FilePath)
			}
			sb.WriteString(taskLine + "\n")
		}
	}
	sb.WriteString("\n")
}

func renderByStatus(sb *strings.Builder, r *roadmap.Roadmap, opts Options) {
	statusOrder := []roadmap.Status{roadmap.StatusCompleted, roadmap.StatusInProgress, roadmap.StatusPlanned, roadmap.StatusFuture}
	itemsByStatus := r.ItemsByStatus()

	for _, status := range statusOrder {
		items := itemsByStatus[status]
		if len(items) == 0 {
			continue
		}
		if status == roadmap.StatusCompleted && !opts.ShowCompleted {
			continue
		}

		legend := r.GetLegend()
		header := legend[status].Description
		if opts.UseEmoji {
			header = legend[status].Emoji + " " + header
		}
		fmt.Fprintf(sb, "## %s\n\n", header)
		renderItems(sb, items, r, opts)

		if opts.HorizontalRules {
			sb.WriteString("---\n\n")
		}
	}
}

func renderByQuarter(sb *strings.Builder, r *roadmap.Roadmap, opts Options) {
	itemsByQuarter := r.ItemsByQuarter()

	// Get sorted quarters
	quarters := make([]string, 0, len(itemsByQuarter))
	for q := range itemsByQuarter {
		quarters = append(quarters, q)
	}
	sort.Strings(quarters)

	for _, quarter := range quarters {
		items := itemsByQuarter[quarter]
		if len(items) == 0 {
			continue
		}

		header := quarter
		if quarter == "_unscheduled" {
			header = "Unscheduled"
		}
		fmt.Fprintf(sb, "## %s\n\n", header)
		renderItems(sb, items, r, opts)

		if opts.HorizontalRules {
			sb.WriteString("---\n\n")
		}
	}
}

func renderByPriority(sb *strings.Builder, r *roadmap.Roadmap, opts Options) {
	priorityOrder := []roadmap.Priority{
		roadmap.PriorityCritical,
		roadmap.PriorityHigh,
		roadmap.PriorityMedium,
		roadmap.PriorityLow,
	}
	itemsByPriority := r.ItemsByPriority()

	for _, priority := range priorityOrder {
		items := itemsByPriority[priority]
		if len(items) == 0 {
			continue
		}

		header := roadmap.PriorityLabelFull(priority)
		fmt.Fprintf(sb, "## %s\n\n", header)
		renderItems(sb, items, r, opts)

		if opts.HorizontalRules {
			sb.WriteString("---\n\n")
		}
	}

	// Unspecified priority items
	if items, ok := itemsByPriority["_unspecified"]; ok && len(items) > 0 {
		sb.WriteString("## Other\n\n")
		renderItems(sb, items, r, opts)
	}
}

func renderItems(sb *strings.Builder, items []roadmap.Item, r *roadmap.Roadmap, opts Options) {
	// Sort items by order, then by priority, then by title
	sorted := make([]roadmap.Item, len(items))
	copy(sorted, items)
	sort.Slice(sorted, func(i, j int) bool {
		// First by explicit order
		if sorted[i].Order != sorted[j].Order {
			return sorted[i].Order < sorted[j].Order
		}
		// Then by priority level
		pi := roadmap.PriorityOrder(sorted[i].Priority)
		pj := roadmap.PriorityOrder(sorted[j].Priority)
		if pi != pj {
			return pi < pj
		}
		// Then alphabetically
		return sorted[i].Title < sorted[j].Title
	})

	for i, item := range sorted {
		if item.Status == roadmap.StatusCompleted && !opts.ShowCompleted {
			continue
		}
		renderItem(sb, item, i+1, r, opts)
	}
}

func renderItem(sb *strings.Builder, item roadmap.Item, num int, r *roadmap.Roadmap, opts Options) {
	// Determine if item is "done" - either completed status or all tasks done
	isComplete := isItemComplete(item)

	// Item header with checkbox
	var title string
	if opts.UseCheckboxes {
		checkbox := "[ ]"
		if isComplete {
			checkbox = "[x]"
		}
		if opts.NumberItems {
			title = fmt.Sprintf("%d. %s %s", num, checkbox, item.Title)
		} else {
			title = fmt.Sprintf("%s %s", checkbox, item.Title)
		}
	} else {
		if opts.NumberItems {
			title = fmt.Sprintf("%d. %s", num, item.Title)
		} else {
			title = item.Title
		}
	}
	// Only add emoji suffix if not using checkboxes (avoid redundancy)
	if opts.UseEmoji && !opts.UseCheckboxes {
		title += " " + r.GetStatusEmoji(item.Status)
	}
	fmt.Fprintf(sb, "### %s\n\n", title)

	// Description
	if item.Description != "" {
		sb.WriteString(item.Description + "\n\n")
	}

	// Version/date info
	if item.Version != "" {
		fmt.Fprintf(sb, "**Version:** %s", item.Version)
		if item.CompletedDate != "" {
			fmt.Fprintf(sb, " (%s)", item.CompletedDate)
		}
		sb.WriteString("\n\n")
	} else if item.TargetVersion != "" {
		fmt.Fprintf(sb, "**Target:** %s", item.TargetVersion)
		if item.TargetQuarter != "" {
			fmt.Fprintf(sb, " (%s)", item.TargetQuarter)
		}
		sb.WriteString("\n\n")
	} else if item.TargetQuarter != "" {
		fmt.Fprintf(sb, "**Target:** %s\n\n", item.TargetQuarter)
	}

	// Tasks
	if len(item.Tasks) > 0 {
		for _, task := range item.Tasks {
			checkbox := "[ ]"
			if task.Completed {
				checkbox = "[x]"
			}
			if opts.UseCheckboxes {
				line := fmt.Sprintf("- %s %s", checkbox, task.Description)
				if task.FilePath != "" {
					line += fmt.Sprintf(" (`%s`)", task.FilePath)
				}
				sb.WriteString(line + "\n")
			} else {
				prefix := "- "
				if task.Completed {
					prefix = "- ✅ "
				}
				sb.WriteString(prefix + task.Description + "\n")
			}
		}
		sb.WriteString("\n")
	}

	// Content blocks
	for _, block := range item.Content {
		renderContentBlock(sb, block)
	}
}

func renderContentBlock(sb *strings.Builder, block roadmap.ContentBlock) {
	switch block.Type {
	case roadmap.ContentTypeText:
		sb.WriteString(block.Value + "\n\n")

	case roadmap.ContentTypeCode:
		lang := block.Language
		if lang == "" {
			lang = ""
		}
		fmt.Fprintf(sb, "```%s\n%s\n```\n\n", lang, block.Value)

	case roadmap.ContentTypeDiagram:
		sb.WriteString("```\n" + block.Value + "\n```\n\n")

	case roadmap.ContentTypeTable:
		if len(block.Headers) > 0 {
			sb.WriteString("| " + strings.Join(block.Headers, " | ") + " |\n")
			sb.WriteString("|" + strings.Repeat("--------|", len(block.Headers)) + "\n")
			for _, row := range block.Rows {
				sb.WriteString("| " + strings.Join(row, " | ") + " |\n")
			}
			sb.WriteString("\n")
		}

	case roadmap.ContentTypeList:
		for _, item := range block.Items {
			sb.WriteString("- " + item + "\n")
		}
		sb.WriteString("\n")

	case roadmap.ContentTypeBlockquote:
		// Prefix each line with "> "
		lines := strings.Split(block.Value, "\n")
		for _, line := range lines {
			sb.WriteString("> " + line + "\n")
		}
		sb.WriteString("\n")
	}
}

func renderSections(sb *strings.Builder, r *roadmap.Roadmap, opts Options) {
	// Filter out overview (already rendered)
	sections := make([]roadmap.Section, 0)
	for _, s := range r.Sections {
		if s.ID != "overview" {
			sections = append(sections, s)
		}
	}

	// Sort by order
	sort.Slice(sections, func(i, j int) bool {
		return sections[i].Order < sections[j].Order
	})

	for _, section := range sections {
		renderSection(sb, section, r, opts)
	}
}

func renderSection(sb *strings.Builder, section roadmap.Section, _ *roadmap.Roadmap, _ Options) {
	fmt.Fprintf(sb, "## %s\n\n", section.Title)
	for _, block := range section.Content {
		renderContentBlock(sb, block)
	}
}

func renderVersionHistory(sb *strings.Builder, r *roadmap.Roadmap, opts Options) {
	sb.WriteString("## Version History\n\n")
	sb.WriteString("| Version | Date | Status | Summary |\n")
	sb.WriteString("|---------|------|--------|--------|\n")
	for _, v := range r.VersionHistory {
		date := v.Date
		if date == "" {
			date = "TBD"
		}
		status := ""
		if opts.UseEmoji && v.Status != "" {
			status = r.GetStatusEmoji(v.Status)
		} else if v.Status != "" {
			status = string(v.Status)
		}
		fmt.Fprintf(sb, "| %s | %s | %s | %s |\n", v.Version, date, status, v.Summary)
	}
	sb.WriteString("\n")
}

func renderDependencies(sb *strings.Builder, r *roadmap.Roadmap, _ Options) {
	if r.Dependencies == nil {
		return
	}

	sb.WriteString("## Dependencies\n\n")

	if len(r.Dependencies.External) > 0 {
		sb.WriteString("### External\n\n")
		sb.WriteString("| Name | Status | Note |\n")
		sb.WriteString("|------|--------|------|\n")
		for _, dep := range r.Dependencies.External {
			fmt.Fprintf(sb, "| %s | %s | %s |\n", dep.Name, dep.Status, dep.Note)
		}
		sb.WriteString("\n")
	}

	if len(r.Dependencies.Internal) > 0 {
		sb.WriteString("### Internal\n\n")
		sb.WriteString("| Package | Depends On |\n")
		sb.WriteString("|---------|------------|\n")
		for _, dep := range r.Dependencies.Internal {
			deps := strings.Join(dep.DependsOn, ", ")
			fmt.Fprintf(sb, "| %s | %s |\n", dep.Package, deps)
		}
		sb.WriteString("\n")
	}
}
