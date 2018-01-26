// Copyright 2018 Istio Authors.
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

//go:generate $GOPATH/src/istio.io/contrib/bin/codegen.sh -f adapters/istio/mixer/example/config/config.proto
//go:generate $GOPATH/src/istio.io/contrib/bin/codegen.sh -t adapters/istio/mixer/example/template/template.proto

// Package tmpl provides an tmpl Mixer adapter implementation. It is used
// to provide code samples for future adapter developers on how to structure
// their code and build scripts. It is also intended to be used in automated
// Mixer testing to ensure that Mixer can be built and extended by adapters in
// the extensions repo (as well as other external repos).
//
// The adapter implementation in this package handles instances from existing
// Mixer templates: checknothing and reportnothing. For each
// template, the adapter logs the request, including provided params.
package example

import (
	"context"
	"time"

	"istio.io/istio/mixer/pkg/adapter"
	"istio.io/istio/mixer/pkg/status"
	"istio.io/istio/mixer/template/checknothing"
	"istio.io/istio/mixer/template/reportnothing"

	"istio.io/contrib/adapters/istio/mixer/example/config"
	tmpl "istio.io/contrib/adapters/istio/mixer/example/template"
)

const (
	validDuration = 10 * time.Second
	validUses     = 1000
)

type (
	builder struct {
		cfg *config.Params
	}

	handler struct {
		prefix string
		log    adapter.Logger
	}
)

// Ensure support all required interfaces.
var (
	_ tmpl.HandlerBuilder = &builder{}
	_ tmpl.Handler        = &handler{}

	_ checknothing.HandlerBuilder = &builder{}
	_ checknothing.Handler        = &handler{}

	_ reportnothing.HandlerBuilder = &builder{}
	_ reportnothing.Handler        = &handler{}
)

// GetInfo returns the adapter.Info associated with this implementation.
func GetInfo() adapter.Info {
	return adapter.Info{
		Name:        "example",
		Impl:        "istio.io/extensions/adapter/mixer/example",
		Description: "An example mixer adapter to showcase adapter development.",
		SupportedTemplates: []string{
			checknothing.TemplateName,
			reportnothing.TemplateName,
			tmpl.TemplateName,
		},
		DefaultConfig: &config.Params{LogPrefix: "mixer example adapter"},
		NewBuilder:    func() adapter.HandlerBuilder { return &builder{} },
	}
}

func (*builder) SetCheckNothingTypes(map[string]*checknothing.Type)   {}
func (*builder) SetReportNothingTypes(map[string]*reportnothing.Type) {}
func (*builder) SetExampleReportTypes(map[string]*tmpl.Type)          {}
func (b *builder) SetAdapterConfig(cfg adapter.Config)                { b.cfg = cfg.(*config.Params) }
func (*builder) Validate() (ce *adapter.ConfigErrors)                 { return }

func (b *builder) Build(context context.Context, env adapter.Env) (adapter.Handler, error) {
	return &handler{prefix: b.cfg.LogPrefix, log: env.Logger()}, nil
}

func (h *handler) HandleCheckNothing(ctx context.Context, instance *checknothing.Instance) (adapter.CheckResult, error) {
	h.log.Infof("%s: received checknothing instance: %#v", h.prefix, instance)
	return adapter.CheckResult{status.OK, validDuration, validUses}, nil
}

func (h *handler) HandleReportNothing(ctx context.Context, instances []*reportnothing.Instance) error {
	for _, instance := range instances {
		h.log.Infof("%s: received reportnothing instance: %#v", h.prefix, instance)
	}
	return nil
}

func (h *handler) HandleExampleReport(ctx context.Context, instances []*tmpl.Instance) error {
	for _, instance := range instances {
		h.log.Infof("%s: received example instance: %#v", h.prefix, instance)
	}
	return nil
}

func (*handler) Close() error { return nil }
