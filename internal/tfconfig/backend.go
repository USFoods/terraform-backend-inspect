package tfconfig

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
)

type Backend struct {
	Type       string
	Attributes map[string]string
}

func decodeBackendConfig(block *hcl.Block) (*Backend, hcl.Diagnostics) {
	label := block.Labels[0]

	backend := &Backend{
		Type:       label,
		Attributes: map[string]string{},
	}

	attrs, diags := block.Body.JustAttributes()

	switch label {
	case "local":
		if attr, defined := attrs["path"]; defined {
			var path string
			valDiags := gohcl.DecodeExpression(attr.Expr, nil, &path)
			diags = append(diags, valDiags...)

			backend.Attributes["bucket"] = path
		}
	case "s3":
		if attr, defined := attrs["bucket"]; defined {
			var bucket string
			valDiags := gohcl.DecodeExpression(attr.Expr, nil, &bucket)
			diags = append(diags, valDiags...)

			backend.Attributes["bucket"] = bucket
		}

		if attr, defined := attrs["key"]; defined {
			var key string
			valDiags := gohcl.DecodeExpression(attr.Expr, nil, &key)
			diags = append(diags, valDiags...)

			backend.Attributes["key"] = key
		}
	}

	return backend, diags
}
