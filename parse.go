package chart

import (
	"strconv"
	"strings"
	"time"
)

// ParseFloats parses a list of floats.
func ParseFloats(values ...string) ([]float64, error) {
	output := make([]float64, 0, len(values))
	var parsedValue float64
	var err error
	var cleaned string
	for _, value := range values {
		cleaned = strings.TrimSpace(strings.ReplaceAll(value, ",", ""))
		if cleaned == "" {
			continue
		}
		if parsedValue, err = strconv.ParseFloat(cleaned, 64); err != nil {
			return nil, err
		}
		output = append(output, parsedValue)
	}
	return output, nil
}

// ParseTimes parses a list of times with a given format.
func ParseTimes(layout string, values ...string) ([]time.Time, error) {
	output := make([]time.Time, 0, len(values))
	var parsedValue time.Time
	var err error
	for _, value := range values {
		if parsedValue, err = time.Parse(layout, value); err != nil {
			return nil, err
		}
		output = append(output, parsedValue)
	}
	return output, nil
}
