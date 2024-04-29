package main

import (
	"log"
	"os"
	"strings"
)
import v8 "rogchap.com/v8go"

type Obfuscator struct {
	params     string
	runtime    *v8.Context
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
	vm := v8.NewContext()
	if _, err := vm.RunScript(obfuscator, "obfuscator.js"); err != nil {
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
	JavaScript = strings.ReplaceAll(JavaScript, "\n", "\\\n")
	JavaScript = strings.ReplaceAll(JavaScript, "'", "\"")
	source := `Encryption('` + JavaScript + `',` + o.params + `)`

	res, err := o.runtime.RunScript(source, "Encryption.js")
	if err != nil {
		log.Printf(`Encryption error:%s`, err.Error())
		return ""
	}
	return res.String()
}
