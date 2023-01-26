package tfconfig

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/hcl/v2"
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
	assert.Equal(t, "basic/state.tfstate", mod.Backends[0].Attributes["key"])
}

func TestLoadModules(t *testing.T) {
	mods, _ := LoadModules([]string{
		"testdata/complex",
		"testdata/complex/nonprod",
		"testdata/complex/nonprod/project_one",
		"testdata/complex/nonprod/project_two",
		"testdata/complex/prod",
		"testdata/complex/prod/project_one",
		"testdata/complex/prod/project_two",
	})

	expectedModuleOne := &Module{
		Path:         "testdata/complex/nonprod/project_one",
		RequiredCore: []string{"~> 1.3"},
		Backends: []*Backend{
			{
				Type: "s3",
				Attributes: map[string]string{
					"bucket": "aws.basic.bucket",
					"key":    "project_one/nonprod/state.tfstate",
				},
				Range: hcl.Range{
					Filename: "testdata/complex/nonprod/project_one/main.tf",
					Start:    hcl.Pos{Line: 3, Column: 3},
					End:      hcl.Pos{Line: 3, Column: 15},
				},
			},
			{
				Type: "local",
				Attributes: map[string]string{
					"path": "integration.tfstate",
				},
				Range: hcl.Range{
					Filename: "testdata/complex/nonprod/project_one/main_override.tf",
					Start:    hcl.Pos{Line: 2, Column: 3},
					End:      hcl.Pos{Line: 2, Column: 18},
				},
			},
		},
	}

	expectedModuleTwo := &Module{
		Path:         "testdata/complex/nonprod/project_two",
		RequiredCore: []string{"~> 1.3"},
		Backends: []*Backend{
			{
				Type: "s3",
				Attributes: map[string]string{
					"bucket": "aws.basic.bucket",
					"key":    "project_two/nonprod/state.tfstate",
				},
				Range: hcl.Range{
					Filename: "testdata/complex/nonprod/project_two/main.tf",
					Start:    hcl.Pos{Line: 3, Column: 3},
					End:      hcl.Pos{Line: 3, Column: 15},
				},
			},
			{
				Type: "local",
				Attributes: map[string]string{
					"path": "integration.tfstate",
				},
				Range: hcl.Range{
					Filename: "testdata/complex/nonprod/project_two/main_override.tf",
					Start:    hcl.Pos{Line: 2, Column: 3},
					End:      hcl.Pos{Line: 2, Column: 18},
				},
			},
		},
	}

	expectedModuleThree := &Module{
		Path:         "testdata/complex/prod/project_one",
		RequiredCore: []string{"~> 1.3"},
		Backends: []*Backend{
			{
				Type: "s3",
				Attributes: map[string]string{
					"bucket": "aws.basic.bucket",
					"key":    "project_one/prod/state.tfstate",
				},
				Range: hcl.Range{
					Filename: "testdata/complex/prod/project_one/main.tf",
					Start:    hcl.Pos{Line: 3, Column: 3},
					End:      hcl.Pos{Line: 3, Column: 15},
				},
			},
			{
				Type: "local",
				Attributes: map[string]string{
					"path": "integration.tfstate",
				},
				Range: hcl.Range{
					Filename: "testdata/complex/prod/project_one/main_override.tf",
					Start:    hcl.Pos{Line: 2, Column: 3},
					End:      hcl.Pos{Line: 2, Column: 18},
				},
			},
		},
	}

	expectedModuleFour := &Module{
		Path:         "testdata/complex/prod/project_two",
		RequiredCore: []string{"~> 1.3"},
		Backends: []*Backend{
			{
				Type: "s3",
				Attributes: map[string]string{
					"bucket": "aws.basic.bucket",
					"key":    "project_two/prod/state.tfstate",
				},
				Range: hcl.Range{
					Filename: "testdata/complex/prod/project_two/main.tf",
					Start:    hcl.Pos{Line: 3, Column: 3},
					End:      hcl.Pos{Line: 3, Column: 15},
				},
			},
			{
				Type: "local",
				Attributes: map[string]string{
					"path": "integration.tfstate",
				},
				Range: hcl.Range{
					Filename: "testdata/complex/prod/project_two/main_override.tf",
					Start:    hcl.Pos{Line: 2, Column: 3},
					End:      hcl.Pos{Line: 2, Column: 18},
				},
			},
		},
	}

	opts := []cmp.Option{cmpopts.IgnoreFields(hcl.Pos{}, "Byte")}

	assert.True(t, cmp.Equal(expectedModuleOne, mods[0], opts...), cmp.Diff(expectedModuleOne, mods[0], opts...))
	assert.True(t, cmp.Equal(expectedModuleTwo, mods[1], opts...), cmp.Diff(expectedModuleTwo, mods[1], opts...))
	assert.True(t, cmp.Equal(expectedModuleThree, mods[2], opts...), cmp.Diff(expectedModuleThree, mods[2], opts...))
	assert.True(t, cmp.Equal(expectedModuleFour, mods[3], opts...), cmp.Diff(expectedModuleFour, mods[3], opts...))
}
