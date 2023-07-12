//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package pkg

import (
	"time"

	"github.com/google/uuid"
)

type Metric struct {
	Timestamp         time.Time `ch:"timestamp"`
	MetricGroup       string    `ch:"metric_group"`
	NumberFieldKeys   []string  `ch:"number_field_keys"`
	NumberFieldValues []float64 `ch:"number_field_values"`
	StringFieldKeys   []string  `ch:"string_field_keys"`
	StringFieldValues []string  `ch:"string_field_values"`
	TagKeys           []string  `ch:"tag_keys"`
	TagValues         []string  `ch:"tag_values"`
}

// generateMetrics generates a random Metrics object with the given timestamp

func generateMetric(timestamp time.Time, randomColumn bool) Metric {
	metric := Metric{
		Timestamp:         timestamp,
		MetricGroup:       "sample_metric_group",
		NumberFieldKeys:   []string{"number_field_key_1", "number_field_key_2"},
		NumberFieldValues: []float64{1.23, 4.56},
		StringFieldKeys:   []string{"string_field_key_1", "string_field_key_2"},
		StringFieldValues: []string{"value1", "value2"},
		TagKeys:           []string{"tag_key_1", "tag_key_2"},
		TagValues:         []string{"tag_value_1", "tag_value_2"},
	}

	if randomColumn {
		for i := 0; i < 10; i++ {
			guid := uuid.New()
			key := guid.String()

			metric.StringFieldKeys = append(metric.StringFieldKeys, key)
			metric.StringFieldValues = append(metric.StringFieldValues, key)
			metric.NumberFieldKeys = append(metric.NumberFieldKeys, key)
			metric.NumberFieldValues = append(metric.NumberFieldValues, float64(i))

			metric.TagKeys = append(metric.TagKeys, key)
			metric.TagValues = append(metric.TagValues, key)
		}
	}

	return metric
}
