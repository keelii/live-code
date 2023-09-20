package main

import (
	"github.com/evanw/esbuild/pkg/api"
	"testing"
)

func TestBuild(t *testing.T) {
	code := `var a: number = 1`
	result := Build(code, api.TransformOptions{
		Sourcemap: api.SourceMapNone,
	})

	if !result.Ok {
		t.Error(result.Msg)
	} else {
		if result.Data != "(() => {\n})();\n" {
			t.Error("build failed")
		} else {
			t.Log("build success")
		}
	}
}
