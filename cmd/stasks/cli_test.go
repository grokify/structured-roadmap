package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// Helper to execute a cobra command and capture output
func executeCommand(root *cobra.Command, args ...string) (string, string, error) {
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	root.SetOut(stdout)
	root.SetErr(stderr)
	root.SetArgs(args)

	err := root.Execute()

	return stdout.String(), stderr.String(), err
}

func TestValidateCommand(t *testing.T) {
	// Create a temporary valid JSON file
	tmpDir := t.TempDir()
	validJSON := `{
		"irVersion": "1.0",
		"project": "test-project",
		"tasks": [
			{"id": "task-1", "title": "Feature 1", "status": "completed"}
		]
	}`
	validFile := filepath.Join(tmpDir, "valid.json")
	if err := os.WriteFile(validFile, []byte(validJSON), 0600); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create an invalid JSON file
	invalidJSON := `{
		"irVersion": "1.0"
	}`
	invalidFile := filepath.Join(tmpDir, "invalid.json")
	if err := os.WriteFile(invalidFile, []byte(invalidJSON), 0600); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name      string
		args      []string
		wantErr   bool
		wantInOut string
	}{
		{
			name:      "valid file",
			args:      []string{"validate", validFile},
			wantErr:   false,
			wantInOut: "valid",
		},
		{
			name:    "invalid file - missing project",
			args:    []string{"validate", invalidFile},
			wantErr: true,
		},
		{
			name:    "nonexistent file",
			args:    []string{"validate", "/nonexistent/file.json"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a fresh root command for each test
			cmd := &cobra.Command{Use: "stasks"}
			cmd.AddCommand(validateCmd)

			_, stderr, err := executeCommand(cmd, tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("validate error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantInOut != "" && !strings.Contains(stderr, tt.wantInOut) {
				t.Errorf("Expected output to contain %q", tt.wantInOut)
			}
		})
	}
}

func TestGenerateCommand(t *testing.T) {
	// Create a temporary valid JSON file
	tmpDir := t.TempDir()
	validJSON := `{
		"irVersion": "1.0",
		"project": "Test Project",
		"areas": [
			{"id": "core", "name": "Core Features"}
		],
		"tasks": [
			{"id": "task-1", "title": "Feature 1", "status": "completed", "area": "core"}
		]
	}`
	inputFile := filepath.Join(tmpDir, "TASKS.json")
	if err := os.WriteFile(inputFile, []byte(validJSON), 0600); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	t.Run("generate to stdout", func(t *testing.T) {
		cmd := &cobra.Command{Use: "stasks"}
		cmd.AddCommand(generateCmd)

		stdout, _, err := executeCommand(cmd, "generate", "-i", inputFile)
		if err != nil {
			t.Fatalf("generate failed: %v", err)
		}

		if !strings.Contains(stdout, "# Task List") {
			t.Error("Expected task list title in output")
		}
		if !strings.Contains(stdout, "Core Features") {
			t.Error("Expected area heading in output")
		}
		if !strings.Contains(stdout, "[x]") {
			t.Error("Expected checkbox in output")
		}
	})

	t.Run("generate to file", func(t *testing.T) {
		outputFile := filepath.Join(tmpDir, "TASKS.md")

		cmd := &cobra.Command{Use: "stasks"}
		cmd.AddCommand(generateCmd)

		_, _, err := executeCommand(cmd, "generate", "-i", inputFile, "-o", outputFile)
		if err != nil {
			t.Fatalf("generate failed: %v", err)
		}

		content, err := os.ReadFile(outputFile)
		if err != nil {
			t.Fatalf("Failed to read output file: %v", err)
		}

		if !strings.Contains(string(content), "# Task List") {
			t.Error("Expected task list title in output file")
		}
	})

	t.Run("generate with options", func(t *testing.T) {
		// Reset global flags
		genTOC = false
		genLegend = false
		genInput = "TASKS.json"
		genOutput = ""

		cmd := &cobra.Command{Use: "stasks"}
		cmd.AddCommand(generateCmd)

		stdout, _, err := executeCommand(cmd, "generate", "-i", inputFile, "--toc", "--legend")
		if err != nil {
			t.Fatalf("generate failed: %v", err)
		}

		if !strings.Contains(stdout, "## Table of Contents") {
			t.Error("Expected TOC in output")
		}
		if !strings.Contains(stdout, "## Legend") {
			t.Error("Expected legend in output")
		}
	})
}

func TestStatsCommand(t *testing.T) {
	// Create a temporary JSON file
	tmpDir := t.TempDir()
	validJSON := `{
		"irVersion": "1.0",
		"project": "Test Project",
		"tasks": [
			{"id": "1", "title": "Task 1", "status": "completed"},
			{"id": "2", "title": "Task 2", "status": "planned"},
			{"id": "3", "title": "Task 3", "status": "planned"}
		]
	}`
	inputFile := filepath.Join(tmpDir, "TASKS.json")
	if err := os.WriteFile(inputFile, []byte(validJSON), 0600); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cmd := &cobra.Command{Use: "stasks"}
	cmd.AddCommand(statsCmd)

	stdout, _, err := executeCommand(cmd, "stats", inputFile)
	if err != nil {
		t.Fatalf("stats failed: %v", err)
	}

	expectedOutputs := []string{
		"Test Project",
		"Total tasks: 3",
		"Completed",
		"Planned",
	}

	for _, expected := range expectedOutputs {
		if !strings.Contains(stdout, expected) {
			t.Errorf("Expected output to contain %q, got:\n%s", expected, stdout)
		}
	}
}

func TestDepsCommand(t *testing.T) {
	// Create a temporary JSON file with dependencies
	tmpDir := t.TempDir()
	validJSON := `{
		"irVersion": "1.0",
		"project": "Test Project",
		"tasks": [
			{"id": "task-1", "title": "Foundation", "status": "completed"},
			{"id": "task-2", "title": "Feature A", "status": "planned", "dependsOn": ["task-1"]},
			{"id": "task-3", "title": "Feature B", "status": "planned", "dependsOn": ["task-1", "task-2"]}
		]
	}`
	inputFile := filepath.Join(tmpDir, "TASKS.json")
	if err := os.WriteFile(inputFile, []byte(validJSON), 0600); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	t.Run("mermaid format", func(t *testing.T) {
		cmd := &cobra.Command{Use: "stasks"}
		cmd.AddCommand(depsCmd)

		stdout, _, err := executeCommand(cmd, "deps", inputFile, "--format", "mermaid")
		if err != nil {
			t.Fatalf("deps failed: %v", err)
		}

		if !strings.Contains(stdout, "graph TD") {
			t.Error("Expected Mermaid graph declaration")
		}
		if !strings.Contains(stdout, "-->") {
			t.Error("Expected dependency arrows")
		}
	})

	t.Run("dot format", func(t *testing.T) {
		cmd := &cobra.Command{Use: "stasks"}
		cmd.AddCommand(depsCmd)

		stdout, _, err := executeCommand(cmd, "deps", inputFile, "--format", "dot")
		if err != nil {
			t.Fatalf("deps failed: %v", err)
		}

		if !strings.Contains(stdout, "digraph") {
			t.Error("Expected DOT digraph declaration")
		}
		if !strings.Contains(stdout, "->") {
			t.Error("Expected dependency arrows")
		}
	})
}
