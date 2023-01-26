package tfconfig

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/assert"
)

func TestShouldRaiseIssueOnDuplicateBackend(t *testing.T) {

	modules, diags := LoadModules([]string{
		"testdata/rule_no_duplicates/nonprod/project_one",
		"testdata/rule_no_duplicates/nonprod/project_two",
		"testdata/rule_no_duplicates/prod/project_one",
		"testdata/rule_no_duplicates/prod/project_two",
	})

	if diags.HasErrors() {
		t.Fatalf("error loading module: %s", diags.Error())
	}

	issues := ModuleShouldHaveUniqueBackend(modules)

	expectedIssues := Issues{
		{
			Message:    "Duplicate remote backend configuration",
			ModulePath: "testdata/rule_no_duplicates/nonprod/project_one",
			Range: &hcl.Range{
				Filename: "testdata/rule_no_duplicates/nonprod/project_one/main.tf",
				Start:    hcl.Pos{Line: 3, Column: 3},
				End:      hcl.Pos{Line: 3, Column: 15},
			},
			Severity: ERROR,
		},
		{
			Message:    "Duplicate remote backend configuration",
			ModulePath: "testdata/rule_no_duplicates/prod/project_one",
			Range: &hcl.Range{
				Filename: "testdata/rule_no_duplicates/prod/project_one/main.tf",
				Start:    hcl.Pos{Line: 3, Column: 3},
				End:      hcl.Pos{Line: 3, Column: 15},
			},
			Severity: ERROR,
		},
		{
			Message:    "Duplicate remote backend configuration",
			ModulePath: "testdata/rule_no_duplicates/nonprod/project_two",
			Range: &hcl.Range{
				Filename: "testdata/rule_no_duplicates/nonprod/project_two/main.tf",
				Start:    hcl.Pos{Line: 3, Column: 3},
				End:      hcl.Pos{Line: 3, Column: 15},
			},
			Severity: ERROR,
		},
		{
			Message:    "Duplicate remote backend configuration",
			ModulePath: "testdata/rule_no_duplicates/prod/project_two",
			Range: &hcl.Range{
				Filename: "testdata/rule_no_duplicates/prod/project_two/main.tf",
				Start:    hcl.Pos{Line: 3, Column: 3},
				End:      hcl.Pos{Line: 3, Column: 15},
			},
			Severity: ERROR,
		},
	}

	opts := []cmp.Option{cmpopts.IgnoreFields(hcl.Pos{}, "Byte")}

	assert.True(t, cmp.Equal(expectedIssues, issues, opts...), cmp.Diff(expectedIssues, issues, opts...))
}

func TestShouldRaiseIssueOnLocalBackend(t *testing.T) {

	modules, diags := LoadModules([]string{
		"testdata/rule_no_locals/nonprod/project_one",
		"testdata/rule_no_locals/nonprod/project_two",
		"testdata/rule_no_locals/prod/project_one",
		"testdata/rule_no_locals/prod/project_two",
	})

	if diags.HasErrors() {
		t.Fatalf("error loading module: %s", diags.Error())
	}

	issues := ModuleShouldHaveRemoteBackend(modules)

	expectedIssues := Issues{
		{
			Message:    "No remote backend configured",
			ModulePath: "testdata/rule_no_locals/nonprod/project_one",
			Range: &hcl.Range{
				Filename: "testdata/rule_no_locals/nonprod/project_one/main.tf",
				Start:    hcl.Pos{Line: 3, Column: 3},
				End:      hcl.Pos{Line: 3, Column: 18},
			},
			Severity: ERROR,
		},
		{
			Message:    "No remote backend configured",
			ModulePath: "testdata/rule_no_locals/prod/project_one",
			Range:      nil,
			Severity:   ERROR,
		},
	}

	opts := []cmp.Option{cmpopts.IgnoreFields(hcl.Pos{}, "Byte")}

	assert.True(t, cmp.Equal(expectedIssues, issues, opts...), cmp.Diff(expectedIssues, issues, opts...))
}

func TestShouldRaiseIssueOnNoLocalOverride(t *testing.T) {

	modules, diags := LoadModules([]string{
		"testdata/rule_local_overrides/nonprod/project_one",
		"testdata/rule_local_overrides/nonprod/project_two",
		"testdata/rule_local_overrides/prod/project_one",
		"testdata/rule_local_overrides/prod/project_two",
	})

	if diags.HasErrors() {
		t.Fatalf("error loading module: %s", diags.Error())
	}

	issues := ModuleShouldHaveLocalBackendOverride(modules)

	expectedIssues := Issues{
		{
			Message:    "No local backend override configured",
			ModulePath: "testdata/rule_local_overrides/prod/project_one",
			Range:      nil,
			Severity:   ERROR,
		},
		{
			Message:    "No local backend override configured",
			ModulePath: "testdata/rule_local_overrides/prod/project_two",
			Range:      nil,
			Severity:   ERROR,
		},
	}

	opts := []cmp.Option{cmpopts.IgnoreFields(hcl.Pos{}, "Byte")}

	assert.True(t, cmp.Equal(expectedIssues, issues, opts...), cmp.Diff(expectedIssues, issues, opts...))
}

func TestShouldRaiseIssueNoRemoteWithOverride(t *testing.T) {
	modules, diags := LoadModules([]string{
		"testdata/local_with_override",
	})

	if diags.HasErrors() {
		t.Fatalf("error loading module: %s", diags.Error())
	}

	issues := ModuleShouldHaveRemoteBackend(modules)

	expectedIssues := Issues{
		{
			Message:    "No remote backend configured",
			ModulePath: "testdata/local_with_override",
			Range:      nil,
			Severity:   ERROR,
		},
	}

	opts := []cmp.Option{cmpopts.IgnoreFields(hcl.Pos{}, "Byte")}

	assert.True(t, cmp.Equal(expectedIssues, issues, opts...), cmp.Diff(expectedIssues, issues, opts...))
}
