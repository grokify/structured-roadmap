package schema

import (
	"encoding/json"
	"testing"
)

func TestSchemaV1Embedded(t *testing.T) {
	if len(SchemaV1) == 0 {
		t.Fatal("SchemaV1 should not be empty")
	}

	// Verify it's valid JSON
	var schema map[string]interface{}
	if err := json.Unmarshal(SchemaV1, &schema); err != nil {
		t.Fatalf("SchemaV1 is not valid JSON: %v", err)
	}

	// Check for expected schema properties
	if schema["$schema"] != "http://json-schema.org/draft-07/schema#" {
		t.Errorf("Expected JSON Schema draft-07, got %v", schema["$schema"])
	}

	if schema["title"] != "Structured Tasks IR" {
		t.Errorf("Expected title 'Structured Tasks IR', got %v", schema["title"])
	}

	// Check required fields are defined
	required, ok := schema["required"].([]interface{})
	if !ok {
		t.Fatal("Expected required field to be an array")
	}

	hasIRVersion := false
	hasProject := false
	for _, r := range required {
		if r == "irVersion" {
			hasIRVersion = true
		}
		if r == "project" {
			hasProject = true
		}
	}
	if !hasIRVersion {
		t.Error("Expected 'irVersion' in required fields")
	}
	if !hasProject {
		t.Error("Expected 'project' in required fields")
	}
}

func TestSchemaVersion(t *testing.T) {
	v := SchemaVersion()
	if v != "1.0" {
		t.Errorf("SchemaVersion() = %q, want %q", v, "1.0")
	}
}
