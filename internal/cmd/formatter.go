package cmd

import (
	"fmt"
	"io"

	"github.com/fatih/color"
	"github.com/hashicorp/hcl/v2"
	"github.com/usfoods/terraform-backend-inspect/internal/tfconfig"
)

// Formatter outputs results to stdout and stderr
type Formatter struct {
	Stdout io.Writer
	Stderr io.Writer
}

var colorBold = color.New(color.Bold).SprintfFunc()
var colorError = color.New(color.FgRed).SprintFunc()
var colorWarning = color.New(color.FgYellow).SprintFunc()
var colorNotice = color.New(color.FgHiWhite).SprintFunc()

// Print outputs the given issues and errors according to configured format
func (f *Formatter) Print(issues tfconfig.Issues, diags hcl.Diagnostics) {
	if len(issues) > 0 {
		fmt.Fprintf(f.Stdout, "%d issue(s) found:\n\n", len(issues))

		for _, issue := range issues {
			fmt.Fprintf(
				f.Stdout,
				"%s: %s \n\n",
				colorSeverity(issue.Severity), colorBold(issue.Message),
			)
		}
	}

	if diags.HasErrors() {
		for _, err := range diags {
			fmt.Fprintf(f.Stderr, "%s: %s\n\n", colorSeverity(tfconfig.ERROR), err.Summary)
		}
	}
}

func colorSeverity(severity tfconfig.Severity) string {
	switch severity {
	case tfconfig.ERROR:
		return colorError(severity)
	case tfconfig.WARNING:
		return colorWarning(severity)
	case tfconfig.NOTICE:
		return colorNotice(severity)
	default:
		panic("Unreachable")
	}
}
