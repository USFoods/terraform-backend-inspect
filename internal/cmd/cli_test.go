package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldFindWorkingDirs(t *testing.T) {
	dirs, _ := findWorkingDirs("testdata/terraform")

	expected := []string{
		"testdata/terraform",
		"testdata/terraform/environments",
		"testdata/terraform/environments/nonprod",
		"testdata/terraform/environments/nonprod/project_one",
		"testdata/terraform/environments/prod",
		"testdata/terraform/environments/prod/project_one",
	}

	assert.Equal(t, expected, dirs)
}
