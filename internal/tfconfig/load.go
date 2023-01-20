package tfconfig

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/spf13/afero"
)

func loadModule(fs afero.OsFs, dir string) *Module {
	mod := &Module{Path: dir}

	filesNames, diags := directoryFiles(fs, dir)

	parser := hclparse.NewParser()

	for _, filename := range filesNames {
		var file *hcl.File
		var fileDiags hcl.Diagnostics

		b, err := afero.ReadFile(fs, filename)

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

	return mod
}

func directoryFiles(fs afero.OsFs, dir string) (files []string, diags hcl.Diagnostics) {
	infos, err := afero.ReadDir(fs, dir)

	if err != nil {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Failed to read module directory",
			Detail:   fmt.Sprintf("Module directory %s does not exist or cannot be read.", dir),
		})
		return
	}

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

		if isOverride {
			continue
		}

		fullPath := filepath.Join(dir, name)

		files = append(files, fullPath)
	}

	return
}

func IsModuleDir(dir string) bool {
	files, _ := directoryFiles(afero.OsFs{}, dir)
	return len(files) != 0
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
