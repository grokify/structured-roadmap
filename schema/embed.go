// Package schema provides embedded JSON schema for roadmap validation.
package schema

import (
	_ "embed"
)

// SchemaV1 contains the embedded JSON schema for roadmap v1.0.
//
//go:embed roadmap.v1.schema.json
var SchemaV1 []byte

// SchemaVersion returns the current schema version.
func SchemaVersion() string {
	return "1.0"
}
