package gomodule


import (
	"regexp"
	"github.com/google/blueprint"
)

var (
	// Package context used to define Ninja build rules.
	pctx = blueprint.NewPackageContext("github.com/jn-lp/se-lab21/build/gomodule")
)

func patternsToPaths(ctx blueprint.ModuleContext, include []string, exclude []string) (inputsPaths []string, testsPaths []string, withError bool) {
	for _, src := range include {
		if matches, err := ctx.GlobWithDeps(src, exclude); err == nil {
			for _, path := range matches {
				if isTest, _ := regexp.Match("^.*_test.go$", []byte(path)); isTest {
					testsPaths = append(testsPaths, path)
				} else {
					inputsPaths = append(inputsPaths, path)
				}
			}
		} else {
			ctx.PropertyErrorf("srcs", "Cannot resolve files that match pattern %s", src)
			return nil, nil, true
		}
	}
	return
}
