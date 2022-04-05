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

package probe

import (
	"time"

	"github.com/sirupsen/logrus"

	"d8.io/upmeter/pkg/check"
	"d8.io/upmeter/pkg/kubernetes"
	"d8.io/upmeter/pkg/probe/checker"
	"d8.io/upmeter/pkg/set"
)

func NewLoader(filter Filter, access kubernetes.Access, logger *logrus.Logger) *Loader {
	return &Loader{
		filter: filter,
		access: access,
		logger: logger,
	}
}

type Loader struct {
	filter Filter
	access kubernetes.Access
	logger *logrus.Logger

	groups []string
	probes []check.ProbeRef

	configs []runnerConfig
}

func (l *Loader) collectConfigs() []runnerConfig {
	if l.configs != nil {
		// Already inited
		return l.configs
	}

	l.configs = make([]runnerConfig, 0)
	l.configs = append(l.configs, initSynthetic(l.access, l.logger)...)
	l.configs = append(l.configs, initControlPlane(l.access)...)
	l.configs = append(l.configs, initMonitoringAndAutoscaling(l.access)...)
	l.configs = append(l.configs, initScaling(l.access)...)
	l.configs = append(l.configs, initLoadBalancing(l.access)...)
	l.configs = append(l.configs, initDeckhouse(l.access, l.logger)...)

	return l.configs
}

func (l *Loader) Load() []*check.Runner {
	runners := make([]*check.Runner, 0)
	for _, rc := range l.collectConfigs() {
		if !l.filter.Enabled(rc.Ref()) {
			continue
		}

		runnerLogger := l.logger.WithFields(map[string]interface{}{
			"group": rc.group,
			"probe": rc.probe,
			"check": rc.check,
		})

		runner := check.NewRunner(rc.group, rc.probe, rc.check, rc.period, rc.config.Checker(), runnerLogger)

		runners = append(runners, runner)
		l.logger.Infof("Register probe %s", runner.ProbeRef().Id())
	}

	return runners
}

func (l *Loader) Groups() []string {
	if l.groups != nil {
		return l.groups
	}

	groups := set.New()
	for _, rc := range l.collectConfigs() {
		if !l.filter.Enabled(rc.Ref()) {
			continue
		}
		groups.Add(rc.group)

	}

	l.groups = groups.Slice()
	return l.groups
}

func (l *Loader) Probes() []check.ProbeRef {
	if l.probes != nil {
		return l.probes
	}

	l.probes = make([]check.ProbeRef, 0)
	for _, rc := range l.collectConfigs() {
		ref := rc.Ref()
		if !l.filter.Enabled(ref) {
			continue
		}

		l.probes = append(l.probes, ref)

	}

	return l.probes
}

type runnerConfig struct {
	group  string
	probe  string
	check  string
	period time.Duration
	config checker.Config
}

func (rc runnerConfig) Ref() check.ProbeRef {
	return check.ProbeRef{Group: rc.group, Probe: rc.probe}
}

func NewProbeFilter(disabled []string) Filter {
	return Filter{refs: set.New(disabled...)}
}

type Filter struct {
	refs set.StringSet
}

func (f Filter) Enabled(ref check.ProbeRef) bool {
	return !(f.refs.Has(ref.Id()) || f.refs.Has(ref.Group) || f.refs.Has(ref.Group+"/"))
}
