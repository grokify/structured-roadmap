package roadmap

import (
	"fmt"
	"regexp"

	"github.com/grokify/structured-changelog/changelog"
)

// ValidationError represents a validation error.
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidationResult holds the results of validation.
type ValidationResult struct {
	Valid  bool
	Errors []ValidationError
}

// Validate checks a Roadmap for validity.
func Validate(r *Roadmap) ValidationResult {
	result := ValidationResult{Valid: true}

	// Required fields
	if r.IRVersion == "" {
		result.addError("ir_version", "required field is missing")
	} else if r.IRVersion != "1.0" {
		result.addError("ir_version", fmt.Sprintf("unsupported version: %s", r.IRVersion))
	}

	if r.Project == "" {
		result.addError("project", "required field is missing")
	}

	// Validate items
	itemIDs := make(map[string]bool)
	for i, item := range r.Items {
		prefix := fmt.Sprintf("items[%d]", i)

		if item.ID == "" {
			result.addError(prefix+".id", "required field is missing")
		} else if itemIDs[item.ID] {
			result.addError(prefix+".id", fmt.Sprintf("duplicate ID: %s", item.ID))
		} else {
			itemIDs[item.ID] = true
		}

		if item.Title == "" {
			result.addError(prefix+".title", "required field is missing")
		}

		if item.Status == "" {
			result.addError(prefix+".status", "required field is missing")
		} else if !isValidStatus(item.Status) {
			result.addError(prefix+".status", fmt.Sprintf("invalid status: %s", item.Status))
		}

		// Validate target_quarter format
		if item.TargetQuarter != "" {
			if !isValidQuarter(item.TargetQuarter) {
				result.addError(prefix+".target_quarter", fmt.Sprintf("invalid format: %s (expected 'Q1 2026')", item.TargetQuarter))
			}
		}

		// Validate type against structured-changelog change types
		if item.Type != "" {
			if !changelog.DefaultRegistry.IsValidName(item.Type) {
				result.addError(prefix+".type", fmt.Sprintf("invalid change type: %s (see structured-changelog for valid types)", item.Type))
			}
		}

		// Validate tasks
		for j, task := range item.Tasks {
			taskPrefix := fmt.Sprintf("%s.tasks[%d]", prefix, j)
			if task.Description == "" {
				result.addError(taskPrefix+".description", "required field is missing")
			}
		}

		// Validate content blocks
		for j, block := range item.Content {
			blockPrefix := fmt.Sprintf("%s.content[%d]", prefix, j)
			if err := validateContentBlock(block, blockPrefix); err != nil {
				result.Errors = append(result.Errors, *err)
				result.Valid = false
			}
		}
	}

	// Validate depends_on references
	for i, item := range r.Items {
		for _, dep := range item.DependsOn {
			if !itemIDs[dep] {
				result.addError(fmt.Sprintf("items[%d].depends_on", i), fmt.Sprintf("references unknown item: %s", dep))
			}
		}
	}

	// Validate areas
	areaIDs := make(map[string]bool)
	for i, area := range r.Areas {
		prefix := fmt.Sprintf("areas[%d]", i)
		if area.ID == "" {
			result.addError(prefix+".id", "required field is missing")
		} else if areaIDs[area.ID] {
			result.addError(prefix+".id", fmt.Sprintf("duplicate ID: %s", area.ID))
		} else {
			areaIDs[area.ID] = true
		}
		if area.Name == "" {
			result.addError(prefix+".name", "required field is missing")
		}
	}

	// Validate phases
	phaseIDs := make(map[string]bool)
	for i, phase := range r.Phases {
		prefix := fmt.Sprintf("phases[%d]", i)
		if phase.ID == "" {
			result.addError(prefix+".id", "required field is missing")
		} else if phaseIDs[phase.ID] {
			result.addError(prefix+".id", fmt.Sprintf("duplicate ID: %s", phase.ID))
		} else {
			phaseIDs[phase.ID] = true
		}
		if phase.Name == "" {
			result.addError(prefix+".name", "required field is missing")
		}
		if phase.Status != "" && !isValidStatus(phase.Status) {
			result.addError(prefix+".status", fmt.Sprintf("invalid status: %s", phase.Status))
		}
	}

	// Validate item area/phase references
	for i, item := range r.Items {
		if item.Area != "" && len(r.Areas) > 0 && !areaIDs[item.Area] {
			result.addError(fmt.Sprintf("items[%d].area", i), fmt.Sprintf("references unknown area: %s", item.Area))
		}
		if item.Phase != "" && len(r.Phases) > 0 && !phaseIDs[item.Phase] {
			result.addError(fmt.Sprintf("items[%d].phase", i), fmt.Sprintf("references unknown phase: %s", item.Phase))
		}
	}

	// Validate sections
	sectionIDs := make(map[string]bool)
	for i, section := range r.Sections {
		prefix := fmt.Sprintf("sections[%d]", i)
		if section.ID == "" {
			result.addError(prefix+".id", "required field is missing")
		} else if sectionIDs[section.ID] {
			result.addError(prefix+".id", fmt.Sprintf("duplicate ID: %s", section.ID))
		} else {
			sectionIDs[section.ID] = true
		}
		if section.Title == "" {
			result.addError(prefix+".title", "required field is missing")
		}

		// Validate section content blocks
		for j, block := range section.Content {
			blockPrefix := fmt.Sprintf("%s.content[%d]", prefix, j)
			if err := validateContentBlock(block, blockPrefix); err != nil {
				result.Errors = append(result.Errors, *err)
				result.Valid = false
			}
		}
	}

	return result
}

func (r *ValidationResult) addError(field, message string) {
	r.Errors = append(r.Errors, ValidationError{Field: field, Message: message})
	r.Valid = false
}

func isValidStatus(s Status) bool {
	switch s {
	case StatusCompleted, StatusInProgress, StatusPlanned, StatusFuture:
		return true
	}
	return false
}

var quarterRegex = regexp.MustCompile(`^Q[1-4] \d{4}$`)

func isValidQuarter(q string) bool {
	return quarterRegex.MatchString(q)
}

func validateContentBlock(block ContentBlock, prefix string) *ValidationError {
	switch block.Type {
	case ContentTypeText, ContentTypeCode, ContentTypeDiagram, ContentTypeBlockquote:
		if block.Value == "" {
			return &ValidationError{Field: prefix + ".value", Message: "required for type " + string(block.Type)}
		}
	case ContentTypeTable:
		if len(block.Headers) == 0 {
			return &ValidationError{Field: prefix + ".headers", Message: "required for type table"}
		}
	case ContentTypeList:
		if len(block.Items) == 0 {
			return &ValidationError{Field: prefix + ".items", Message: "required for type list"}
		}
	case "":
		return &ValidationError{Field: prefix + ".type", Message: "required field is missing"}
	default:
		return &ValidationError{Field: prefix + ".type", Message: fmt.Sprintf("unknown type: %s", block.Type)}
	}
	return nil
}
