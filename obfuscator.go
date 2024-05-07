package main

import (
	js "github.com/dop251/goja"
	"log"
	"os"
	"strconv"
)

type Obfuscator struct {
	params     string
	runtime    *js.Runtime
	obfuscator []byte
}

func NewEngine() *Obfuscator {

	// 读取加密核心引擎
	_obfuscator, err := os.ReadFile(`./index.browser.js`)
	if err != nil {
		panic(err)
	}
	// 加入模板
	template := `
	function Encryption(js,params) {
		var obfuscationResult = JavaScriptObfuscator.obfuscate(js,params);
		return obfuscationResult.getObfuscatedCode();
	}`
	obfuscator := string(_obfuscator) + ";" + template

	// 初始化引擎
	vm := js.New()
	if _, err := vm.RunScript("obfuscator.js", obfuscator); err != nil {
		panic(err)
	}

	obj := &Obfuscator{
		runtime:    vm,
		obfuscator: []byte(obfuscator),
	}

	return obj
}

func (o *Obfuscator) SetParams(params string) *Obfuscator {
	o.params = params
	return o
}

func (o *Obfuscator) Encryption(JavaScript string) string {
	//JavaScript = strings.ReplaceAll(JavaScript, "\n", "\\\n")
	//JavaScript = strings.ReplaceAll(JavaScript, "'", "\"")
	source := `Encryption(` + strconv.Quote(JavaScript) + `,` + o.params + `)`

	res, err := o.runtime.RunScript("Encryption.js", source)
	if err != nil {
		log.Printf(`Encryption error:%s`, err.Error())
		return ""
	}
	return res.String()
}
