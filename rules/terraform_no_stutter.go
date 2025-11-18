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

	// Find all "resource" and "data" blocks
	body, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{Type: "resource", LabelNames: []string{"type", "name"}, Body: &hclext.BodySchema{}},
			{Type: "data", LabelNames: []string{"type", "name"}, Body: &hclext.BodySchema{}},
		},
	}, nil)

	if err != nil {
		return err
	}

	// For each block, find the longest common suffix between the block "type" (e.g. "gogle_storage_bucket") and the resource "name"
	for _, block := range body.Blocks {

		// block.Labels[0] contains the "type"
		// block.Labels[1] contains the "name"
		blockType := block.Labels[0]
		blockName := block.Labels[1]

		common := findLongestCommonSuffix(blockType, blockName)

		// We found common substring between the type and name
		if common != "" {

			message := fmt.Sprintf("Resource type (\"%s\") is repeated in resource name (\"%s\") (specifically \"%s\").", blockType, blockName, common)

			// Compute a recommended name for the resource.
			// It is just the current resource name minus the common prefix we found
			recommendedName, _ := strings.CutSuffix(blockName, fmt.Sprintf("_%s", common))

			// If the recommended name is the same as the current resource name, we can't automatically fix the issue. We only emit an error.
			if recommendedName == blockName {
				err := runner.EmitIssue(
					r,
					message,
					block.DefRange,
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
					block.DefRange,
					func(f tflint.Fixer) error {
						err := f.ReplaceText(block.LabelRanges[1], `"`, fmt.Sprintf(`"%s"`, recommendedName))
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
