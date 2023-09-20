package main

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wasmerio/wasmer-go/wasmer"
	"log"
)

type GlobalConfiguration struct {
	LineWidth   int    `json:"lineWidth"`
	IndentWidth int    `json:"indentWidth"`
	UseTabs     bool   `json:"useTabs"`
	NewLineKind string `json:"newLineKind"`
	// "auto" | "lf" | "crlf" | "system"
}
type PluginConfig struct {
	// "prefer" | "asi"
	SemiColons string `json:"semiColons"`
	// preferSingle | alwaysDouble
	QuoteStyle string `json:"quoteStyle"`
}

var instance *wasmer.Instance

//go:embed static/typescript-0.84.4.wasm
var tsWasmBytes []byte

func RunFunctionMust(name string) interface{} {
	fn := GetFunctionMust(name)
	ret, err := fn()
	if err != nil {
		log.Fatalln("error RunFunctionMust:", name, err)
	}
	return ret
}
func GetFunctionMust(name string) wasmer.NativeFunction {
	fn, err := instance.Exports.GetFunction(name)
	if err != nil {
		log.Fatalln("error GetFunctionMust:", name, err)
	}
	return fn
}
func SetFormatConfig(globalConfig GlobalConfiguration, pluginConfig PluginConfig) {
	setConfig(globalConfig, pluginConfig)
}

func DprintFormat(fileName string, fileText string) (string, error) {
	SetFormatConfig(GlobalConfiguration{}, PluginConfig{})

	setFilePath := GetFunctionMust("set_file_path")

	format := GetFunctionMust("format")

	resetConfig := GetFunctionMust("reset_config")
	if resetConfig != nil {
		resetConfig()
	}

	//setOverrideConfig()

	sendString(fileName)
	setFilePath()
	sendString(fileText)
	code, err := format()

	if err != nil {
		fmt.Println("format error: ", err, code)
		return "", err
	}

	switch int(code.(int32)) {
	case 0:
		return fileText, nil
	case 1:
		return recvString(RunFunctionMust("get_formatted_text").(int32)), nil
	case 2:
		return "", errors.New(recvString(RunFunctionMust("get_error_text").(int32)))
	default:
		return "", errors.New("unknown error")
	}
}

func setConfig(globalConfig GlobalConfiguration, pluginConfig PluginConfig) {
	setGlobalConfig := GetFunctionMust("set_global_config")
	setPluginConfig := GetFunctionMust("set_plugin_config")

	if globalConfig.LineWidth == 0 {
		globalConfig.LineWidth = 80
	}
	if globalConfig.IndentWidth == 0 {
		globalConfig.IndentWidth = 4
	}
	if globalConfig.NewLineKind == "" {
		globalConfig.NewLineKind = "auto"
	}

	if pluginConfig.QuoteStyle == "" {
		pluginConfig.QuoteStyle = "preferSingle"
	}
	if pluginConfig.SemiColons == "" {
		pluginConfig.SemiColons = "asi"
	}

	buf1, err := json.Marshal(globalConfig)
	if err != nil {
		log.Println("Marshal globalConfig error", err)
		return
	}

	buf2, err := json.Marshal(pluginConfig)
	if err != nil {
		log.Println("Marshal pluginConfig error", err)
		return
	}

	log.Println("globalConfig", string(buf1))
	log.Println("pluginConfig", string(buf2))

	sendString(string(buf1))
	setGlobalConfig()
	sendString(string(buf2))
	setPluginConfig()
}
func getWasmBuffer(length int32) []byte {
	getWasmMemoryBuffer := GetFunctionMust("get_wasm_memory_buffer")
	pointer, _ := getWasmMemoryBuffer()
	offset := pointer.(int32)
	mem, e := instance.Exports.GetMemory("memory")
	if e != nil {
		log.Fatalln("error getWasmBuffer:", e)
	}
	return mem.Data()[offset : offset+length]
}
func sendString(text string) {
	fmt.Println("sendString:", text)
	clearSharedBytes := GetFunctionMust("clear_shared_bytes")
	addToSharedBytesFromBuffer := GetFunctionMust("add_to_shared_bytes_from_buffer")
	getWasmMemoryBufferSize := GetFunctionMust("get_wasm_memory_buffer_size")
	bufferSize, _ := getWasmMemoryBufferSize()
	bufferIntSize := bufferSize.(int32)

	encodedText := []byte(text)
	length := int32(len(encodedText))

	clearSharedBytes(length)

	var index int32
	for index = 0; index < length; index++ {
		writeCount := bufferIntSize
		if length-index < bufferIntSize {
			writeCount = length - index
		}
		wasmBuffer := getWasmBuffer(writeCount)
		var i int32
		for i = 0; i < writeCount; i++ {
			wasmBuffer[i] = encodedText[index+i]
		}
		addToSharedBytesFromBuffer(writeCount)
		index += writeCount
	}
}
func recvString(length int32) string {
	setBufferWithSharedBytes := GetFunctionMust("set_buffer_with_shared_bytes")
	getWasmMemoryBufferSize := GetFunctionMust("get_wasm_memory_buffer_size")
	bufferSize, _ := getWasmMemoryBufferSize()
	bufferIntSize := bufferSize.(int32)

	buffer := make([]byte, length)

	var index int32
	for index = 0; index < length; index++ {
		readCount := bufferIntSize
		if length-index < bufferIntSize {
			readCount = length - index
		}
		setBufferWithSharedBytes(index, readCount)
		wasmBuffer := getWasmBuffer(readCount)

		var i int32
		for i = 0; i < readCount; i++ {
			buffer[index+i] = wasmBuffer[i]
		}

		index += readCount
	}
	return string(buffer)
}

func createInstance(store *wasmer.Store, wasmBytes []byte) *wasmer.Instance {
	// 编译模块
	module, err := wasmer.NewModule(store, wasmBytes)
	if err != nil {
		log.Fatalln("failed to compile module:", err)
	}

	// 实例化模块
	inst, err := wasmer.NewInstance(module, wasmer.NewImportObject())
	if err != nil {
		log.Fatalln("failed to instantiate module:", err)
	}
	return inst
}

func init() {
	engine := wasmer.NewEngine()
	store := wasmer.NewStore(engine)

	instance = createInstance(store, tsWasmBytes)
}
