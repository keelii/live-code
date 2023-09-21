package main

import (
	"dario.cat/mergo"
	_ "embed"
	"github.com/dop251/goja"
	"log"
)

//go:embed static/prettier.js
var prettierJs string

var polyfill = `
function RangeError() {}
`
var runtimeFunction = `
function format(str, opt) {
	return prettier.format(str, JSON.parse(opt))
}
`

type PrettierFormatOptions struct {
	UseTabs       bool   `json:"useTabs"`
	TabWidth      int    `json:"tabWidth"`
	PrintWidth    int    `json:"printWidth"`
	SingleQuote   bool   `json:"singleQuote"`
	TrailingComma string `json:"trailingComma"`
	Semi          bool   `json:"semi"`
}

var format func(code string, optJsonString string) string

var defaults = PrettierFormatOptions{
	UseTabs:       false,
	TabWidth:      4,
	PrintWidth:    80,
	SingleQuote:   false,
	TrailingComma: "none",
	Semi:          false,
}

func PrettierFormat(code string, opts PrettierFormatOptions) (string, error) {
	if err := mergo.Map(&opts, defaults); err != nil {
		log.Fatalln("mergo defaults error:", err)
		return code, err
	}
	return format(code, ToJsonString(opts)), nil
}

func init() {
	vm := goja.New()
	_, err := vm.RunString(polyfill + prettierJs + runtimeFunction)
	if err != nil {
		log.Fatalln("prettierJs error:", err)
	}

	err = vm.ExportTo(vm.Get("format"), &format)
	if err != nil {
		panic(err)
	}

	//fmt.Println(PrettierFormat("var a= 1;var b: string = <div>1\n\t\t</div>", PrettierFormatOptions{
	//}))
	//fmt.Println(PrettierFormat(`{'a':1 }`, PrettierFormatOptions{}))
}
