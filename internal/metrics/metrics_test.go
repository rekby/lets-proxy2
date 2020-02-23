package metrics

import (
	"sort"
	"testing"

	"github.com/maxatome/go-testdeep"

	"github.com/prometheus/client_golang/prometheus"
)

//go:generate minimock -i github.com/prometheus/client_golang/prometheus.Registerer -o ./registerer_mock_test.go

func TestToefCounters(t *testing.T) {
	r := NewRegistererMock(t)
	defer r.MinimockFinish()

	td := testdeep.NewT(t)

	r.MustRegisterMock.Set(func(args ...prometheus.Collector) {
		descriptions := make([]string, len(args))

		for i, arg := range args {
			iDescChan := make(chan *prometheus.Desc, 1)
			arg.Describe(iDescChan)
			descriptions[i] = (<-iDescChan).String()
		}

		sort.Strings(descriptions)

		td.Len(descriptions, 4)
		td.Contains(descriptions[0], "asd")
		td.Contains(descriptions[1], "asd")
		td.Contains(descriptions[2], "asd")
		td.Contains(descriptions[3], "asd")
		td.Contains(descriptions[0], "test")
		td.Contains(descriptions[1], "test_err")
		td.Contains(descriptions[2], "test_inflight")
		td.Contains(descriptions[3], "test_ok")
	})
	ToefCounters(r, "test", "asd")
}
