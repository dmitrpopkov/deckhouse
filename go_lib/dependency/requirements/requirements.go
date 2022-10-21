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

package requirements

import (
	"fmt"

	"github.com/pkg/errors"
)

var (
	defaultRegistry  requirementsResolver
	memoryStorage    *MemoryValuesStore
	ErrNotRegistered = errors.New("Not registered")
)

func init() {
	defaultRegistry = newRegistry()
	memoryStorage = newMemoryValuesStore()
}

// RegisterCheck add CheckFunc for some component
func RegisterCheck(key string, f CheckFunc) {
	defaultRegistry.RegisterCheck(key, f)
}

// RegisterDisruption add DisruptionFunc for some component
func RegisterDisruption(key string, f DisruptionFunc) {
	defaultRegistry.RegisterDisruption(key, f)
}

// CheckRequirement run check function for `key` requirement. Returns true if check is passed, false otherwise
func CheckRequirement(key, value string) (bool, error) {
	if defaultRegistry == nil {
		return true, nil
	}

	f, err := defaultRegistry.GetCheckByKey(key)
	if err != nil {
		return false, err
	}

	return f(value, memoryStorage)
}

// HasDisruption run check function for `key` disruption. Returns true if disruption condition is met, false otherwise. Returns reason for true response.
func HasDisruption(key string) (bool, string) {
	if defaultRegistry == nil {
		return false, ""
	}

	f, err := defaultRegistry.GetDisruptionByKey(key)
	if err != nil {
		return false, ""
	}

	return f(memoryStorage)
}

// SaveValue could be used in the modules, to store their internal values for updater
// One module does not have access to the other's module values, so we can do it through this interface
func SaveValue(key string, value interface{}) {
	memoryStorage.Set(key, value)
}

// RemoveValue remove previously stored value
func RemoveValue(key string) {
	memoryStorage.Remove(key)
}

// GetValue returns saved value. !Attention: Please don't use it in hooks, only for tests
func GetValue(key string) (interface{}, bool) {
	return memoryStorage.Get(key)
}

// CheckFunc check come precondition, comparing desired value (requirementValue) with current value (getter)
type CheckFunc func(requirementValue string, getter ValueGetter) (bool, error)

// DisruptionFunc implements inner logic to warn users about potentially dangerous changes
type DisruptionFunc func(getter ValueGetter) (bool, string)

type ValueGetter interface {
	Get(path string) (interface{}, bool)
}

type requirementsResolver interface {
	RegisterCheck(key string, f CheckFunc)
	GetCheckByKey(key string) (CheckFunc, error)

	RegisterDisruption(key string, f DisruptionFunc)
	GetDisruptionByKey(key string) (DisruptionFunc, error)
}

type requirementsRegistry struct {
	checkers    map[string]CheckFunc
	disruptions map[string]DisruptionFunc
}

func newRegistry() *requirementsRegistry {
	return &requirementsRegistry{
		checkers:    make(map[string]CheckFunc),
		disruptions: make(map[string]DisruptionFunc),
	}
}

func (r *requirementsRegistry) RegisterCheck(key string, f CheckFunc) {
	r.checkers[key] = f
}

func (r *requirementsRegistry) RegisterDisruption(key string, f DisruptionFunc) {
	r.disruptions[key] = f
}

func (r *requirementsRegistry) GetCheckByKey(key string) (CheckFunc, error) {
	f, ok := r.checkers[key]
	if !ok {
		return nil, errors.Wrap(ErrNotRegistered, fmt.Sprintf("requirement with a key: %s", key))
	}

	return f, nil
}

func (r *requirementsRegistry) GetDisruptionByKey(key string) (DisruptionFunc, error) {
	f, ok := r.disruptions[key]
	if !ok {
		return nil, errors.Wrap(ErrNotRegistered, fmt.Sprintf("disruption with a key: %s", key))
	}

	return f, nil
}
