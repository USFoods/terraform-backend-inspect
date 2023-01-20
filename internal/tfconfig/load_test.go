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
	mod := loadModule(afero.OsFs{}, "testdata/basic")

	assert.Equal(t, []string{"~> 1.3"}, mod.RequiredCore)
	assert.Equal(t, "aws.basic.bucket", mod.Backend.Attributes["bucket"])
}
