package tfconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindconfigFiles(t *testing.T) {
	files, _ := findConfigFiles("testdata/basic")

	expected := []string{
		"testdata/basic/main.tf",
		"testdata/basic/output.tf",
		"testdata/basic/variables.tf",
		"testdata/basic/main_override.tf",
	}

	assert.Equal(t, expected, files)
}

func TestLoadModule(t *testing.T) {
	files := []string{
		"testdata/basic/main.tf",
		"testdata/basic/output.tf",
		"testdata/basic/variables.tf",
		"testdata/basic/main_override.tf",
	}

	mod, _ := loadModule(files, "testdata/basic")

	assert.Equal(t, []string{"~> 1.3"}, mod.RequiredCore)
	assert.Equal(t, "aws.basic.bucket", mod.Backends[0].Attributes["bucket"])
	assert.Equal(t, "environments/nonprod/state.tfstate", mod.Backends[0].Attributes["key"])
}

func TestLoadModules(t *testing.T) {
	mods, _ := LoadModules([]string{
		"testdata/complex",
		"testdata/complex/environments",
		"testdata/complex/environments/nonprod",
		"testdata/complex/environments/nonprod/project_one",
		"testdata/complex/environments/nonprod/project_two",
	})

	expectedModuleOne := &Module{
		Path:         "testdata/complex/environments/nonprod/project_one",
		RequiredCore: []string{"~> 1.3"},
		Backends: []*Backend{
			{
				Type: "s3",
				Attributes: map[string]string{
					"bucket": "aws.basic.bucket",
					"key":    "project_one/nonprod/state.tfstate",
				},
			},
			{
				Type: "local",
				Attributes: map[string]string{
					"path": "integration.tfstate",
				},
			},
		},
	}

	expectedModuleTwo := &Module{
		Path:         "testdata/complex/environments/nonprod/project_two",
		RequiredCore: []string{"~> 1.3"},
		Backends: []*Backend{
			{
				Type: "s3",
				Attributes: map[string]string{
					"bucket": "aws.basic.bucket",
					"key":    "project_two/nonprod/state.tfstate",
				},
			},
			{
				Type: "local",
				Attributes: map[string]string{
					"path": "integration.tfstate",
				},
			},
		},
	}

	assert.Equal(t, expectedModuleOne, mods[0])
	assert.Equal(t, expectedModuleTwo, mods[1])
}
