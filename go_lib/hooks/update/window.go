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

package update

import (
	"encoding/json"
	"strings"
	"time"
)

const (
	hh_mm = "15:04" // nolint: revive
)

// Windows update windows
type Windows []Window

// Window single window
type Window struct {
	From string   `json:"from"`
	To   string   `json:"to"`
	Days []string `json:"days"`
}

// FromJSON returns update Windows from json
func FromJSON(data []byte) (Windows, error) {
	var w Windows

	err := json.Unmarshal(data, &w)

	return w, err
}

// IsAllowed returns if specified time get into windows
func (ws Windows) IsAllowed(t time.Time) bool {
	if len(ws) == 0 {
		return true
	}

	for _, window := range ws {
		if window.IsAllowed(t) {
			return true
		}
	}

	return false
}

// IsAllowed check if specified window is allowed at the moment or not
func (uw Window) IsAllowed(now time.Time) bool {
	now = now.UTC()
	// input is validated through the openapi spec
	// we must have only a valid time here
	fromInput, _ := time.Parse(hh_mm, uw.From)
	toInput, _ := time.Parse(hh_mm, uw.To)

	fromTime := time.Date(now.Year(), now.Month(), now.Day(), fromInput.Hour(), fromInput.Minute(), 0, 0, time.UTC)
	toTime := time.Date(now.Year(), now.Month(), now.Day(), toInput.Hour(), toInput.Minute(), 0, 0, time.UTC)

	updateToday := uw.isTodayAllowed(now, uw.Days)

	if !updateToday {
		return false
	}

	if now.After(fromTime) && now.Before(toTime) {
		return true
	}

	return false
}

// NextAllowedTime calculates next update window with respect on minimalTime
// if minimal time is out of window - this function checks next days to find the nearest one
func (ws Windows) NextAllowedTime(min time.Time) time.Time {
	min = min.UTC()

	if len(ws) == 0 {
		return min
	}

	var minTime time.Time

	for _, window := range ws {
		var windowMinTime time.Time

		fromInput, _ := time.Parse(hh_mm, window.From)
		toInput, _ := time.Parse(hh_mm, window.To)

		fromTime := time.Date(min.Year(), min.Month(), min.Day(), fromInput.Hour(), fromInput.Minute(), 0, 0, time.UTC)
		toTime := time.Date(min.Year(), min.Month(), min.Day(), toInput.Hour(), toInput.Minute(), 0, 0, time.UTC)

		if window.isTodayAllowed(min, window.Days) {
			if min.After(fromTime) && min.Before(toTime) {
				windowMinTime = min
				if minTime.IsZero() || windowMinTime.Before(minTime) {
					minTime = windowMinTime
				}
				continue
			}
		}

		// if not today
		nextDay := min.AddDate(0, 0, 1)
		for {
			if window.isTodayAllowed(nextDay, window.Days) {
				fromTime = time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(), fromInput.Hour(), fromInput.Minute(), 0, 0, time.UTC)
				toTime = time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(), toInput.Hour(), toInput.Minute(), 0, 0, time.UTC)

				if min.Before(toTime) {
					if min.After(fromTime) {
						windowMinTime = min
						break
					} else {
						windowMinTime = fromTime
						break
					}

				}
			}
			nextDay = nextDay.AddDate(0, 0, 1)
		}

		if minTime.IsZero() || windowMinTime.Before(minTime) {
			minTime = windowMinTime
		}
	}

	return minTime
}

func (uw Window) isDayEqual(today time.Time, dayString string) bool {
	var day time.Weekday

	switch strings.ToLower(dayString) {
	case "mon":
		day = time.Monday

	case "tue":
		day = time.Tuesday

	case "wed":
		day = time.Wednesday

	case "thu":
		day = time.Thursday

	case "fri":
		day = time.Friday

	case "sat":
		day = time.Saturday

	case "sun":
		day = time.Sunday
	}

	return today.Weekday() == day
}

func (uw Window) isTodayAllowed(now time.Time, days []string) bool {
	if len(days) == 0 {
		return true
	}

	for _, day := range days {
		if uw.isDayEqual(now, day) {
			return true
		}
	}

	return false
}

// below is generated code to use update windows in CRDs with DeepCopy funcs

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (uw *Window) DeepCopyInto(out *Window) {
	*out = *uw
	if uw.Days != nil {
		in, out := &uw.Days, &out.Days
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new UpdateWindow.
func (uw *Window) DeepCopy() *Window {
	if uw == nil {
		return nil
	}
	out := new(Window)
	uw.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (ws Windows) DeepCopyInto(out *Windows) {
	{
		in := &ws
		*out = make(Windows, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
		return
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new UpdateWindows.
func (ws Windows) DeepCopy() Windows {
	if ws == nil {
		return nil
	}
	out := new(Windows)
	ws.DeepCopyInto(out)
	return *out
}
