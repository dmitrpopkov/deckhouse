// Copyright 2023 Flant JSC
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

package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sync/atomic"

	"github.com/spf13/fsync"
	"k8s.io/klog/v2"

	"github.com/flant/docs-builder/pkg/hugo"
)

func newBuildHandler(src, dst string, wasCalled *atomic.Bool) *buildHandler {
	return &buildHandler{
		src:       src,
		dst:       dst,
		wasCalled: wasCalled,
	}
}

type buildHandler struct {
	src       string
	dst       string
	wasCalled *atomic.Bool
}

func (b *buildHandler) ServeHTTP(writer http.ResponseWriter, _ *http.Request) {
	err := b.build()
	if err != nil {
		klog.Error(err)
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
}

func (b *buildHandler) build() error {
	err := b.buildHugo()
	if err != nil {
		return fmt.Errorf("hugo build: %w", err)
	}

	for _, lang := range []string{"ru", "en"} {
		glob := filepath.Join(b.dst, "public", lang, "modules/*")
		err = removeGlob(glob)
		if err != nil {
			return fmt.Errorf("clear %s: %w", b.dst, err)
		}

		oldLocation := filepath.Join(b.src, "public", lang, "modules")
		newLocation := filepath.Join(b.dst, "public", lang, "modules")
		err = fsync.Sync(newLocation, oldLocation)
		if err != nil {
			return fmt.Errorf("move %s to %s: %w", oldLocation, newLocation, err)
		}
	}

	b.wasCalled.Store(true)
	return nil
}

func (b *buildHandler) buildHugo() error {
	flags := hugo.Flags{
		LogLevel: "debug",
		Source:   b.src,
		CfgDir:   filepath.Join(b.src, "config"),
	}

	for {
		err := hugo.Build(flags)
		if err == nil {
			return nil
		}

		if b.removeBrokenFile(err) {
			continue
		}

		return err
	}
}

func (b *buildHandler) removeBrokenFile(err error) bool {
	match := assembleErrorRegexp.FindStringSubmatch(err.Error())
	if match == nil || len(match) != 4 {
		return false
	}

	path := match[1]
	err = os.Remove(path)
	if err != nil {
		klog.Errorf("remove broken source file: %s", err.Error())
		return false
	}

	klog.Warningf("removed broken source file %q", path)
	return true
}

var assembleErrorRegexp = regexp.MustCompile(`error building site: assemble: \x1b\[1;36m"(?P<path>.+):(?P<line>\d+):(?P<column>\d+)"\x1b\[0m:`)

func removeGlob(path string) error {
	contents, err := filepath.Glob(path)
	if err != nil {
		return fmt.Errorf("glob: %w", err)
	}

	for _, item := range contents {
		err = os.RemoveAll(item)
		if err != nil {
			return fmt.Errorf("remove all: %w", err)
		}
	}

	return nil
}
