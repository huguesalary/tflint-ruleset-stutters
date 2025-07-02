package rules

import (
	"fmt"
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// TerraformNoStutterRule checks whether ...
type TerraformNoStutterRule struct {
	tflint.DefaultRule
}

// NewTerraformNoStutterRule returns a new rule
func NewTerraformNoStutterRule() *TerraformNoStutterRule {
	return &TerraformNoStutterRule{}
}

// Name returns the rule name
func (r *TerraformNoStutterRule) Name() string {
	return "terraform_no_stutter"
}

// Enabled returns whether the rule is enabled by default
func (r *TerraformNoStutterRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *TerraformNoStutterRule) Severity() tflint.Severity {
	return tflint.NOTICE
}

// Link returns the rule reference link
func (r *TerraformNoStutterRule) Link() string {
	return ""
}

// Check checks whether ...
func (r *TerraformNoStutterRule) Check(runner tflint.Runner) error {

	// Find all "resource" blocks
	body, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{Type: "resource", LabelNames: []string{"type", "name"}, Body: &hclext.BodySchema{}},
		},
	}, nil)

	if err != nil {
		return err
	}

	// For each block, find the longest common suffix between the resource "type" (e.g. "gogle_storage_bucket") and the resource "name"
	for _, resource := range body.Blocks {

		// resource.Labels[0] contains the "type"
		// resource.Labels[1] contains the "name"
		common := findLongestCommonSuffix(resource.Labels[0], resource.Labels[1])

		// We found common substring between the type and name
		if common != "" {

			message := fmt.Sprintf("Resource type (\"%s\") is repeated in resource name (\"%s\") (specifically \"%s\").", resource.Labels[0], resource.Labels[1], common)

			// Compute a recommended name for the resource.
			// It is just the current resource name minus the common prefix we found
			recommendedName, _ := strings.CutSuffix(resource.Labels[1], fmt.Sprintf("_%s", common))

			// If the recommended name is the same as the current resource name, we can't automatically fix the issue. We only emit an error.
			if recommendedName == resource.Labels[1] {
				err := runner.EmitIssue(
					r,
					message,
					resource.DefRange,
				)

				if err != nil {
					return err
				}
			} else {

				// Add the recommended name to the recommandation message
				message = fmt.Sprintf("%s Recommended name: \"%s\".", message, recommendedName)

				err := runner.EmitIssueWithFix(
					r,
					message,
					resource.DefRange,
					func(f tflint.Fixer) error {
						err := f.ReplaceText(resource.LabelRanges[1], `"`, fmt.Sprintf(`"%s"`, recommendedName))
						if err != nil {
							return err
						}
						return nil
					},
				)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func findLongestCommonSuffix(a, b string) string {
	wordsA := strings.Split(a, "_")
	wordsB := strings.Split(b, "_")

	i := len(wordsA) - 1
	j := len(wordsB) - 1

	var commonSuffix []string

	for i >= 0 && j >= 0 && wordsA[i] == wordsB[j] {
		commonSuffix = append([]string{wordsA[i]}, commonSuffix...)
		i--
		j--
	}

	return strings.Join(commonSuffix, "_")
}
