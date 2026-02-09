package main

import (
	"fmt"
	"os"

	"github.com/grokify/structured-tasks/renderer"
	"github.com/grokify/structured-tasks/tasks"
	"github.com/spf13/cobra"
)

var (
	genInput           string
	genOutput          string
	genGroupBy         string
	genCheckbox        bool
	genEmoji           bool
	genLegend          bool
	genNoIntro         bool
	genTOC             bool
	genTOCDepth        int
	genOverview        bool
	genAreaSubheadings bool
	genNumbered        bool
	genNoRules         bool
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate TASKS.md from TASKS.json",
	Long:  `Generate a Markdown task list file from a JSON intermediate representation.`,
	RunE:  runGenerate,
}

func init() {
	generateCmd.Flags().StringVarP(&genInput, "input", "i", "TASKS.json", "Input JSON file")
	generateCmd.Flags().StringVarP(&genOutput, "output", "o", "", "Output Markdown file (default: stdout)")
	generateCmd.Flags().StringVar(&genGroupBy, "group-by", "area", "Grouping: area, type, phase, status")
	generateCmd.Flags().BoolVar(&genCheckbox, "checkboxes", true, "Use [x]/[ ] checkbox syntax")
	generateCmd.Flags().BoolVar(&genEmoji, "emoji", true, "Include emoji status indicators")
	generateCmd.Flags().BoolVar(&genLegend, "legend", false, "Show legend table")
	generateCmd.Flags().BoolVar(&genNoIntro, "no-intro", false, "Omit introductory paragraph")
	generateCmd.Flags().BoolVar(&genTOC, "toc", false, "Show table of contents")
	generateCmd.Flags().IntVar(&genTOCDepth, "toc-depth", 1, "TOC depth: 1 = sections only, 2 = sections + items")
	generateCmd.Flags().BoolVar(&genOverview, "status-table", true, "Show status table at top")
	generateCmd.Flags().BoolVar(&genAreaSubheadings, "area-subheadings", false, "Show area sub-sections within phases (use with --group-by phase)")
	generateCmd.Flags().BoolVar(&genNumbered, "numbered", false, "Number items")
	generateCmd.Flags().BoolVar(&genNoRules, "no-rules", false, "Omit horizontal rules between sections")
}

func runGenerate(cmd *cobra.Command, args []string) error {
	r, err := tasks.ParseFile(genInput)
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	// Validate first
	result := tasks.Validate(r)
	if !result.Valid {
		fmt.Fprintf(cmd.ErrOrStderr(), "Validation errors in %s:\n", genInput)
		for _, e := range result.Errors {
			fmt.Fprintf(cmd.ErrOrStderr(), "  • %s: %s\n", e.Field, e.Message)
		}
		return fmt.Errorf("validation failed with %d error(s)", len(result.Errors))
	}

	// Build options
	opts := renderer.DefaultOptions()
	opts.UseCheckboxes = genCheckbox
	opts.UseEmoji = genEmoji
	opts.ShowLegend = genLegend
	opts.ShowIntro = !genNoIntro
	opts.ShowTOC = genTOC
	opts.TOCDepth = genTOCDepth
	opts.ShowOverviewTable = genOverview
	opts.ShowAreaSubheadings = genAreaSubheadings
	opts.NumberItems = genNumbered
	opts.HorizontalRules = !genNoRules

	switch genGroupBy {
	case "area":
		opts.GroupBy = renderer.GroupByArea
	case "type":
		opts.GroupBy = renderer.GroupByType
	case "phase":
		opts.GroupBy = renderer.GroupByPhase
	case "status":
		opts.GroupBy = renderer.GroupByStatus
	default:
		return fmt.Errorf("unknown group-by value: %s", genGroupBy)
	}

	// Render
	output := renderer.Render(r, opts)

	// Write output
	if genOutput == "" {
		fmt.Fprint(cmd.OutOrStdout(), output)
	} else {
		if err := os.WriteFile(genOutput, []byte(output), 0600); err != nil {
			return fmt.Errorf("failed to write output: %w", err)
		}
		fmt.Fprintf(cmd.ErrOrStderr(), "Generated %s\n", genOutput)
	}
	return nil
}
