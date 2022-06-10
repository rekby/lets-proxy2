package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rekby/lets-proxy2/internal/th"

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

	toUnixString := func(source []byte) string {
		s := string(source)
		s = strings.Replace(s, "\r\n", "\n", -1)
		return s
	}

	td.CmpDeeply(toUnixString(defaultConfigContent), toUnixString(sourceConfig))
}

func TestReadConfig(t *testing.T) {
	e, ctx, cancel := th.NewEnv(t)
	defer cancel()

	td := testdeep.NewT(t)
	tmpDir := th.TmpDir(e)

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
	e, ctx, cancel := th.NewEnv(t)
	defer cancel()

	e.NotNil(getConfig(ctx))
}
