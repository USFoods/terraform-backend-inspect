package tfconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
)

func LoadModules(workingDirs []string) (modules []*Module, diags hcl.Diagnostics) {
	for _, wd := range workingDirs {
		// attempt to find any Terraform files
		configFiles, err := findConfigFiles(wd)

		if err != nil {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Failed to read module directory",
				Detail:   fmt.Sprintf("Module directory %s does not exist or cannot be read.", wd),
			})

			continue
		}

		if len(configFiles) > 0 {
			module, moduleDiags := loadModule(configFiles, wd)
			modules = append(modules, module)
			diags = append(diags, moduleDiags...)
		}
	}

	return
}

func loadModule(moduleFiles []string, moduleDir string) (mod *Module, diags hcl.Diagnostics) {
	mod = &Module{Path: moduleDir}

	parser := hclparse.NewParser()

	for _, filename := range moduleFiles {
		var file *hcl.File
		var fileDiags hcl.Diagnostics

		b, err := os.ReadFile(filename)

		if err != nil {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Failed to read file",
				Detail:   fmt.Sprintf("The configuration file %q could not be read.", filename),
			})

			continue
		}

		if strings.HasSuffix(filename, ".json") {
			file, fileDiags = parser.ParseJSON(b, filename)
		} else {
			file, fileDiags = parser.ParseHCL(b, filename)
		}

		diags = append(diags, fileDiags...)

		if file == nil {
			continue
		}

		contentDiags := loadConfig(file, mod)

		diags = append(diags, contentDiags...)
	}

	return
}

func findConfigFiles(dir string) ([]string, error) {
	infos, err := os.ReadDir(dir)

	if err != nil {
		return []string{}, err
	}

	var primary []string
	var override []string

	for _, info := range infos {
		if info.IsDir() {
			continue
		}

		name := info.Name()

		if !isTerraformFile(name) || isIgnoredFile(name) {
			continue
		}

		baseName := strings.TrimSuffix(name, filepath.Ext(name))
		isOverride := baseName == "override" || strings.HasSuffix(baseName, "_override")

		fullPath := filepath.Join(dir, name)
		if isOverride {
			override = append(override, fullPath)
		} else {
			primary = append(primary, fullPath)
		}
	}

	primary = append(primary, override...)

	return primary, nil
}

func isTerraformFile(path string) bool {
	return strings.HasSuffix(path, ".tf") ||
		strings.HasSuffix(path, ".tf.json")
}

func isIgnoredFile(name string) bool {
	return strings.HasPrefix(name, ".") || // Unix-like hidden files
		strings.HasSuffix(name, "~") || // vim
		strings.HasPrefix(name, "#") && strings.HasSuffix(name, "#") // emacs
}
