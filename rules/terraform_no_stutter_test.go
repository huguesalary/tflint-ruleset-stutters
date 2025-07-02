package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformNoStutter(t *testing.T) {
	tests := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "issue found with recommended name",
			Content: `
resource "google_storage_bucket" "super_duper_important_storage_bucket" {}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformNoStutterRule(),
					Message: "Resource type (\"google_storage_bucket\") is repeated in resource name (\"super_duper_important_storage_bucket\") (specifically \"storage_bucket\"). Recommended name: \"super_duper_important\".",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 72},
					},
				},
			},
		},
		{
			Name: "issue found without recommended name",
			Content: `
resource "google_storage_bucket" "storage_bucket" {}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformNoStutterRule(),
					Message: "Resource type (\"google_storage_bucket\") is repeated in resource name (\"storage_bucket\") (specifically \"storage_bucket\").",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 50},
					},
				},
			},
		},
	}

	rule := NewTerraformNoStutterRule()

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"resource.tf": test.Content})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssues(t, test.Expected, runner.Issues)
		})
	}
}
