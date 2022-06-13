package main

import "flag"

var (
	configFileP      = flag.String("config", "config.tom[l]", "Path to config file. Internally expand glob syntax.")
	debugLog         = flag.Bool("debug", false, "Enable debug logging")
	defaultConfigP   = flag.Bool("print-default-config", false, "Write default config to stdout and exit.")
	versionP         = flag.Bool("version", false, "print version and exit")
	testAcmeServerP  = flag.Bool("test-acme-server", false, "Use test acme server, instead address from config")
	manualAcmeServer = flag.String("acme-server", "", "Override acme server")
)
