package main

import (
	"fmt"
	"github.com/evanw/esbuild/pkg/api"
	"strings"
)

type BuildResult = Result[string]

func okResult(data string) BuildResult {
	return BuildResult{
		Ok:   true,
		Msg:  "",
		Data: data,
	}
}
func errResult(errs []api.Message) BuildResult {
	msg := make([]string, 0)
	for _, e := range errs {
		msg = append(msg, e.Text)
	}

	return BuildResult{
		Ok:   false,
		Msg:  strings.Join(msg, "\n"),
		Data: "",
	}
}

func Build(code string, options api.TransformOptions) BuildResult {
	config := api.TransformOptions{
		Loader:    api.LoaderTSX,
		Format:    api.FormatIIFE,
		Sourcemap: api.SourceMapInline,
	}

	if options.Sourcemap != config.Sourcemap {
		config.Sourcemap = options.Sourcemap
	}

	ret := api.Transform(code, config)

	if len(ret.Errors) > 0 {
		return errResult(ret.Errors)
	} else {
		return okResult(string(ret.Code))
	}
}
func BuildReactCode(code string, options api.TransformOptions) BuildResult {
	if code == "" {
		code = `function App() { return <div>No code to run...</div> }`
	}

	reactRuntime := fmt.Sprintf(`
try {
  %s
  window._REACT_ROOT_ = ReactDOM.createRoot(document.getElementById("root"))
  window._REACT_ROOT_.render(<App />)
} catch (e) { 
  console.error(e)
  if (e instanceof ReferenceError && e.message.includes("App is not defined")) {
    alert("App component is not defined, please define it like this:\n" + "function App() {\n    return <div>Hello, live-code.</div>\n}")	
  } else {
    alert(e); 
  }
}
`, code)

	return Build(reactRuntime, options)
}
