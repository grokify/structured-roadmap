package roadmap

import (
	"encoding/json"
	"fmt"
	"os"
)

// ParseFile reads and parses a ROADMAP.json file.
func ParseFile(path string) (*Roadmap, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrReadFile, err)
	}
	return Parse(data)
}

// Parse parses JSON data into a Roadmap.
func Parse(data []byte) (*Roadmap, error) {
	var r Roadmap
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrParseJSON, err)
	}
	return &r, nil
}

// WriteFile writes a Roadmap to a JSON file.
func WriteFile(path string, r *Roadmap) error {
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return fmt.Errorf("%w: %v", ErrWriteFile, err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("%w: %v", ErrWriteFile, err)
	}
	return nil
}

// ToJSON converts a Roadmap to JSON bytes.
func ToJSON(r *Roadmap) ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}
