package vsphere

// Code generated by http://github.com/gojuno/minimock (3.0.10). DO NOT EDIT.

//go:generate minimock -i github.com/deckhouse/deckhouse/go_lib/dependency/vsphere.Client -o ./vsphere_mock.go -n ClientMock

import (
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
)

// ClientMock implements Client
type ClientMock struct {
	t minimock.Tester

	funcGetZonesDatastores          func() (op1 *Output, err error)
	inspectFuncGetZonesDatastores   func()
	afterGetZonesDatastoresCounter  uint64
	beforeGetZonesDatastoresCounter uint64
	GetZonesDatastoresMock          mClientMockGetZonesDatastores
}

// NewClientMock returns a mock for Client
func NewClientMock(t minimock.Tester) *ClientMock {
	m := &ClientMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetZonesDatastoresMock = mClientMockGetZonesDatastores{mock: m}

	return m
}

type mClientMockGetZonesDatastores struct {
	mock               *ClientMock
	defaultExpectation *ClientMockGetZonesDatastoresExpectation
	expectations       []*ClientMockGetZonesDatastoresExpectation
}

// ClientMockGetZonesDatastoresExpectation specifies expectation struct of the Client.GetZonesDatastores
type ClientMockGetZonesDatastoresExpectation struct {
	mock *ClientMock

	results *ClientMockGetZonesDatastoresResults
	Counter uint64
}

// ClientMockGetZonesDatastoresResults contains results of the Client.GetZonesDatastores
type ClientMockGetZonesDatastoresResults struct {
	op1 *Output
	err error
}

// Expect sets up expected params for Client.GetZonesDatastores
func (mmGetZonesDatastores *mClientMockGetZonesDatastores) Expect() *mClientMockGetZonesDatastores {
	if mmGetZonesDatastores.mock.funcGetZonesDatastores != nil {
		mmGetZonesDatastores.mock.t.Fatalf("ClientMock.GetZonesDatastores mock is already set by Set")
	}

	if mmGetZonesDatastores.defaultExpectation == nil {
		mmGetZonesDatastores.defaultExpectation = &ClientMockGetZonesDatastoresExpectation{}
	}

	return mmGetZonesDatastores
}

// Inspect accepts an inspector function that has same arguments as the Client.GetZonesDatastores
func (mmGetZonesDatastores *mClientMockGetZonesDatastores) Inspect(f func()) *mClientMockGetZonesDatastores {
	if mmGetZonesDatastores.mock.inspectFuncGetZonesDatastores != nil {
		mmGetZonesDatastores.mock.t.Fatalf("Inspect function is already set for ClientMock.GetZonesDatastores")
	}

	mmGetZonesDatastores.mock.inspectFuncGetZonesDatastores = f

	return mmGetZonesDatastores
}

// Return sets up results that will be returned by Client.GetZonesDatastores
func (mmGetZonesDatastores *mClientMockGetZonesDatastores) Return(op1 *Output, err error) *ClientMock {
	if mmGetZonesDatastores.mock.funcGetZonesDatastores != nil {
		mmGetZonesDatastores.mock.t.Fatalf("ClientMock.GetZonesDatastores mock is already set by Set")
	}

	if mmGetZonesDatastores.defaultExpectation == nil {
		mmGetZonesDatastores.defaultExpectation = &ClientMockGetZonesDatastoresExpectation{mock: mmGetZonesDatastores.mock}
	}
	mmGetZonesDatastores.defaultExpectation.results = &ClientMockGetZonesDatastoresResults{op1, err}
	return mmGetZonesDatastores.mock
}

// Set uses given function f to mock the Client.GetZonesDatastores method
func (mmGetZonesDatastores *mClientMockGetZonesDatastores) Set(f func() (op1 *Output, err error)) *ClientMock {
	if mmGetZonesDatastores.defaultExpectation != nil {
		mmGetZonesDatastores.mock.t.Fatalf("Default expectation is already set for the Client.GetZonesDatastores method")
	}

	if len(mmGetZonesDatastores.expectations) > 0 {
		mmGetZonesDatastores.mock.t.Fatalf("Some expectations are already set for the Client.GetZonesDatastores method")
	}

	mmGetZonesDatastores.mock.funcGetZonesDatastores = f
	return mmGetZonesDatastores.mock
}

// GetZonesDatastores implements Client
func (mmGetZonesDatastores *ClientMock) GetZonesDatastores() (op1 *Output, err error) {
	mm_atomic.AddUint64(&mmGetZonesDatastores.beforeGetZonesDatastoresCounter, 1)
	defer mm_atomic.AddUint64(&mmGetZonesDatastores.afterGetZonesDatastoresCounter, 1)

	if mmGetZonesDatastores.inspectFuncGetZonesDatastores != nil {
		mmGetZonesDatastores.inspectFuncGetZonesDatastores()
	}

	if mmGetZonesDatastores.GetZonesDatastoresMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmGetZonesDatastores.GetZonesDatastoresMock.defaultExpectation.Counter, 1)

		mm_results := mmGetZonesDatastores.GetZonesDatastoresMock.defaultExpectation.results
		if mm_results == nil {
			mmGetZonesDatastores.t.Fatal("No results are set for the ClientMock.GetZonesDatastores")
		}
		return (*mm_results).op1, (*mm_results).err
	}
	if mmGetZonesDatastores.funcGetZonesDatastores != nil {
		return mmGetZonesDatastores.funcGetZonesDatastores()
	}
	mmGetZonesDatastores.t.Fatalf("Unexpected call to ClientMock.GetZonesDatastores.")
	return
}

// GetZonesDatastoresAfterCounter returns a count of finished ClientMock.GetZonesDatastores invocations
func (mmGetZonesDatastores *ClientMock) GetZonesDatastoresAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetZonesDatastores.afterGetZonesDatastoresCounter)
}

// GetZonesDatastoresBeforeCounter returns a count of ClientMock.GetZonesDatastores invocations
func (mmGetZonesDatastores *ClientMock) GetZonesDatastoresBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetZonesDatastores.beforeGetZonesDatastoresCounter)
}

// MinimockGetZonesDatastoresDone returns true if the count of the GetZonesDatastores invocations corresponds
// the number of defined expectations
func (m *ClientMock) MinimockGetZonesDatastoresDone() bool {
	for _, e := range m.GetZonesDatastoresMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetZonesDatastoresMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterGetZonesDatastoresCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetZonesDatastores != nil && mm_atomic.LoadUint64(&m.afterGetZonesDatastoresCounter) < 1 {
		return false
	}
	return true
}

// MinimockGetZonesDatastoresInspect logs each unmet expectation
func (m *ClientMock) MinimockGetZonesDatastoresInspect() {
	for _, e := range m.GetZonesDatastoresMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Error("Expected call to ClientMock.GetZonesDatastores")
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetZonesDatastoresMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterGetZonesDatastoresCounter) < 1 {
		m.t.Error("Expected call to ClientMock.GetZonesDatastores")
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetZonesDatastores != nil && mm_atomic.LoadUint64(&m.afterGetZonesDatastoresCounter) < 1 {
		m.t.Error("Expected call to ClientMock.GetZonesDatastores")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *ClientMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockGetZonesDatastoresInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *ClientMock) MinimockWait(timeout mm_time.Duration) {
	timeoutCh := mm_time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-mm_time.After(10 * mm_time.Millisecond):
		}
	}
}

func (m *ClientMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockGetZonesDatastoresDone()
}
