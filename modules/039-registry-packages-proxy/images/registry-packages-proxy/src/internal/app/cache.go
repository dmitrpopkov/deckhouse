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

package app

import (
	"context"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"registry-packages-proxy/internal/cache"
)

func NewCache(ctx context.Context, config *Config, metrics *cache.Metrics) (*cache.Cache, error) {
	if config.DisableCache {
		return nil, nil
	}

	cache, err := cache.New(config.CacheDirectory, uint64(config.CacheRetentionSize.Value()), metrics)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create cache")
	}

	go func() {
		err := cache.Run(ctx)
		if err != nil {
			log.Errorf("Run cache: %v", err)
		}
	}()

	return cache, nil
}
