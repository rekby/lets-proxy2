package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/gobuffalo/packr"

	"github.com/maxatome/go-testdeep"
)

func TestConfigEmbed(t *testing.T) {
	td := testdeep.NewT(t)
	sourceConfig, err := ioutil.ReadFile("static/default-config.toml")
	td.CmpNoError(err)

	// force remove file - for prevent box read from disk
	err = os.Rename("static/default-config.toml", "static/default-config.toml.tmp")
	td.CmpNoError(err)
	defer os.Rename("static/default-config.toml.tmp", "static/default-config.toml")

	box := packr.NewBox("static")
	boxBytes, err := box.Find("default-config.toml")
	td.CmpNoError(err)
	td.CmpDeeply(string(boxBytes), string(sourceConfig))
}
