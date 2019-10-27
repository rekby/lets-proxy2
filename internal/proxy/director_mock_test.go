package proxy

// DO NOT EDIT!
// The code below was generated with http://github.com/gojuno/minimock (dev)

//go:generate minimock -i github.com/rekby/lets-proxy2/internal/proxy.Director -o ./director_mock_test.go

import (
	"sync/atomic"
	"time"

	"net/http"

	"github.com/gojuno/minimock/v3"
)

// DirectorMock implements Director
type DirectorMock struct {
	t minimock.Tester

	funcDirector          func(request *http.Request)
	afterDirectorCounter  uint64
	beforeDirectorCounter uint64
	DirectorMock          mDirectorMockDirector
}

// NewDirectorMock returns a mock for Director
func NewDirectorMock(t minimock.Tester) *DirectorMock {
	m := &DirectorMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}
	m.DirectorMock = mDirectorMockDirector{mock: m}

	return m
}

type mDirectorMockDirector struct {
	mock               *DirectorMock
	defaultExpectation *DirectorMockDirectorExpectation
	expectations       []*DirectorMockDirectorExpectation
}

// DirectorMockDirectorExpectation specifies expectation struct of the Director.Director
type DirectorMockDirectorExpectation struct {
	mock   *DirectorMock
	params *DirectorMockDirectorParams

	Counter uint64
}

// DirectorMockDirectorParams contains parameters of the Director.Director
type DirectorMockDirectorParams struct {
	request *http.Request
}

// Expect sets up expected params for Director.Director
func (m *mDirectorMockDirector) Expect(request *http.Request) *mDirectorMockDirector {
	if m.mock.funcDirector != nil {
		m.mock.t.Fatalf("DirectorMock.Director mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &DirectorMockDirectorExpectation{}
	}

	m.defaultExpectation.params = &DirectorMockDirectorParams{request}
	for _, e := range m.expectations {
		if minimock.Equal(e.params, m.defaultExpectation.params) {
			m.mock.t.Fatalf("Expectation set by When has same params: %#v", *m.defaultExpectation.params)
		}
	}

	return m
}

// Return sets up results that will be returned by Director.Director
func (m *mDirectorMockDirector) Return() *DirectorMock {
	if m.mock.funcDirector != nil {
		m.mock.t.Fatalf("DirectorMock.Director mock is already set by Set")
	}

	if m.defaultExpectation == nil {
		m.defaultExpectation = &DirectorMockDirectorExpectation{mock: m.mock}
	}

	return m.mock
}

//Set uses given function f to mock the Director.Director method
func (m *mDirectorMockDirector) Set(f func(request *http.Request)) *DirectorMock {
	if m.defaultExpectation != nil {
		m.mock.t.Fatalf("Default expectation is already set for the Director.Director method")
	}

	if len(m.expectations) > 0 {
		m.mock.t.Fatalf("Some expectations are already set for the Director.Director method")
	}

	m.mock.funcDirector = f
	return m.mock
}

// Director implements Director
func (m *DirectorMock) Director(request *http.Request) {
	atomic.AddUint64(&m.beforeDirectorCounter, 1)
	defer atomic.AddUint64(&m.afterDirectorCounter, 1)

	for _, e := range m.DirectorMock.expectations {
		if minimock.Equal(*e.params, DirectorMockDirectorParams{request}) {
			atomic.AddUint64(&e.Counter, 1)
			return
		}
	}

	if m.DirectorMock.defaultExpectation != nil {
		atomic.AddUint64(&m.DirectorMock.defaultExpectation.Counter, 1)
		want := m.DirectorMock.defaultExpectation.params
		got := DirectorMockDirectorParams{request}
		if want != nil && !minimock.Equal(*want, got) {
			m.t.Errorf("DirectorMock.Director got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		return

	}
	if m.funcDirector != nil {
		m.funcDirector(request)
		return
	}
	m.t.Fatalf("Unexpected call to DirectorMock.Director. %v", request)

}

// DirectorAfterCounter returns a count of finished DirectorMock.Director invocations
func (m *DirectorMock) DirectorAfterCounter() uint64 {
	return atomic.LoadUint64(&m.afterDirectorCounter)
}

// DirectorBeforeCounter returns a count of DirectorMock.Director invocations
func (m *DirectorMock) DirectorBeforeCounter() uint64 {
	return atomic.LoadUint64(&m.beforeDirectorCounter)
}

// MinimockDirectorDone returns true if the count of the Director invocations corresponds
// the number of defined expectations
func (m *DirectorMock) MinimockDirectorDone() bool {
	for _, e := range m.DirectorMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.DirectorMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterDirectorCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcDirector != nil && atomic.LoadUint64(&m.afterDirectorCounter) < 1 {
		return false
	}
	return true
}

// MinimockDirectorInspect logs each unmet expectation
func (m *DirectorMock) MinimockDirectorInspect() {
	for _, e := range m.DirectorMock.expectations {
		if atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to DirectorMock.Director with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.DirectorMock.defaultExpectation != nil && atomic.LoadUint64(&m.afterDirectorCounter) < 1 {
		m.t.Errorf("Expected call to DirectorMock.Director with params: %#v", *m.DirectorMock.defaultExpectation.params)
	}
	// if func was set then invocations count should be greater than zero
	if m.funcDirector != nil && atomic.LoadUint64(&m.afterDirectorCounter) < 1 {
		m.t.Error("Expected call to DirectorMock.Director")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *DirectorMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockDirectorInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *DirectorMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-time.After(10 * time.Millisecond):
		}
	}
}

func (m *DirectorMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockDirectorDone()
}
