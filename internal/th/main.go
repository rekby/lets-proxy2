package th

import (
	"testing"

	"github.com/rekby/fixenv"
)

func InitMain(m *testing.M) {
	fixenv.CreateMainTestEnv(nil)
}
