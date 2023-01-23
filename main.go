package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/spf13/afero"
	"github.com/usfoods/terraform-backend-inspect/internal/tfconfig"
)

func main() {
	flag.Parse()

	var dir string
	if flag.NArg() > 0 {
		dir = flag.Arg(0)
	} else {
		dir = "."
	}

	modules, diagnostics := tfconfig.LoadModules(afero.OsFs{}, dir)

	if diagnostics.HasErrors() {
		os.Exit(1)
	}

	for _, mod := range modules {
		fmt.Printf("Module: %s\n", mod.Path)
		fmt.Printf("State File: %s/%s\n", mod.Backend.Attributes["bucket"], mod.Backend.Attributes["key"])
	}
}
