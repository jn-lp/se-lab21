package gomodule

import (
	"fmt"
	"path"

	"github.com/google/blueprint"
	"github.com/roman-mazur/bood"
)

var (
	// Ninja rule to execute go test for generating test coverage report.
	goTestCoverage = pctx.StaticRule("coverage", blueprint.RuleParams{
		Command:     "cd $workDir && go test ${pkg} -coverprofile=$outputPath && go tool cover -html=$outputPath",
		Description: "generating test coverage of $pkg",
	}, "workDir", "pkg", "outputPath")
)

type testCoverageModule struct {
	blueprint.SimpleName

	properties struct {
        // Go package name to generate coverage
		Pkg string
		// List of source files.
		Srcs []string
		// Exclude patterns.
		SrcsExclude []string
	}
}

func (tb *testCoverageModule) GenerateBuildActions(ctx blueprint.ModuleContext) {
	name := ctx.ModuleName()
	config := bood.ExtractConfig(ctx)
	config.Debug.Printf("Adding build actions for generating '%s.html'", name)

	outputPath := path.Join(config.BaseOutputDir, fmt.Sprintf("reports/%s.html", name))

	inputs, testInputs, withError := patternsToPaths(ctx, tb.properties.Srcs, tb.properties.SrcsExclude)
	if withError {
		return
	}

	ctx.Build(pctx, blueprint.BuildParams{
		Description: fmt.Sprintf("Generating %s.html", name),
		Rule:        goTestCoverage,
		Outputs:     []string{outputPath},
		Implicits:   append(inputs, testInputs...),
		Args: map[string]string{
			"outputPath": outputPath,
			"workDir":    ctx.ModuleDir(),
			"pkg":        tb.properties.Pkg,
		},
	})
}

// TestCoverageFactory is a factory for test coverage generating blueprints
func TestCoverageFactory() (blueprint.Module, []interface{}) {
	mType := &testCoverageModule{}
	return mType, []interface{}{&mType.SimpleName.Properties, &mType.properties}
}
