package tfconfig

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestDirectoryFiles(t *testing.T) {
	files, _ := directoryFiles(afero.OsFs{}, "testdata/basic")

	expected := []string{
		"testdata/basic/main.tf",
		"testdata/basic/output.tf",
		"testdata/basic/variables.tf",
	}

	assert.Equal(t, expected, files)
}

func TestLoadModule(t *testing.T) {
	mod, _ := loadModule(afero.OsFs{}, "testdata/basic")

	assert.Equal(t, []string{"~> 1.3"}, mod.RequiredCore)
	assert.Equal(t, "aws.basic.bucket", mod.Backend.Attributes["bucket"])
	assert.Equal(t, "environments/nonprod/state.tfstate", mod.Backend.Attributes["key"])
}

func TestLoadModules(t *testing.T) {
	mods, _ := LoadModules(afero.OsFs{}, "testdata/complex")

	expectedModuleOne := &Module{
		Path:         "testdata/complex/environments/nonprod/project_one",
		RequiredCore: []string{"~> 1.3"},
		Backend: &Backend{
			Type: "s3",
			Attributes: map[string]string{
				"bucket": "aws.basic.bucket",
				"key":    "project_one/nonprod/state.tfstate",
			},
		},
	}

	expectedModuleTwo := &Module{
		Path:         "testdata/complex/environments/nonprod/project_two",
		RequiredCore: []string{"~> 1.3"},
		Backend: &Backend{
			Type: "s3",
			Attributes: map[string]string{
				"bucket": "aws.basic.bucket",
				"key":    "project_two/nonprod/state.tfstate",
			},
		},
	}

	assert.Equal(t, expectedModuleOne, &mods[0])
	assert.Equal(t, expectedModuleTwo, &mods[1])
}
