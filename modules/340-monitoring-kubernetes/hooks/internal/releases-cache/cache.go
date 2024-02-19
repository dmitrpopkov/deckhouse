/*
Copyright 2024 Flant JSC

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

package releasescache

import (
	"errors"
	"sync"
	"time"
)

type cache struct {
	m             sync.RWMutex
	timemark      int64
	val           []byte
	helm3Releases uint32
	helm2Releases uint32
}

var ch cache

func GetInstance() *cache {
	return &ch
}

func (c *cache) Set(helm3Releases, helm2Releases uint32, val []byte) {
	c.m.Lock()
	defer c.m.Unlock()

	ch.timemark = time.Now().Unix()
	copy(ch.val, val)
	ch.helm3Releases = helm3Releases
	ch.helm2Releases = helm2Releases
}

func (c *cache) Get(ttl time.Duration) (helm3Releases, helm2Releases uint32, val []byte, err error) {
	c.m.RLock()
	defer c.m.RUnlock()

	timeDeep := time.Now().Add(-ttl).Unix()
	if ttl == 0 || c.timemark > timeDeep {
		copy(val, c.val)
		return ch.helm3Releases, ch.helm2Releases, val, nil
	}

	return 0, 0, nil, errors.New("time to live expired")
}
