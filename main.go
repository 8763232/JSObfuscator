package main

import (
	"fmt"
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func WalkPath(input string) []string {
	var files []string
	filepath.Walk(input, func(path string, info os.FileInfo, err error) error {
		if path == input {
			return nil
		}
		files = append(files, path)
		return nil
	})
	return files
}

func Encryption(o *Obfuscator, input, output string) {
	_input, err := os.ReadFile(input)
	if err != nil {
		log.Printf(`Encryption:%s`, err.Error())
		return
	}
	_output := o.Encryption(string(_input))
	if _output == "" {
		log.Printf(`Encryption: _output is nil`)
		return
	}

	if err := os.WriteFile(output, []byte(_output), 0644); err != nil {
		log.Printf(`Encryption: WriteFile :%s`, err.Error())
	}
}

func main() {

	config.WithOptions(config.ParseEnv)
	config.AddDriver(yaml.Driver)
	err := config.LoadFiles("./config.yml")
	if err != nil {
		panic(err)
	}

	params := config.String(`Params.obfuscate`, "{}")

	engine := NewEngine().SetParams(params)

	output := config.String(`Env.output`, "./output")
	input := config.String(`Env.input`, "./input")
	files := WalkPath(input)

	for _, f := range files {
		if IsDir(f) {
			continue
		}

		_output := strings.ReplaceAll(f, input, output)
		CreateDir(_output)
		Encryption(engine, f, _output)
		fmt.Printf(f)
	}
}
