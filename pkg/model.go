package pkg

import (
	"time"
)

type Metrics struct {
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
func generateMetrics(timestamp time.Time) Metrics {
	return Metrics{
		Timestamp:         timestamp,
		MetricGroup:       "sample_metric_group",
		NumberFieldKeys:   []string{"number_field_key_1", "number_field_key_2"},
		NumberFieldValues: []float64{1.23, 4.56},
		StringFieldKeys:   []string{"string_field_key_1", "string_field_key_2"},
		StringFieldValues: []string{"value1", "value2"},
		TagKeys:           []string{"tag_key_1", "tag_key_2"},
		TagValues:         []string{"tag_value_1", "tag_value_2"},
	}
}
