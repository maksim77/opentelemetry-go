// Copyright The OpenTelemetry Authors
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

package prometheus // import "go.opentelemetry.io/otel/exporters/prometheus"

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/view"
)

func TestNewConfig(t *testing.T) {
	registry := prometheus.NewRegistry()

	aggregationSelector := func(view.InstrumentKind) aggregation.Aggregation { return nil }

	testCases := []struct {
		name            string
		options         []Option
		wantRegisterer  prometheus.Registerer
		wantAggregation metric.AggregationSelector
	}{
		{
			name:           "Default",
			options:        nil,
			wantRegisterer: prometheus.DefaultRegisterer,
		},
		{
			name: "WithRegisterer",
			options: []Option{
				WithRegisterer(registry),
			},
			wantRegisterer: registry,
		},
		{
			name: "WithAggregationSelector",
			options: []Option{
				WithAggregationSelector(aggregationSelector),
			},
			wantRegisterer:  prometheus.DefaultRegisterer,
			wantAggregation: aggregationSelector,
		},
		{
			name: "With Multiple Options",
			options: []Option{
				WithRegisterer(registry),
				WithAggregationSelector(aggregationSelector),
			},
			wantRegisterer:  registry,
			wantAggregation: aggregationSelector,
		},
		{
			name: "nil options do nothing",
			options: []Option{
				WithRegisterer(nil),
			},
			wantRegisterer: prometheus.DefaultRegisterer,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			cfg := newConfig(tt.options...)

			assert.Equal(t, tt.wantRegisterer, cfg.registerer)
		})
	}
}

func TestConfigManualReaderOptions(t *testing.T) {
	aggregationSelector := func(view.InstrumentKind) aggregation.Aggregation { return nil }

	testCases := []struct {
		name            string
		config          config
		wantOptionCount int
	}{
		{
			name:            "Default",
			config:          config{},
			wantOptionCount: 0,
		},

		{
			name:            "WithAggregationSelector",
			config:          config{aggregation: aggregationSelector},
			wantOptionCount: 1,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			opts := tt.config.manualReaderOptions()
			assert.Len(t, opts, tt.wantOptionCount)
		})
	}
}
