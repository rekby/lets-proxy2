package th

import (
	"io/ioutil"
	"os"

	"github.com/rekby/fixenv"
)

func TmpDir(e fixenv.Env) string {
	var dirPath string
	return e.Cache(nil, &fixenv.FixtureOptions{CleanupFunc: func() {
		_ = os.RemoveAll(dirPath)
	}}, func() (res interface{}, err error) {
		dirPath, err = ioutil.TempDir("", "lets-proxy2-test-")
		return dirPath, err
	}).(string)
}
