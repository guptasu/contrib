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

package example

import (
	"context"
	"strings"
	"testing"

	tmpl "istio.io/contrib/adapters/istio/mixer/example/template"
	"istio.io/istio/mixer/pkg/adapter/test"
	"istio.io/istio/mixer/pkg/status"
	"istio.io/istio/mixer/template/checknothing"
	"istio.io/istio/mixer/template/reportnothing"
)

func TestValidateBuild(t *testing.T) {
	b := GetInfo().NewBuilder().(*builder)
	b.SetAdapterConfig(GetInfo().DefaultConfig)

	if err := b.Validate(); err != nil {
		t.Errorf("Validate() resulted in unexpected error: %v", err)
	}

	// invoke the empty set methods for coverage
	b.SetCheckNothingTypes(map[string]*checknothing.Type{})
	b.SetReportNothingTypes(map[string]*reportnothing.Type{})
	b.SetExampleReportTypes(map[string]*tmpl.Type{})

	if _, err := b.Build(context.Background(), test.NewEnv(t)); err != nil {
		t.Errorf("Build() resulted in unexpected error: %v", err)
	}
}

func TestHandler(t *testing.T) {
	ctx := context.Background()
	testEnv := test.NewEnv(t)
	h := &handler{prefix: "Test", log: testEnv.Logger()}

	got, err := h.HandleCheckNothing(ctx, nil)
	if err != nil {
		t.Errorf("HandleCheckNothing(ctx, nil) resulted in an unexpected error: %v", err)
	}
	if !status.IsOK(got.Status) {
		t.Errorf("HandleCheckNothing(ctx, nil) => %#v, want %#v", got.Status, status.OK)
	}

	rni := &reportnothing.Instance{"rni"}
	if err := h.HandleReportNothing(ctx, []*reportnothing.Instance{rni}); err != nil {
		t.Errorf("HandleReportNothing(ctx, []{%#v}) resulted in an unexpected error: %v", rni, err)
	}

	ei := &tmpl.Instance{"ei", map[string]interface{}{"foo": "bar"}}
	if err := h.HandleExampleReport(ctx, []*tmpl.Instance{ei}); err != nil {
		t.Errorf("HandleExampleReport(ctx, []{%#v}) resulted in an unexpected error: %v", ei, err)
	}

	logs := testEnv.GetLogs()
	if len(logs) != 3 {
		t.Errorf("Got %d logs; want %d", len(logs), 3)
	}

	for _, l := range logs {
		if !strings.HasPrefix(l, "Test:") {
			t.Errorf("Log does not start with expected prefix, want 'Test:': %s", l)
		}
	}

	if err := h.Close(); err != nil {
		t.Errorf("Close() returned an unexpected error")
	}
}
