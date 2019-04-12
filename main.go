package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/pelletier/go-toml"
)

func main() {
	globalContext := context.Background()
	_ = globalContext
	flag.Parse()

	if true {
		var config configType
		err := toml.Unmarshal([]byte(``), &config)
		if err != nil {
			panic(err)
		}
		configBytes, _ := toml.Marshal(config)
		fmt.Println(string(configBytes))
	}
}
