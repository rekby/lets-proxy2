//go:build go1.17
// +build go1.17

package cert_manager

import (
	"crypto/tls"
	"testing"

	"github.com/rekby/lets-proxy2/internal/th"
)

func TestInsecureChipers(t *testing.T) {
	e, _, flush := th.NewEnv(t)
	defer flush()

	localInsecureMap := make(map[uint16]struct{})
	for _, suite := range tls.InsecureCipherSuites() {
		localInsecureMap[suite.ID] = struct{}{}
	}
	e.Cmp(localInsecureMap, insecureChipers)
}
