package cmd

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/usfoods/terraform-backend-inspect/internal/tfconfig"
)

// Exit codes are int values that represent an exit code for a particular error.
const (
	ExitCodeOK int = iota
	ExitCodeError
	ExitCodeIssuesFound
)

// CLI is the command line object
type CLI struct {
	// outStream and errStream are the stdout and stderr
	// to write message from the CLI.
	outStream, errStream io.Writer
	originalWorkingDir   string

	// // fields for each module
	formatter *Formatter
}

// NewCLI returns new CLI initialized by input streams
func NewCLI(outStream io.Writer, errStream io.Writer) (*CLI, error) {
	wd, err := os.Getwd()

	return &CLI{
		outStream:          outStream,
		errStream:          errStream,
		originalWorkingDir: wd,
	}, err
}

// Run invokes the CLI with the given arguments.
func (cli *CLI) Run() int {
	// Set up output formatter
	cli.formatter = &Formatter{
		Stdout: cli.outStream,
		Stderr: cli.errStream,
	}

	// find all directories under current to inspect
	workingDirs, err := findWorkingDirs(cli.originalWorkingDir)

	if err != nil {
		diags := hcl.Diagnostics{
			&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Failed to find working directories",
				Detail:   err.Error(),
			},
		}

		cli.formatter.Print(tfconfig.Issues{}, diags)

		return ExitCodeError
	}

	// load modules from the working directories
	modules, diags := tfconfig.LoadModules(workingDirs)

	if diags.HasErrors() {
		cli.formatter.Print(tfconfig.Issues{}, diags)
		return ExitCodeError
	}

	rules := []func([]*tfconfig.Module) tfconfig.Issues{
		tfconfig.ModuleShouldHaveRemoteBackend,
		tfconfig.ModuleShouldHaveLocalBackendOverride,
		tfconfig.ModuleShouldHaveUniqueBackend,
	}

	issues := tfconfig.ParseRules(modules, rules)

	if len(issues) > 0 {
		cli.formatter.Print(issues, nil)
		return ExitCodeIssuesFound
	}

	return ExitCodeOK
}

func findWorkingDirs(dir string) ([]string, error) {
	workingDirs := []string{}

	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			return nil
		}

		// making an assumption that modules are kept in a directory
		// that is aptly named "modules"
		if d.Name() == "modules" {
			return filepath.SkipDir
		}

		// hidden directories are skipped
		if path != "." && strings.HasPrefix(d.Name(), ".") {
			return filepath.SkipDir
		}

		workingDirs = append(workingDirs, path)
		return nil
	})

	if err != nil {
		return []string{}, err
	}

	return workingDirs, nil
}
