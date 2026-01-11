package roadmap

import (
	"errors"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		wantErr bool
	}{
		{
			name: "minimal valid roadmap",
			json: `{"ir_version": "1.0", "project": "test-project"}`,
		},
		{
			name: "roadmap with areas",
			json: `{
				"ir_version": "1.0",
				"project": "test-project",
				"areas": [
					{"id": "core", "name": "Core", "priority": 1}
				]
			}`,
		},
		{
			name: "roadmap with items",
			json: `{
				"ir_version": "1.0",
				"project": "test-project",
				"items": [
					{"id": "item-1", "title": "Feature 1", "status": "completed"},
					{"id": "item-2", "title": "Feature 2", "status": "planned"}
				]
			}`,
		},
		{
			name:    "invalid json",
			json:    `{invalid}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := Parse([]byte(tt.json))
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && r == nil {
				t.Error("Parse() returned nil roadmap for valid input")
			}
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name      string
		roadmap   *Roadmap
		wantValid bool
	}{
		{
			name: "valid minimal roadmap",
			roadmap: &Roadmap{
				IRVersion: "1.0",
				Project:   "test",
			},
			wantValid: true,
		},
		{
			name: "missing ir_version",
			roadmap: &Roadmap{
				Project: "test",
			},
			wantValid: false,
		},
		{
			name: "missing project",
			roadmap: &Roadmap{
				IRVersion: "1.0",
			},
			wantValid: false,
		},
		{
			name: "unsupported ir_version",
			roadmap: &Roadmap{
				IRVersion: "2.0",
				Project:   "test",
			},
			wantValid: false,
		},
		{
			name: "valid roadmap with items",
			roadmap: &Roadmap{
				IRVersion: "1.0",
				Project:   "test",
				Items: []Item{
					{ID: "item-1", Title: "Feature", Status: StatusCompleted},
				},
			},
			wantValid: true,
		},
		{
			name: "item missing id",
			roadmap: &Roadmap{
				IRVersion: "1.0",
				Project:   "test",
				Items: []Item{
					{Title: "Feature", Status: StatusCompleted},
				},
			},
			wantValid: false,
		},
		{
			name: "item missing title",
			roadmap: &Roadmap{
				IRVersion: "1.0",
				Project:   "test",
				Items: []Item{
					{ID: "item-1", Status: StatusCompleted},
				},
			},
			wantValid: false,
		},
		{
			name: "item missing status",
			roadmap: &Roadmap{
				IRVersion: "1.0",
				Project:   "test",
				Items: []Item{
					{ID: "item-1", Title: "Feature"},
				},
			},
			wantValid: false,
		},
		{
			name: "duplicate item ids",
			roadmap: &Roadmap{
				IRVersion: "1.0",
				Project:   "test",
				Items: []Item{
					{ID: "item-1", Title: "Feature 1", Status: StatusCompleted},
					{ID: "item-1", Title: "Feature 2", Status: StatusPlanned},
				},
			},
			wantValid: false,
		},
		{
			name: "invalid status",
			roadmap: &Roadmap{
				IRVersion: "1.0",
				Project:   "test",
				Items: []Item{
					{ID: "item-1", Title: "Feature", Status: "invalid"},
				},
			},
			wantValid: false,
		},
		{
			name: "valid target quarter",
			roadmap: &Roadmap{
				IRVersion: "1.0",
				Project:   "test",
				Items: []Item{
					{ID: "item-1", Title: "Feature", Status: StatusPlanned, TargetQuarter: "Q2 2026"},
				},
			},
			wantValid: true,
		},
		{
			name: "invalid target quarter format",
			roadmap: &Roadmap{
				IRVersion: "1.0",
				Project:   "test",
				Items: []Item{
					{ID: "item-1", Title: "Feature", Status: StatusPlanned, TargetQuarter: "2026-Q2"},
				},
			},
			wantValid: false,
		},
		{
			name: "valid type from structured-changelog",
			roadmap: &Roadmap{
				IRVersion: "1.0",
				Project:   "test",
				Items: []Item{
					{ID: "item-1", Title: "Feature", Status: StatusCompleted, Type: "Added"},
				},
			},
			wantValid: true,
		},
		{
			name: "invalid type",
			roadmap: &Roadmap{
				IRVersion: "1.0",
				Project:   "test",
				Items: []Item{
					{ID: "item-1", Title: "Feature", Status: StatusCompleted, Type: "InvalidType"},
				},
			},
			wantValid: false,
		},
		{
			name: "valid dependency reference",
			roadmap: &Roadmap{
				IRVersion: "1.0",
				Project:   "test",
				Items: []Item{
					{ID: "item-1", Title: "Feature 1", Status: StatusCompleted},
					{ID: "item-2", Title: "Feature 2", Status: StatusPlanned, DependsOn: []string{"item-1"}},
				},
			},
			wantValid: true,
		},
		{
			name: "invalid dependency reference",
			roadmap: &Roadmap{
				IRVersion: "1.0",
				Project:   "test",
				Items: []Item{
					{ID: "item-1", Title: "Feature 1", Status: StatusCompleted},
					{ID: "item-2", Title: "Feature 2", Status: StatusPlanned, DependsOn: []string{"nonexistent"}},
				},
			},
			wantValid: false,
		},
		{
			name: "valid area reference",
			roadmap: &Roadmap{
				IRVersion: "1.0",
				Project:   "test",
				Areas:     []Area{{ID: "core", Name: "Core"}},
				Items: []Item{
					{ID: "item-1", Title: "Feature", Status: StatusCompleted, Area: "core"},
				},
			},
			wantValid: true,
		},
		{
			name: "invalid area reference",
			roadmap: &Roadmap{
				IRVersion: "1.0",
				Project:   "test",
				Areas:     []Area{{ID: "core", Name: "Core"}},
				Items: []Item{
					{ID: "item-1", Title: "Feature", Status: StatusCompleted, Area: "nonexistent"},
				},
			},
			wantValid: false,
		},
		{
			name: "valid phase reference",
			roadmap: &Roadmap{
				IRVersion: "1.0",
				Project:   "test",
				Phases:    []Phase{{ID: "phase-1", Name: "Phase 1"}},
				Items: []Item{
					{ID: "item-1", Title: "Feature", Status: StatusCompleted, Phase: "phase-1"},
				},
			},
			wantValid: true,
		},
		{
			name: "invalid phase reference",
			roadmap: &Roadmap{
				IRVersion: "1.0",
				Project:   "test",
				Phases:    []Phase{{ID: "phase-1", Name: "Phase 1"}},
				Items: []Item{
					{ID: "item-1", Title: "Feature", Status: StatusCompleted, Phase: "nonexistent"},
				},
			},
			wantValid: false,
		},
		{
			name: "valid content block - text",
			roadmap: &Roadmap{
				IRVersion: "1.0",
				Project:   "test",
				Items: []Item{
					{
						ID: "item-1", Title: "Feature", Status: StatusCompleted,
						Content: []ContentBlock{{Type: ContentTypeText, Value: "Some text"}},
					},
				},
			},
			wantValid: true,
		},
		{
			name: "valid content block - blockquote",
			roadmap: &Roadmap{
				IRVersion: "1.0",
				Project:   "test",
				Items: []Item{
					{
						ID: "item-1", Title: "Feature", Status: StatusCompleted,
						Content: []ContentBlock{{Type: ContentTypeBlockquote, Value: "A quote"}},
					},
				},
			},
			wantValid: true,
		},
		{
			name: "content block missing value",
			roadmap: &Roadmap{
				IRVersion: "1.0",
				Project:   "test",
				Items: []Item{
					{
						ID: "item-1", Title: "Feature", Status: StatusCompleted,
						Content: []ContentBlock{{Type: ContentTypeText}},
					},
				},
			},
			wantValid: false,
		},
		{
			name: "content block missing type",
			roadmap: &Roadmap{
				IRVersion: "1.0",
				Project:   "test",
				Items: []Item{
					{
						ID: "item-1", Title: "Feature", Status: StatusCompleted,
						Content: []ContentBlock{{Value: "text"}},
					},
				},
			},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Validate(tt.roadmap)
			if result.Valid != tt.wantValid {
				t.Errorf("Validate() valid = %v, want %v", result.Valid, tt.wantValid)
				if len(result.Errors) > 0 {
					for _, e := range result.Errors {
						t.Logf("  Error: %s: %s", e.Field, e.Message)
					}
				}
			}
		})
	}
}

func TestPriorityOrder(t *testing.T) {
	tests := []struct {
		priority Priority
		want     int
	}{
		{PriorityCritical, 1},
		{PriorityHigh, 2},
		{PriorityMedium, 3},
		{PriorityLow, 4},
		{"unknown", 5},
	}

	for _, tt := range tests {
		t.Run(string(tt.priority), func(t *testing.T) {
			if got := PriorityOrder(tt.priority); got != tt.want {
				t.Errorf("PriorityOrder(%q) = %v, want %v", tt.priority, got, tt.want)
			}
		})
	}
}

func TestPriorityLabel(t *testing.T) {
	tests := []struct {
		priority Priority
		want     string
	}{
		{PriorityCritical, "Critical"},
		{PriorityHigh, "High"},
		{PriorityMedium, "Medium"},
		{PriorityLow, "Low"},
		{"unknown", "-"},
	}

	for _, tt := range tests {
		t.Run(string(tt.priority), func(t *testing.T) {
			if got := PriorityLabel(tt.priority); got != tt.want {
				t.Errorf("PriorityLabel(%q) = %v, want %v", tt.priority, got, tt.want)
			}
		})
	}
}

func TestStats(t *testing.T) {
	r := &Roadmap{
		IRVersion: "1.0",
		Project:   "test",
		Items: []Item{
			{ID: "1", Title: "Item 1", Status: StatusCompleted, Area: "core", Type: "Added", Priority: PriorityHigh},
			{ID: "2", Title: "Item 2", Status: StatusCompleted, Area: "core", Type: "Added", Priority: PriorityHigh},
			{ID: "3", Title: "Item 3", Status: StatusInProgress, Area: "api", Type: "Changed", Priority: PriorityMedium},
			{ID: "4", Title: "Item 4", Status: StatusPlanned, Area: "api", Type: "Added", Priority: PriorityLow},
			{ID: "5", Title: "Item 5", Status: StatusFuture},
		},
	}

	stats := r.Stats()

	if stats.Total != 5 {
		t.Errorf("Total = %d, want 5", stats.Total)
	}
	if stats.ByStatus[StatusCompleted] != 2 {
		t.Errorf("ByStatus[completed] = %d, want 2", stats.ByStatus[StatusCompleted])
	}
	if stats.ByStatus[StatusInProgress] != 1 {
		t.Errorf("ByStatus[in_progress] = %d, want 1", stats.ByStatus[StatusInProgress])
	}
	if stats.ByArea["core"] != 2 {
		t.Errorf("ByArea[core] = %d, want 2", stats.ByArea["core"])
	}
	if stats.ByType["Added"] != 3 {
		t.Errorf("ByType[Added] = %d, want 3", stats.ByType["Added"])
	}
	if stats.CompletedCount() != 2 {
		t.Errorf("CompletedCount() = %d, want 2", stats.CompletedCount())
	}
	if stats.CompletedPercent() != 40.0 {
		t.Errorf("CompletedPercent() = %f, want 40.0", stats.CompletedPercent())
	}
}

func TestItemsBy(t *testing.T) {
	r := &Roadmap{
		IRVersion: "1.0",
		Project:   "test",
		Items: []Item{
			{ID: "1", Title: "Item 1", Status: StatusCompleted, Area: "core", Phase: "phase-1", Priority: PriorityHigh, TargetQuarter: "Q1 2026"},
			{ID: "2", Title: "Item 2", Status: StatusPlanned, Area: "api", Phase: "phase-1", Priority: PriorityMedium, TargetQuarter: "Q2 2026"},
			{ID: "3", Title: "Item 3", Status: StatusPlanned, Area: "core", Phase: "phase-2", Priority: PriorityLow},
		},
	}

	t.Run("ItemsByArea", func(t *testing.T) {
		byArea := r.ItemsByArea()
		if len(byArea["core"]) != 2 {
			t.Errorf("ItemsByArea[core] = %d items, want 2", len(byArea["core"]))
		}
		if len(byArea["api"]) != 1 {
			t.Errorf("ItemsByArea[api] = %d items, want 1", len(byArea["api"]))
		}
	})

	t.Run("ItemsByPhase", func(t *testing.T) {
		byPhase := r.ItemsByPhase()
		if len(byPhase["phase-1"]) != 2 {
			t.Errorf("ItemsByPhase[phase-1] = %d items, want 2", len(byPhase["phase-1"]))
		}
	})

	t.Run("ItemsByStatus", func(t *testing.T) {
		byStatus := r.ItemsByStatus()
		if len(byStatus[StatusCompleted]) != 1 {
			t.Errorf("ItemsByStatus[completed] = %d items, want 1", len(byStatus[StatusCompleted]))
		}
		if len(byStatus[StatusPlanned]) != 2 {
			t.Errorf("ItemsByStatus[planned] = %d items, want 2", len(byStatus[StatusPlanned]))
		}
	})

	t.Run("ItemsByQuarter", func(t *testing.T) {
		byQuarter := r.ItemsByQuarter()
		if len(byQuarter["Q1 2026"]) != 1 {
			t.Errorf("ItemsByQuarter[Q1 2026] = %d items, want 1", len(byQuarter["Q1 2026"]))
		}
		if len(byQuarter["_unscheduled"]) != 1 {
			t.Errorf("ItemsByQuarter[_unscheduled] = %d items, want 1", len(byQuarter["_unscheduled"]))
		}
	})

	t.Run("ItemsByPriority", func(t *testing.T) {
		byPriority := r.ItemsByPriority()
		if len(byPriority[PriorityHigh]) != 1 {
			t.Errorf("ItemsByPriority[high] = %d items, want 1", len(byPriority[PriorityHigh]))
		}
	})
}

func TestDefaultLegend(t *testing.T) {
	legend := DefaultLegend()
	if legend[StatusCompleted].Emoji != "✅" {
		t.Errorf("DefaultLegend[completed].Emoji = %q, want ✅", legend[StatusCompleted].Emoji)
	}
	if legend[StatusInProgress].Emoji != "🚧" {
		t.Errorf("DefaultLegend[in_progress].Emoji = %q, want 🚧", legend[StatusInProgress].Emoji)
	}
}

func TestGetLegend(t *testing.T) {
	t.Run("uses default when no legend", func(t *testing.T) {
		r := &Roadmap{IRVersion: "1.0", Project: "test"}
		legend := r.GetLegend()
		if legend[StatusCompleted].Emoji != "✅" {
			t.Error("Expected default legend")
		}
	})

	t.Run("merges custom legend", func(t *testing.T) {
		r := &Roadmap{
			IRVersion: "1.0",
			Project:   "test",
			Legend: map[Status]LegendEntry{
				StatusCompleted: {Emoji: "✓", Description: "Done"},
			},
		}
		legend := r.GetLegend()
		if legend[StatusCompleted].Emoji != "✓" {
			t.Errorf("Expected custom emoji, got %q", legend[StatusCompleted].Emoji)
		}
		// Should still have defaults for other statuses
		if legend[StatusInProgress].Emoji != "🚧" {
			t.Errorf("Expected default emoji for in_progress, got %q", legend[StatusInProgress].Emoji)
		}
	})
}

func TestSentinelErrors(t *testing.T) {
	t.Run("Parse returns ErrParseJSON for invalid JSON", func(t *testing.T) {
		_, err := Parse([]byte(`{invalid}`))
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if !errors.Is(err, ErrParseJSON) {
			t.Errorf("Expected error to wrap ErrParseJSON, got %v", err)
		}
	})

	t.Run("ParseFile returns ErrReadFile for missing file", func(t *testing.T) {
		_, err := ParseFile("/nonexistent/file.json")
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if !errors.Is(err, ErrReadFile) {
			t.Errorf("Expected error to wrap ErrReadFile, got %v", err)
		}
	})
}

func TestFieldError(t *testing.T) {
	err := NewFieldError("items[0].id", "required field is missing", ErrMissingRequiredField)

	if err.Field != "items[0].id" {
		t.Errorf("Field = %q, want %q", err.Field, "items[0].id")
	}
	if err.Message != "required field is missing" {
		t.Errorf("Message = %q, want %q", err.Message, "required field is missing")
	}
	if !errors.Is(err, ErrMissingRequiredField) {
		t.Error("Expected error to wrap ErrMissingRequiredField")
	}

	expectedStr := "items[0].id: required field is missing"
	if err.Error() != expectedStr {
		t.Errorf("Error() = %q, want %q", err.Error(), expectedStr)
	}
}

func TestGetStatusEmoji(t *testing.T) {
	r := &Roadmap{IRVersion: "1.0", Project: "test"}

	tests := []struct {
		status Status
		want   string
	}{
		{StatusCompleted, "✅"},
		{StatusInProgress, "🚧"},
		{StatusPlanned, "📋"},
		{StatusFuture, "💡"},
		{"unknown", ""},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			got := r.GetStatusEmoji(tt.status)
			if got != tt.want {
				t.Errorf("GetStatusEmoji(%q) = %q, want %q", tt.status, got, tt.want)
			}
		})
	}
}

func TestItemsByType(t *testing.T) {
	r := &Roadmap{
		IRVersion: "1.0",
		Project:   "test",
		Items: []Item{
			{ID: "1", Title: "Item 1", Status: StatusCompleted, Type: "Added"},
			{ID: "2", Title: "Item 2", Status: StatusCompleted, Type: "Added"},
			{ID: "3", Title: "Item 3", Status: StatusCompleted, Type: "Changed"},
			{ID: "4", Title: "Item 4", Status: StatusCompleted}, // No type
		},
	}

	byType := r.ItemsByType()

	if len(byType["Added"]) != 2 {
		t.Errorf("ItemsByType[Added] = %d items, want 2", len(byType["Added"]))
	}
	if len(byType["Changed"]) != 1 {
		t.Errorf("ItemsByType[Changed] = %d items, want 1", len(byType["Changed"]))
	}
	if len(byType["_unspecified"]) != 1 {
		t.Errorf("ItemsByType[_unspecified] = %d items, want 1", len(byType["_unspecified"]))
	}
}

func TestStatusOrder(t *testing.T) {
	order := StatusOrder()
	if len(order) != 4 {
		t.Errorf("StatusOrder() returned %d items, want 4", len(order))
	}
	if order[0] != StatusCompleted {
		t.Errorf("StatusOrder()[0] = %q, want %q", order[0], StatusCompleted)
	}
	if order[3] != StatusFuture {
		t.Errorf("StatusOrder()[3] = %q, want %q", order[3], StatusFuture)
	}
}

func TestPriorityOrderList(t *testing.T) {
	order := PriorityOrderList()
	if len(order) != 4 {
		t.Errorf("PriorityOrderList() returned %d items, want 4", len(order))
	}
	if order[0] != PriorityCritical {
		t.Errorf("PriorityOrderList()[0] = %q, want %q", order[0], PriorityCritical)
	}
	if order[3] != PriorityLow {
		t.Errorf("PriorityOrderList()[3] = %q, want %q", order[3], PriorityLow)
	}
}

func TestCompletedPercentEmpty(t *testing.T) {
	stats := Stats{Total: 0, ByStatus: make(map[Status]int)}
	pct := stats.CompletedPercent()
	if pct != 0 {
		t.Errorf("CompletedPercent() with empty stats = %f, want 0", pct)
	}
}

func TestToJSON(t *testing.T) {
	r := &Roadmap{
		IRVersion: "1.0",
		Project:   "test-project",
		Items: []Item{
			{ID: "item-1", Title: "Test Item", Status: StatusCompleted},
		},
	}

	data, err := ToJSON(r)
	if err != nil {
		t.Fatalf("ToJSON() error = %v", err)
	}

	// Parse it back to verify
	r2, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if r2.Project != "test-project" {
		t.Errorf("Project = %q, want %q", r2.Project, "test-project")
	}
	if len(r2.Items) != 1 {
		t.Errorf("Items = %d, want 1", len(r2.Items))
	}
}

func TestWriteFile(t *testing.T) {
	r := &Roadmap{
		IRVersion: "1.0",
		Project:   "test-project",
	}

	// Write to temp file
	tmpFile := t.TempDir() + "/roadmap.json"
	err := WriteFile(tmpFile, r)
	if err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	// Read it back
	r2, err := ParseFile(tmpFile)
	if err != nil {
		t.Fatalf("ParseFile() error = %v", err)
	}

	if r2.Project != "test-project" {
		t.Errorf("Project = %q, want %q", r2.Project, "test-project")
	}
}

func TestParseError(t *testing.T) {
	underlying := errors.New("connection refused")
	parseErr := &ParseError{
		Op:  "read",
		Err: underlying,
	}

	expectedStr := "read: connection refused"
	if parseErr.Error() != expectedStr {
		t.Errorf("Error() = %q, want %q", parseErr.Error(), expectedStr)
	}

	if !errors.Is(parseErr, underlying) {
		t.Error("Expected Unwrap to return underlying error")
	}

	unwrapped := parseErr.Unwrap()
	if unwrapped != underlying {
		t.Errorf("Unwrap() = %v, want %v", unwrapped, underlying)
	}
}

func TestValidationError(t *testing.T) {
	err := ValidationError{
		Field:   "ir_version",
		Message: "required field is missing",
	}

	errStr := err.Error()
	if !strings.Contains(errStr, "ir_version") {
		t.Error("Expected error to contain ir_version")
	}
	if !strings.Contains(errStr, "required field is missing") {
		t.Error("Expected error to contain message")
	}
	// Verify the format: "field: message"
	expected := "ir_version: required field is missing"
	if errStr != expected {
		t.Errorf("Error() = %q, want %q", errStr, expected)
	}
}

func TestValidationResultWithErrors(t *testing.T) {
	result := ValidationResult{
		Valid: false,
		Errors: []ValidationError{
			{Field: "ir_version", Message: "required field is missing"},
			{Field: "project", Message: "required field is missing"},
		},
	}

	if result.Valid {
		t.Error("Expected result.Valid to be false")
	}
	if len(result.Errors) != 2 {
		t.Errorf("Expected 2 errors, got %d", len(result.Errors))
	}
}

func TestValidationResultValid(t *testing.T) {
	result := ValidationResult{
		Valid:  true,
		Errors: nil,
	}

	if !result.Valid {
		t.Error("Expected result.Valid to be true")
	}
	if len(result.Errors) != 0 {
		t.Errorf("Expected 0 errors, got %d", len(result.Errors))
	}
}

func TestItemsByAreaUnspecified(t *testing.T) {
	r := &Roadmap{
		IRVersion: "1.0",
		Project:   "test",
		Items: []Item{
			{ID: "1", Title: "Item 1", Status: StatusCompleted, Area: "core"},
			{ID: "2", Title: "Item 2", Status: StatusCompleted}, // No area
		},
	}

	byArea := r.ItemsByArea()
	if len(byArea["core"]) != 1 {
		t.Errorf("ItemsByArea[core] = %d items, want 1", len(byArea["core"]))
	}
	if len(byArea["_unspecified"]) != 1 {
		t.Errorf("ItemsByArea[_unspecified] = %d items, want 1", len(byArea["_unspecified"]))
	}
}

func TestItemsByPhaseUnphased(t *testing.T) {
	r := &Roadmap{
		IRVersion: "1.0",
		Project:   "test",
		Items: []Item{
			{ID: "1", Title: "Item 1", Status: StatusCompleted, Phase: "phase-1"},
			{ID: "2", Title: "Item 2", Status: StatusCompleted}, // No phase
		},
	}

	byPhase := r.ItemsByPhase()
	if len(byPhase["phase-1"]) != 1 {
		t.Errorf("ItemsByPhase[phase-1] = %d items, want 1", len(byPhase["phase-1"]))
	}
	if len(byPhase["_unphased"]) != 1 {
		t.Errorf("ItemsByPhase[_unphased] = %d items, want 1", len(byPhase["_unphased"]))
	}
}

func TestItemsByPriorityUnspecified(t *testing.T) {
	r := &Roadmap{
		IRVersion: "1.0",
		Project:   "test",
		Items: []Item{
			{ID: "1", Title: "Item 1", Status: StatusCompleted, Priority: PriorityHigh},
			{ID: "2", Title: "Item 2", Status: StatusCompleted}, // No priority
		},
	}

	byPriority := r.ItemsByPriority()
	if len(byPriority[PriorityHigh]) != 1 {
		t.Errorf("ItemsByPriority[high] = %d items, want 1", len(byPriority[PriorityHigh]))
	}
	if len(byPriority["_unspecified"]) != 1 {
		t.Errorf("ItemsByPriority[_unspecified] = %d items, want 1", len(byPriority["_unspecified"]))
	}
}

func TestValidateMoreCases(t *testing.T) {
	tests := []struct {
		name      string
		roadmap   *Roadmap
		wantValid bool
	}{
		{
			name: "duplicate area ids",
			roadmap: &Roadmap{
				IRVersion: "1.0",
				Project:   "test",
				Areas: []Area{
					{ID: "core", Name: "Core"},
					{ID: "core", Name: "Core 2"}, // Duplicate
				},
			},
			wantValid: false,
		},
		{
			name: "duplicate phase ids",
			roadmap: &Roadmap{
				IRVersion: "1.0",
				Project:   "test",
				Phases: []Phase{
					{ID: "phase-1", Name: "Phase 1"},
					{ID: "phase-1", Name: "Phase 1 Again"}, // Duplicate
				},
			},
			wantValid: false,
		},
		{
			name: "duplicate section ids",
			roadmap: &Roadmap{
				IRVersion: "1.0",
				Project:   "test",
				Sections: []Section{
					{ID: "intro", Title: "Intro"},
					{ID: "intro", Title: "Intro 2"}, // Duplicate
				},
			},
			wantValid: false,
		},
		{
			name: "valid with all priorities",
			roadmap: &Roadmap{
				IRVersion: "1.0",
				Project:   "test",
				Items: []Item{
					{ID: "1", Title: "Item 1", Status: StatusCompleted, Priority: PriorityCritical},
					{ID: "2", Title: "Item 2", Status: StatusCompleted, Priority: PriorityHigh},
					{ID: "3", Title: "Item 3", Status: StatusCompleted, Priority: PriorityMedium},
					{ID: "4", Title: "Item 4", Status: StatusCompleted, Priority: PriorityLow},
				},
			},
			wantValid: true,
		},
		{
			name: "item with tasks",
			roadmap: &Roadmap{
				IRVersion: "1.0",
				Project:   "test",
				Items: []Item{
					{
						ID:     "1",
						Title:  "Item with tasks",
						Status: StatusInProgress,
						Tasks: []Task{
							{Description: "Task 1", Completed: true},
							{Description: "Task 2", Completed: false},
						},
					},
				},
			},
			wantValid: true,
		},
		{
			name: "section with content blocks",
			roadmap: &Roadmap{
				IRVersion: "1.0",
				Project:   "test",
				Sections: []Section{
					{
						ID:    "intro",
						Title: "Intro",
						Content: []ContentBlock{
							{Type: ContentTypeText, Value: "Hello"},
						},
					},
				},
			},
			wantValid: true,
		},
		{
			name: "section with invalid content block",
			roadmap: &Roadmap{
				IRVersion: "1.0",
				Project:   "test",
				Sections: []Section{
					{
						ID:    "intro",
						Title: "Intro",
						Content: []ContentBlock{
							{Type: ContentTypeText}, // Missing value
						},
					},
				},
			},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Validate(tt.roadmap)
			if result.Valid != tt.wantValid {
				t.Errorf("Validate() valid = %v, want %v", result.Valid, tt.wantValid)
				for _, e := range result.Errors {
					t.Logf("  Error: %s: %s", e.Field, e.Message)
				}
			}
		})
	}
}

func TestValidateContentBlocks(t *testing.T) {
	tests := []struct {
		name      string
		content   []ContentBlock
		wantValid bool
	}{
		{
			name:      "valid code block",
			content:   []ContentBlock{{Type: ContentTypeCode, Value: "func main() {}", Language: "go"}},
			wantValid: true,
		},
		{
			name:      "valid diagram",
			content:   []ContentBlock{{Type: ContentTypeDiagram, Value: "A --> B", Format: "mermaid"}},
			wantValid: true,
		},
		{
			name:      "valid table",
			content:   []ContentBlock{{Type: ContentTypeTable, Headers: []string{"A", "B"}, Rows: [][]string{{"1", "2"}}}},
			wantValid: true,
		},
		{
			name:      "table missing headers",
			content:   []ContentBlock{{Type: ContentTypeTable, Rows: [][]string{{"1", "2"}}}},
			wantValid: false,
		},
		{
			name:      "valid list",
			content:   []ContentBlock{{Type: ContentTypeList, Items: []string{"item1", "item2"}}},
			wantValid: true,
		},
		{
			name:      "list missing items",
			content:   []ContentBlock{{Type: ContentTypeList}},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Roadmap{
				IRVersion: "1.0",
				Project:   "test",
				Items: []Item{
					{ID: "item-1", Title: "Test", Status: StatusCompleted, Content: tt.content},
				},
			}
			result := Validate(r)
			if result.Valid != tt.wantValid {
				t.Errorf("Validate() = %v, want %v", result.Valid, tt.wantValid)
				for _, e := range result.Errors {
					t.Logf("  Error: %s: %s", e.Field, e.Message)
				}
			}
		})
	}
}
