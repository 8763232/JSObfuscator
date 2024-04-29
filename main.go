package main

import (
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
	"log"
)

func main() {

	config.WithOptions(config.ParseEnv)
	config.AddDriver(yaml.Driver)
	err := config.LoadFiles("./config.yml")
	if err != nil {
		panic(err)
	}

	JavaScript := `
        (function(){
            var variable1 = '5' - 3;
            var variable2 = '5' + 3;
            var variable3 = '5' + - '2';
            var variable4 = ['10','10','10','10','10'].map(parseInt);
            var variable5 = 'foo ' + 1 + 1;
            console.log(variable1);
            console.log(variable2);
            console.log(variable3);
            console.log(variable4);
            console.log(variable5);
        })();
    `
	params := config.String(`Params.obfuscate`, "{}")
	engine := NewEngine().SetParams(params)
	result := engine.Encryption(JavaScript)
	log.Printf(`restult:%s`, result)
}
