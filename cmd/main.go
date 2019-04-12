package main

import (
	"context"
	"flag"
	"fmt"
	"os"
)

func main() {
	globalContext := context.Background()
	_ = globalContext
	flag.Parse()
	if *defaultConfigP {
		fmt.Println(string(defaultConfig(globalContext)))
		os.Exit(0)
	}

	getConfig(globalContext)

}
