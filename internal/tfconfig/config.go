package tfconfig

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
)

func loadConfig(file *hcl.File, mod *Module) hcl.Diagnostics {
	var diags hcl.Diagnostics
	content, _, contentDiags := file.Body.PartialContent(rootSchema)
	diags = append(diags, contentDiags...)

	for _, block := range content.Blocks {
		switch block.Type {
		case "terraform":
			content, _, contentDiags := block.Body.PartialContent(terraformBlockSchema)
			diags = append(diags, contentDiags...)

			if attr, defined := content.Attributes["required_version"]; defined {
				var version string
				valDiags := gohcl.DecodeExpression(attr.Expr, nil, &version)
				diags = append(diags, valDiags...)
				if !valDiags.HasErrors() {
					mod.RequiredCore = append(mod.RequiredCore, version)
				}
			}

			for _, innerBlock := range content.Blocks {
				switch innerBlock.Type {
				case "backend":
					label := innerBlock.Labels[0]

					mod.Backend = &Backend{
						Type:       label,
						Attributes: map[string]string{},
					}

					attrs, _ := innerBlock.Body.JustAttributes()

					if attr, defined := attrs["bucket"]; defined {
						var bucket string
						valDiags := gohcl.DecodeExpression(attr.Expr, nil, &bucket)
						diags = append(diags, valDiags...)

						mod.Backend.Attributes["bucket"] = bucket
					}

				}
			}
		}
	}

	return diags
}
