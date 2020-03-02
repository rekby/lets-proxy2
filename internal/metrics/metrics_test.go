package metrics

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"

	io_prometheus_client "github.com/prometheus/client_model/go"

	"github.com/maxatome/go-testdeep"

	"github.com/prometheus/client_golang/prometheus"
)

func TestNew(t *testing.T) {
	td := testdeep.NewT(t)
	logger := zap.NewNop()

	gatherer := NewGathererMock(t)
	defer gatherer.MinimockFinish()
	gatherer.GatherMock.Set(func() (metrics []*io_prometheus_client.MetricFamily, err error) {
		return nil, nil
	})

	metrics := New(logger, gatherer)
	td.CmpDeeply(metrics.logger, logger)
	td.NotNil(metrics.metricsHandler)

	req, _ := http.NewRequest(http.MethodGet, "", nil)
	metrics.ServeHTTP(httptest.NewRecorder(), req)
}

func TestToefCounters(t *testing.T) {
	td := testdeep.NewT(t)
	getDesc := func(c prometheus.Collector) string {
		iDescChan := make(chan *prometheus.Desc, 1)
		c.Describe(iDescChan)
		res := (<-iDescChan).String()
		return res
	}
	getCount := func(c prometheus.Collector) int {
		iMetricChan := make(chan prometheus.Metric, 1)
		c.Collect(iMetricChan)
		metric := <-iMetricChan
		metProto := io_prometheus_client.Metric{}
		err := metric.Write(&metProto)
		td.CmpNoError(err)
		if metProto.Counter != nil {
			return int(*metProto.Counter.Value)
		}
		return int(*metProto.Gauge.Value)
	}

	var cntTotal, cntOk, cntErr, cntInFly prometheus.Collector

	r := NewRegistererMock(t)
	defer r.MinimockFinish()

	r.MustRegisterMock.Set(func(args ...prometheus.Collector) {
		cntTotal, cntOk, cntErr, cntInFly = args[0], args[1], args[2], args[3]

		td.Len(args, 4)
		for i, arg := range args {
			td.Contains(getDesc(arg), "asd", i)
		}
		td.Contains(getDesc(cntTotal), "test")
		td.Contains(getDesc(cntOk), "test_ok")
		td.Contains(getDesc(cntErr), "test_err")
		td.Contains(getDesc(cntInFly), "test_inflight")
	})

	start, finish := ToefCounters(r, "test", "asd")

	td.Cmp(getCount(cntTotal), 0)
	td.Cmp(getCount(cntOk), 0)
	td.Cmp(getCount(cntErr), 0)
	td.Cmp(getCount(cntInFly), 0)

	start()
	td.Cmp(getCount(cntTotal), 1)
	td.Cmp(getCount(cntOk), 0)
	td.Cmp(getCount(cntErr), 0)
	td.Cmp(getCount(cntInFly), 1)

	finish(nil)
	td.Cmp(getCount(cntTotal), 1)
	td.Cmp(getCount(cntOk), 1)
	td.Cmp(getCount(cntErr), 0)
	td.Cmp(getCount(cntInFly), 0)

	start()
	td.Cmp(getCount(cntTotal), 2)
	td.Cmp(getCount(cntOk), 1)
	td.Cmp(getCount(cntErr), 0)
	td.Cmp(getCount(cntInFly), 1)

	finish(errors.New("test"))
	td.Cmp(getCount(cntTotal), 2)
	td.Cmp(getCount(cntOk), 1)
	td.Cmp(getCount(cntErr), 1)
	td.Cmp(getCount(cntInFly), 0)

	start, finish = ToefCounters(nil, "qwe", "ddd")
	start()
	td.Cmp(getCount(cntTotal), 2)
	td.Cmp(getCount(cntOk), 1)
	td.Cmp(getCount(cntErr), 1)
	td.Cmp(getCount(cntInFly), 0)

	finish(nil)
	td.Cmp(getCount(cntTotal), 2)
	td.Cmp(getCount(cntOk), 1)
	td.Cmp(getCount(cntErr), 1)
	td.Cmp(getCount(cntInFly), 0)

	finish(errors.New("test"))
	td.Cmp(getCount(cntTotal), 2)
	td.Cmp(getCount(cntOk), 1)
	td.Cmp(getCount(cntErr), 1)
	td.Cmp(getCount(cntInFly), 0)
}

func TestErrorLoggger_Println(t *testing.T) {
	loggerMock := NewLoggerErrorMock(t)
	defer loggerMock.MinimockFinish()

	loggerMock.ErrorMock.Expect("a", "s", 3).Return()
	logger := errorLoggger{loggerMock}
	logger.Println("a", "s", 3)
}
