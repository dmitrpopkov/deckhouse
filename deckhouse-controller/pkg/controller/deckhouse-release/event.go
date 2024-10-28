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

package deckhouse_release

import (
	"runtime/debug"

	"github.com/flant/shell-operator/pkg/unilogger"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/deckhouse/deckhouse/deckhouse-controller/pkg/apis/deckhouse.io/v1alpha1"
)

func newEventFilter() predicate.Predicate {
	return predicate.And(
		predicate.Or(predicate.GenerationChangedPredicate{}, predicate.AnnotationChangedPredicate{}),
		releasePhasePredicate{},
	)
}

type logWrapper struct {
	l *unilogger.Logger
	p predicate.Predicate
}

func (w logWrapper) Create(createEvent event.CreateEvent) bool {
	logEntry := w.l.With("event", createEvent)
	defer w.recover(logEntry)

	result := w.p.Create(createEvent)
	logEntry.
		With("result", result).
		Debug("processed create event")

	return result
}

func (w logWrapper) Delete(deleteEvent event.DeleteEvent) bool {
	logEntry := w.l.With("event", deleteEvent)
	defer w.recover(logEntry)

	result := w.p.Delete(deleteEvent)
	logEntry.
		With("result", result).
		Debug("processed delete event")

	return result
}

func (w logWrapper) Update(updateEvent event.UpdateEvent) bool {
	logEntry := w.l.With("event", updateEvent)
	defer w.recover(logEntry)

	result := w.p.Update(updateEvent)
	logEntry.
		With("result", result).
		Debug("processed update event")

	return result
}

func (w logWrapper) Generic(genericEvent event.GenericEvent) bool {
	logEntry := w.l.With("event", genericEvent)
	defer w.recover(logEntry)

	result := w.p.Generic(genericEvent)
	logEntry.
		With("result", result).
		Debug("processed generic event")

	return result
}

func (w logWrapper) recover(logEntry *unilogger.Logger) {
	r := recover()
	if r == nil {
		return
	}

	logEntry.
		With("panic", r).
		With("stack", debug.Stack()).
		Error("recovered from panic")
}

type releasePhasePredicate struct{}

func (rp releasePhasePredicate) Create(ev event.CreateEvent) bool {
	if ev.Object == nil {
		return false
	}

	switch ev.Object.(*v1alpha1.DeckhouseRelease).Status.Phase {
	case v1alpha1.PhaseSkipped, v1alpha1.PhaseSuperseded, v1alpha1.PhaseSuspended, v1alpha1.PhaseDeployed:
		return false
	}
	return true
}

// Delete returns true if the Delete event should be processed
func (rp releasePhasePredicate) Delete(_ event.DeleteEvent) bool {
	return false
}

// Update returns true if the Update event should be processed
func (rp releasePhasePredicate) Update(ev event.UpdateEvent) bool {
	if ev.ObjectNew == nil {
		return false
	}

	switch ev.ObjectNew.(*v1alpha1.DeckhouseRelease).Status.Phase {
	case v1alpha1.PhaseSkipped, v1alpha1.PhaseSuperseded, v1alpha1.PhaseSuspended, v1alpha1.PhaseDeployed:
		return false
	}
	return true
}

// Generic returns true if the Generic event should be processed
func (rp releasePhasePredicate) Generic(_ event.GenericEvent) bool {
	return true
}
