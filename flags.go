package main

import "flag"

var (
	configFileP    = flag.String("config", "", "Path to config file")
	defaultConfigP = flag.Bool("write-default-config", false, "Write default config to stdout and exit.")
)
