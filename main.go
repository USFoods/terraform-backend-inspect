package main

import (
	"fmt"
	"os"

	"github.com/mattn/go-colorable"
	"github.com/usfoods/terraform-backend-inspect/internal/cmd"
)

func main() {
	cli, err := cmd.NewCLI(colorable.NewColorable(os.Stdout), colorable.NewColorable(os.Stderr))

	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(cmd.ExitCodeError)
	}

	os.Exit(cli.Run())

	// modIssues := map[string][]string{}

	// modules, diagnostics := tfconfig.LoadModules(afero.OsFs{}, dir)

	// if diagnostics.HasErrors() {
	// 	for _, err := range diagnostics.Errs() {
	// 		print(err.Error())
	// 	}

	// 	os.Exit(1)
	// }

	// stateFiles := map[string][]string{}

	// for _, mod := range modules {
	// 	backend := mod.Backend
	// 	modPath := strings.ReplaceAll(mod.Path, dir, "")

	// 	if backend == nil || backend.Type == "local" {
	// 		modIssues[modPath] = append(modIssues[modPath], "Found local backend configuration")
	// 		continue
	// 	}

	// 	attrs := backend.Attributes
	// 	statePath := fmt.Sprintf("%s/%s", attrs["bucket"], attrs["key"])

	// 	stateFiles[statePath] = append(stateFiles[statePath], modPath)
	// }

	// for _, modules := range stateFiles {
	// 	if len(modules) > 1 {
	// 		for _, module := range modules {
	// 			modIssues[module] = append(modIssues[module], "Found duplicate backend configuration")
	// 		}
	// 	}
	// }

	// for key, issues := range modIssues {
	// 	for _, issue := range issues {
	// 		fmt.Printf("%s: Warning - %s\n", key, issue)
	// 	}
	// }

	// if len(modIssues) > 0 {
	// 	os.Exit(2)
	// }
}
