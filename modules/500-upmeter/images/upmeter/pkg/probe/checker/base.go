/*
Copyright 2021 Flant JSC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package checker

import (
	"fmt"
	"time"

	"d8.io/upmeter/pkg/check"
	utime "d8.io/upmeter/pkg/time"
)

// Config is basically a checker constructor with verbose arguments
type Config interface {
	Checker() check.Checker
}

// sequenceChecker wraps the sequence of checkers. It returns first-met check error.
// sequenceChecker is stateful, thus should not be reused.
type sequenceChecker struct {
	checkers []check.Checker
	current  int
}

func sequence(first check.Checker, others ...check.Checker) check.Checker {
	checkers := make([]check.Checker, 1+len(others))
	checkers[0] = first
	copy(checkers[1:], others)
	return &sequenceChecker{checkers: checkers}
}

func (c *sequenceChecker) BusyWith() string {
	return c.checkers[c.current].BusyWith()
}

func (c *sequenceChecker) Check() check.Error {
	for i, checker := range c.checkers {
		c.current = i
		err := checker.Check()
		if err != nil {
			return err
		}
	}
	return nil
}

// FailChecker wraps a checker and forces any error it returns to have `fail` status
type FailChecker struct {
	checker check.Checker
}

func failOnError(checker check.Checker) check.Checker {
	return &FailChecker{checker}
}

func (c *FailChecker) BusyWith() string {
	return c.checker.BusyWith()
}

func (c *FailChecker) Check() check.Error {
	err := c.checker.Check()
	if err != nil {
		return check.ErrFail(err.Error())
	}
	return nil
}

// timeoutChecker wraps a checker with timer. If the timer finishes before the wrapped checker,
// the check returns unknown result error.
type timeoutChecker struct {
	checker check.Checker
	timeout time.Duration
}

func withTimeout(checker check.Checker, timeout time.Duration) check.Checker {
	return &timeoutChecker{
		checker: checker,
		timeout: timeout,
	}
}

func (c *timeoutChecker) BusyWith() string {
	return c.checker.BusyWith()
}

func (c *timeoutChecker) Check() check.Error {
	var err check.Error
	utime.DoWithTimer(c.timeout,
		func() {
			err = c.checker.Check()
		},
		func() {
			err = check.ErrUnknown("timed out: %s", c.checker.BusyWith())
		},
	)
	return err
}

// retryChecker launches passed checker `tries` times with the given interval between calls. It is up to the user
// to control timeout in the inner checker
type retryChecker struct {
	checker  check.Checker
	tries    int
	interval time.Duration
}

func withRetryEachSeconds(checker check.Checker, timeout time.Duration) check.Checker {
	interval := time.Second
	return &retryChecker{
		checker:  withTimeout(checker, interval),
		tries:    int(timeout / interval),
		interval: interval,
	}
}

func (c *retryChecker) BusyWith() string {
	return fmt.Sprintf("retrying %s", c.checker.BusyWith())
}

func (c *retryChecker) Check() check.Error {
	var err check.Error

	// FIXME do not add the interval to execution time, not exactly what was expected
	for i := c.tries; i > 0; i-- {
		time.Sleep(c.interval)

		err = c.checker.Check()
		if err == nil {
			// success achieved, stop retrying
			break
		}
	}

	return err
}
