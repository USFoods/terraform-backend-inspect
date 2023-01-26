package tfconfig

// Module is the top-level type representing a parsed and processed Terraform
// module.
type Module struct {
	// Path is the local filesystem directory where the module was loaded from.
	Path string `json:"path"`

	RequiredCore []string `json:"required_core,omitempty"`
	Backends     []*Backend
}
