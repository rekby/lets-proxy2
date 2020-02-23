package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/rekby/lets-proxy2/internal/th"

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

func TestReadConfig(t *testing.T) {
	ctx, cancel := th.TestContext(t)
	defer cancel()

	td := testdeep.NewT(t)
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		td.Fatal(err)
	}
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	_ = ioutil.WriteFile(filepath.Join(tmpDir, "config.toml"), []byte(`
[General]
IssueTimeout = 1
StorageDir = "storage1"
IncludeConfigs = ["configs/*.toml"]
`), 0600)
	_ = os.MkdirAll(filepath.Join(tmpDir, "configs"), 0700)
	_ = ioutil.WriteFile(filepath.Join(tmpDir, "configs/config2.toml"), []byte(`
[General]
StorageDir = "storage2"
`), 0600)

	var config configType

	mergeConfigBytes(ctx, &config, defaultConfig(ctx), "")
	mergeConfigByFilepath(ctx, &config, filepath.Join(tmpDir, "config.toml"))
	td.CmpDeeply(config.General.IssueTimeout, 1)
	td.CmpDeeply(config.General.StorageDir, "storage2")
}

func TestGetConfig(t *testing.T) {
	ctx, cancel := th.TestContext(t)
	defer cancel()

	testdeep.CmpNotNil(t, getConfig(ctx))
}
