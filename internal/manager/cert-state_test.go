package manager

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"sort"
	"sync"
	"testing"
	"time"

	td "github.com/maxatome/go-testdeep"
	"github.com/rekby/lets-proxy2/internal/th"
)

func TestCertState(t *testing.T) {
	ctx, flush := th.TestContext()
	defer flush()

	s := &certState{}
	td.CmpTrue(t, s.StartIssue(ctx))
	td.CmpFalse(t, s.StartIssue(ctx))

	cert := &tls.Certificate{Leaf: &x509.Certificate{Subject: pkix.Name{CommonName: "asd"}}}

	td.CmpNotPanic(t, func() {
		s.FinishIssue(ctx, cert)
	})

	td.CmpDeeply(t, s.Cert(), cert)

	cert2 := &tls.Certificate{Leaf: &x509.Certificate{Subject: pkix.Name{CommonName: "asdf"}}}

	td.CmpPanic(t, func() {
		s.FinishIssue(th.NoLog(ctx), cert2)
	}, td.NotEmpty())

	td.CmpDeeply(t, s.Cert(), cert2)

}

func TestCertStateManyIssuers(t *testing.T) {
	ctx, flush := th.TestContext()
	defer flush()

	const cnt = 1000
	const pause = 1
	const checkEvery = 1000

	timeoutCtx, _ := context.WithTimeout(ctx, time.Second)

	type lockTimeStruct struct {
		start time.Time
		end   time.Time
	}

	ctxNoLog := th.NoLog(ctx)

	s := certState{}

	lockFunc := func() []lockTimeStruct {
		res := make([]lockTimeStruct, 0, cnt)
		i := 0
		for {
			if i%checkEvery == 0 {
				if timeoutCtx.Err() != nil {
					return res
				}
			}
			i++

			if s.StartIssue(ctxNoLog) {
				item := lockTimeStruct{start: time.Now()}
				time.Sleep(pause)
				item.end = time.Now()
				s.FinishIssue(ctxNoLog, nil)
				res = append(res, item)
				i = 0 // for check exit
			}
		}
	}

	var wg sync.WaitGroup
	wg.Add(cnt)

	lockTimesChan := make(chan []lockTimeStruct, cnt)

	for i := 0; i < cnt; i++ {
		go func() {
			lockTimesChan <- lockFunc()
			wg.Done()
		}()
	}
	wg.Wait()
	close(lockTimesChan)

	var lockTimesSlice []lockTimeStruct
	for i := 0; i < cnt; i++ {
		items := <-lockTimesChan
		lockTimesSlice = append(lockTimesSlice, items...)
	}

	sort.Slice(lockTimesSlice, func(i, j int) bool {
		left := lockTimesSlice[i]
		right := lockTimesSlice[j]

		if left.start.Before(right.start) {
			return true
		}
		if left.start.Equal(right.start) {
			return left.end.Before(right.end)
		}
		return false
	})

	// check
	for i := 0; i < len(lockTimesSlice)-1; i++ {
		left := lockTimesSlice[i]
		right := lockTimesSlice[i+1]

		if left.start == right.start {
			t.Error(left, right)
		}
		if left.end == right.end {
			t.Error()
		}
		if left.start.Before(right.start) && left.end.After(right.start) {
			t.Error()
		}
		if left.start.Before(right.end) && left.end.After(right.end) {
			t.Error()
		}
	}
	t.Logf("Succesful locks: %d", len(lockTimesSlice))
}

func TestCertState_WaitFinishIssue(t *testing.T) {
	ctx, flush := th.TestContext()
	defer flush()

	s := certState{}

	const timeout = time.Millisecond * 10

	ctxTimeout, _ := context.WithTimeout(ctx, timeout)
	td.CmpNoError(t, s.WaitFinishIssue(ctxTimeout))

	s.StartIssue(context.Background())
	ctxTimeout, _ = context.WithTimeout(ctx, timeout)
	td.CmpError(t, s.WaitFinishIssue(ctxTimeout))

	go func() {
		time.Sleep(timeout / 2)
		s.FinishIssue(ctx, nil)
	}()
	ctxTimeout, _ = context.WithTimeout(context.Background(), timeout)
	td.CmpNoError(t, s.WaitFinishIssue(ctxTimeout))
}
