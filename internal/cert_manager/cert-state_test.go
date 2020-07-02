//nolint:golint
package cert_manager

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/maxatome/go-testdeep"
	"github.com/rekby/lets-proxy2/internal/th"
)

func TestCertState(t *testing.T) {
	ctx, flush := th.TestContext(t)
	defer flush()

	s := &certState{}
	testdeep.CmpTrue(t, s.StartIssue(ctx))
	testdeep.CmpFalse(t, s.StartIssue(ctx))

	cert := &tls.Certificate{Leaf: &x509.Certificate{Subject: pkix.Name{CommonName: "asd"}}}

	s.FinishIssue(ctx, cert, nil)

	rCert, rErr := s.Cert()
	testdeep.CmpDeeply(t, rCert, cert)
	testdeep.CmpNil(t, rErr)

	s = &certState{}
	err1 := errors.New("1")
	testdeep.CmpTrue(t, s.StartIssue(ctx))
	s.FinishIssue(ctx, nil, err1)
	rCert, rErr = s.Cert()
	testdeep.CmpNil(t, rCert)
	testdeep.CmpDeeply(t, rErr, err1)
}

func TestCertStateManyIssuers(t *testing.T) {
	ctx, flush := th.TestContext(t)
	defer flush()

	const cnt = 1000
	const pause = 1
	const checkEvery = 1000

	//nolint:govet
	timeoutCtx, _ := context.WithTimeout(ctx, time.Second)

	type lockTimeStruct struct {
		start time.Time
		end   time.Time
	}

	ctxNoLog := th.NoLog(ctx)

	s := certState{}
	err1 := errors.New("test noerror")

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
				s.FinishIssue(ctxNoLog, nil, err1)

				res = append(res, item)

				i = 0 // for check exit
			}
		}
	}

	var wg sync.WaitGroup
	wg.Add(cnt) //nolint:wsl

	lockTimesChan := make(chan []lockTimeStruct, cnt)

	for i := 0; i < cnt; i++ {
		go func() {
			defer wg.Done()

			lockTimesChan <- lockFunc()
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
	t.Logf("Successful locks: %d", len(lockTimesSlice))
}

func TestCertState_WaitFinishIssue(t *testing.T) {
	ctx, flush := th.TestContext(t)
	defer flush()

	s := certState{}

	const timeout = time.Millisecond * 10

	//nolint:govet
	ctxTimeout, _ := context.WithTimeout(ctx, timeout)
	rCert, rErr := s.WaitFinishIssue(ctxTimeout)
	testdeep.CmpNil(t, rCert)
	testdeep.CmpNil(t, rErr)

	s.StartIssue(ctx)
	//nolint:govet
	ctxTimeout, _ = context.WithTimeout(ctx, timeout)
	rCert, rErr = s.WaitFinishIssue(ctxTimeout)
	testdeep.CmpNil(t, rCert)
	testdeep.CmpError(t, rErr)

	cert1 := &tls.Certificate{Leaf: &x509.Certificate{Subject: pkix.Name{CommonName: "asdasd"}}}
	go func() {
		time.Sleep(timeout / 2)
		s.FinishIssue(ctx, cert1, nil)
	}()
	//nolint:govet
	ctxTimeout, _ = context.WithTimeout(ctx, timeout)
	rCert, rErr = s.WaitFinishIssue(ctxTimeout)
	testdeep.CmpNoError(t, rErr)
	testdeep.CmpDeeply(t, rCert, cert1)

	s.StartIssue(ctx)
	err2 := errors.New("2")
	go func() {
		time.Sleep(timeout / 2)
		s.FinishIssue(ctx, nil, err2)
	}()
	//nolint:govet
	ctxTimeout, _ = context.WithTimeout(ctx, timeout)
	rCert, rErr = s.WaitFinishIssue(ctxTimeout)
	testdeep.CmpNil(t, rCert)
	testdeep.CmpDeeply(t, rErr, err2)
}

func TestCertState_FinishIssuePanic(t *testing.T) {
	ctx, flush := th.TestContext(t)
	defer flush()

	ctx = th.NoLog(ctx)
	s := certState{}

	cert1 := &tls.Certificate{Leaf: &x509.Certificate{Subject: pkix.Name{CommonName: "asdf"}}}
	err1 := errors.New("2")

	testdeep.CmpPanic(t, func() {
		s.FinishIssue(th.NoLog(ctx), cert1, nil)
	}, testdeep.NotEmpty())

	rCert, rErr := s.Cert()
	testdeep.CmpDeeply(t, rCert, cert1)
	testdeep.CmpNil(t, rErr)

	s = certState{}
	s.StartIssue(ctx)
	testdeep.CmpPanic(t, func() {
		s.FinishIssue(th.NoLog(ctx), nil, nil)
	}, testdeep.NotEmpty())

	s = certState{}
	s.StartIssue(ctx)
	testdeep.CmpPanic(t, func() {
		s.FinishIssue(th.NoLog(ctx), cert1, err1)
	}, testdeep.NotEmpty())
}

func TestCertState_CertSet(t *testing.T) {
	ctx, flush := th.TestContext(t)
	defer flush()

	td := testdeep.NewT(t)
	s := certState{}
	cert := &tls.Certificate{
		OCSPStaple: []byte{1, 2, 3},
	}
	s.CertSet(ctx, true, cert)
	td.CmpDeeply(s.cert, cert)
	td.True(s.useAsIs)

	s.CertSet(ctx, false, nil)
	td.Nil(s.cert)
	td.False(s.useAsIs)
}

func TestCertState_GetLocked(t *testing.T) {
	td := testdeep.NewT(t)

	s := certState{}
	td.False(s.GetUseAsIs())

	s.useAsIs = true
	td.True(s.GetUseAsIs())
}
