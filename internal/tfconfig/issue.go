package tfconfig

import "github.com/hashicorp/hcl/v2"

type Issue struct {
	Severity   Severity
	Message    string
	ModulePath string
	Range      hcl.Range
}

// Issues is an alias for the map of Issue
type Issues []*Issue

type Severity int32

// String returns the string representation of the severity.
func (s Severity) String() string {
	switch s {
	case ERROR:
		return "Error"
	case WARNING:
		return "Warning"
	case NOTICE:
		return "Notice"
	}

	return "Unknown"
}

const (
	// ERROR is possible errors
	ERROR Severity = iota
	// WARNING doesn't cause problem immediately, but not good
	WARNING
	// NOTICE is not important, it's mentioned
	NOTICE
)
