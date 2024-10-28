// Copyright 2024 Flant JSC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logger_test

import (
	"bytes"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/deckhouse/deckhouse/dhctl/pkg/server/pkg/logger"
	"github.com/flant/shell-operator/pkg/unilogger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogWriter(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input [][]byte
		lines [][]string
		f     func([]string) []string
	}{
		"one line of text": {
			input: [][]byte{
				[]byte("one line of text\n"),
			},
			lines: [][]string{
				{"one line of text"},
			},
			f: func(lines []string) []string { return lines },
		},
		"one line of text, multiple writes": {
			input: [][]byte{
				[]byte("one line o"), []byte("f text"), []byte("\n"),
			},
			lines: [][]string{
				{"one line of text"},
			},
			f: func(lines []string) []string { return lines },
		},
		"multiple lines of text, multiple writes": {
			input: [][]byte{
				[]byte("first line o"), []byte("f text"), []byte("\n"),
				[]byte("second line of text\nthird line of text\n"),
			},
			lines: [][]string{
				{"first line of text"},
				{"second line of text", "third line of text"},
			},
			f: func(lines []string) []string { return lines },
		},
		"mutate data": {
			input: [][]byte{
				[]byte("first line of text\n"),
				[]byte("second line of text\n"),
				[]byte("third line of text\n"),
			},
			lines: [][]string{
				{"FIRST LINE OF TEXT"},
				{"SECOND LINE OF TEXT"},
				{"THIRD LINE OF TEXT"},
			},
			f: func(lines []string) []string {
				result := make([]string, 0, len(lines))
				for _, line := range lines {
					result = append(result, strings.ToUpper(line))
				}
				return result
			},
		},
	}

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			sendCh := make(chan []string)
			defer close(sendCh)

			var logBuff bytes.Buffer

			logopts := unilogger.Options{
				Output: &logBuff,
				TimeFunc: func(_ time.Time) time.Time {
					parsedTime, err := time.Parse(time.DateTime, "2006-01-02 15:04:05")
					if err != nil {
						assert.NoError(t, err)
					}

					return parsedTime
				},
			}
			log := unilogger.NewLogger(logopts).With(slog.String("key", "value"))

			w := logger.NewLogWriter(log, sendCh, tt.f)

			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				defer wg.Done()
				var writes int
				for writes < len(tt.lines) {
					lines := <-sendCh
					assert.Equal(t, tt.lines[writes], lines)
					writes++
				}
			}()

			for _, input := range tt.input {
				n, err := w.Write(input)
				require.NoError(t, err)
				assert.EqualValues(t, len(input), n)
			}

			wg.Wait()

			expectedLogLines := make([]string, 0)
			for _, lines := range tt.lines {
				for _, line := range lines {
					lohLine := fmt.Sprintf(`{"level":"info","msg":"%s","key":"value","time":"2006-01-02T15:04:05Z"}`, strings.ToLower(line))
					expectedLogLines = append(expectedLogLines, lohLine)
				}
			}

			logLines := strings.Split(strings.TrimSpace(logBuff.String()), "\n")
			assert.Equal(t, expectedLogLines, logLines)
		})
	}
}
