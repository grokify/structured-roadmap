package roadmap

import (
	"errors"
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
		{PriorityHigh, "High Priority"},
		{PriorityMedium, "Medium Priority"},
		{PriorityLow, "Low Priority"},
		{"unknown", "Unspecified"},
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
