// Copyright 2022 Flant JSC
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

package fs

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func RandomFileName() string {
	// we silent gosec linter here
	// because we do not need security random number
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s) //nolint:gosec

	rndSuf := strconv.FormatUint(r.Uint64(), 10)
	fileName := fmt.Sprintf("dhctl-tst-touch-%s", rndSuf)

	return filepath.Join(os.TempDir(), fileName)
}
