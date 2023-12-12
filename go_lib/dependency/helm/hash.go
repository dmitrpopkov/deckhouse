/*
Copyright 2023 Flant JSC

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

package helm

import (
	"fmt"
	"hash"
	"sort"
)

func hashObject(sum map[string]interface{}, hash *hash.Hash) {

	barrier(hash)
	var keys = sortedKeys(sum)
	for _, key := range keys {
		v := sum[key]

		escapeKey(key, hash)
		switch o := v.(type) {
		case map[string]interface{}:
			hashObject(o, hash)
		case []interface{}:
			hashArray(o, hash)
		default:
			escapeValue(o, hash)
		}
	}
	barrier(hash)
}

func hashArray(array []interface{}, hash *hash.Hash) {

	barrier(hash)
	for _, v := range array {

		switch o := v.(type) {
		case map[string]interface{}:
			hashObject(o, hash)
		case []interface{}:
			hashArray(o, hash)
		default:
			escapeValue(o, hash)
		}
	}
	barrier(hash)
}

func escapeKey(key string, writer *hash.Hash) {
	(*writer).Write([]byte(key))
	barrier(writer)
}

func escapeValue(value interface{}, writer *hash.Hash) {
	(*writer).Write([]byte(fmt.Sprintf(`%T`, value)))
	barrier(writer)
	(*writer).Write([]byte(fmt.Sprintf("%v", value)))
	barrier(writer)
}

// an invalide utf8 byte is used as a barrier
var barrierValue = []byte{byte('\255')}

func barrier(writer *hash.Hash) {
	(*writer).Write(barrierValue)
}

func sortedKeys(json map[string]interface{}) []string {
	var keys []string
	for k := range json {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	return keys
}
