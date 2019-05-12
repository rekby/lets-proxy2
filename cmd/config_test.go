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
	sourceConfig, err := ioutil.ReadFile("cmd/static/default-config.toml")
	td.CmpNoError(err)

	// force remove file - for prevent box read from disk
	err = os.Rename("cmd/static/default-config.toml", "cmd/static/default-config.toml.tmp")
	td.CmpNoError(err)
	defer os.Rename("cmd/static/default-config.toml.tmp", "cmd/static/default-config.toml")

	box := packr.NewBox("static")
	boxBytes, err := box.Find("default-config.toml")
	td.CmpNoError(err)
	td.CmpDeeply(string(boxBytes), string(sourceConfig))
}
