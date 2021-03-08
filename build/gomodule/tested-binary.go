package gomodule

import (
	"fmt"
	"path"

	"github.com/google/blueprint"
	"github.com/roman-mazur/bood"
)

var (
	// Ninja rule to execute go build.
	goBuild = pctx.StaticRule("binaryBuild", blueprint.RuleParams{
		Command:     "cd $workDir && go build -o $outputPath $pkg",
		Description: "build go command $pkg",
	}, "workDir", "outputPath", "pkg")

	// Ninja rule to execute go mod vendor.
	goVendor = pctx.StaticRule("vendor", blueprint.RuleParams{
		Command:     "cd $workDir && go mod vendor",
		Description: "vendor dependencies of $name",
	}, "workDir", "name")

	// Ninja rule to execute go test.
	goTest = pctx.StaticRule("test", blueprint.RuleParams{
		Command:     "cd $workDir && go test -v $testPkg > $reportPath",
		Description: "test $testPkg",
	}, "workDir", "reportPath", "testPkg")
)

type testedBinaryModule struct {
	blueprint.SimpleName

	properties struct {
		// Go package name to build as a command with "go build".
		Pkg string
		// Go package name to test as a command with "go test".
		TestPkg string
		// List of source files.
		Srcs []string
		// Exclude patterns.
		SrcsExclude []string
		// If to call vendor command.
		VendorFirst bool
	}
}

func (tb *testedBinaryModule) GenerateBuildActions(ctx blueprint.ModuleContext) {
	name := ctx.ModuleName()
	config := bood.ExtractConfig(ctx)
	config.Debug.Printf("Adding build actions for go binary module '%s'", name)

	outputPath := path.Join(config.BaseOutputDir, "bin", name)
	// reportPath := path.Join(config.BaseOutputDir, fmt.Sprintf("%s-%v.txt", name, time.Now().UnixNano()))
	reportPath := path.Join(config.BaseOutputDir, "report.log")

	inputs, testInputs, withError := patternsToPaths(ctx, tb.properties.Srcs, tb.properties.SrcsExclude)
	if withError {
		return
	}

	if tb.properties.VendorFirst {
		vendorDirPath := path.Join(ctx.ModuleDir(), "vendor")
		ctx.Build(pctx, blueprint.BuildParams{
			Description: fmt.Sprintf("Vendor dependencies of %s", name),
			Rule:        goVendor,
			Outputs:     []string{vendorDirPath},
			Implicits:   []string{path.Join(ctx.ModuleDir(), "go.mod")},
			Optional:    true,
			Args: map[string]string{
				"workDir": ctx.ModuleDir(),
				"name":    name,
			},
		})
		inputs = append(inputs, vendorDirPath)
	}

	if len(tb.properties.TestPkg) > 0 {
		ctx.Build(pctx, blueprint.BuildParams{
			Description: fmt.Sprintf("Test module %s", tb.properties.TestPkg),
			Rule:        goTest,
			Outputs:     []string{reportPath},
			Implicits:   append(testInputs, inputs...),
			Args: map[string]string{
				"reportPath": reportPath,
				"workDir":    ctx.ModuleDir(),
				"testPkg":    tb.properties.TestPkg,
			},
		})
	}

	ctx.Build(pctx, blueprint.BuildParams{
		Description: fmt.Sprintf("Build %s as Go binary", name),
		Rule:        goBuild,
		Outputs:     []string{outputPath},
		Implicits:   inputs,
		Args: map[string]string{
			"outputPath": outputPath,
			"workDir":    ctx.ModuleDir(),
			"pkg":        tb.properties.Pkg,
		},
	})
}

// TestedBinFactory is a factory for go tested binary module type which supports Go command packages with running tests.
func TestedBinFactory() (blueprint.Module, []interface{}) {
	mType := &testedBinaryModule{}
	return mType, []interface{}{&mType.SimpleName.Properties, &mType.properties}
}
