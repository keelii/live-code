package main

import (
	_ "embed"
	"fmt"
	"github.com/evanw/esbuild/pkg/api"
	"log"
	"net/http"
)

//go:embed static/index.html
var indexHTML []byte

//go:embed static/preview.html
var previewHTML string

//go:embed static/codemirror-editor.js
var codemirrorEditorJS []byte

func main() {
	Get("/", func(req *LcRequest, res *LcResponse) {
		res.WriteHTML(indexHTML)
	})
	Get("/preview", func(req *LcRequest, res *LcResponse) {
		codeParam := req.URL.Query().Get("code")
		originCode := DecompressText(codeParam)
		ret := BuildReactCode(originCode, api.TransformOptions{})
		htmlString := ""
		if ret.Ok {
			htmlString = fmt.Sprintf(previewHTML, "Preivew", ret.Data)
		} else {
			errCode := fmt.Sprintf("alert(%s)", ToJsonString(ret.Msg))
			htmlString = fmt.Sprintf(previewHTML, "Error", errCode)
		}
		res.WriteHTML([]byte(htmlString))
	})

	Get("/static/codemirror-editor.js", func(req *LcRequest, res *LcResponse) {
		res.WriteJS(codemirrorEditorJS)
	})

	Post("/api/format", func(req *LcRequest, res *LcResponse) {
		code := req.GetParam("code")
		ret, err := DprintFormat("main.tsx", code.(string))
		if err != nil {
			fmt.Println("error:", err)
		} else {
			res.WriteJSON(OkResult(ret))
		}
	})
	Post("/api/build", func(req *LcRequest, res *LcResponse) {
		code := req.GetParam("code")
		ret := BuildReactCode(code.(string), api.TransformOptions{})

		if ret.Ok {
			res.WriteJSON(OkResult(ret.Data))
		} else {
			res.WriteJSON(ErrResult(ret.Msg))
		}
	})

	log.Println("Server started at http://localhost:8080")
	log.Fatalln(http.ListenAndServe(":8080", nil))
}
