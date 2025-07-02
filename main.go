package main

import (
	"github.com/huguesalary/tflint-ruleset-stutters/rules"
	"github.com/terraform-linters/tflint-plugin-sdk/plugin"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		RuleSet: &tflint.BuiltinRuleSet{
			Name:    "stutters",
			Version: "0.1.0",
			Rules: []tflint.Rule{
				rules.NewTerraformNoStutterRule(),
			},
		},
	})
}
