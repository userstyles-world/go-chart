package chart

import (
	"fmt"
	"strconv"
	"time"
)

// ValueFormatter is a function that takes a value and produces a string.
type ValueFormatter func(v interface{}) string

var (
	// Unix to string
	cachingTimeFormat = make(map[int64]string)
)

// TimeValueFormatter is a ValueFormatter for timestamps.
func TimeValueFormatter(v interface{}) string {
	if typed, isTyped := v.(time.Time); isTyped {
		if ccache, ok := cachingTimeFormat[typed.Unix()]; ok {
			return ccache
		}
		result := typed.Format(DefaultDateFormat)
		cachingTimeFormat[typed.Unix()] = result
		return result
	}
	if typed, isTyped := v.(int64); isTyped {
		unix := time.Unix(0, typed)
		if ccache, ok := cachingTimeFormat[unix.Unix()]; ok {
			return ccache
		}
		result := unix.Format(DefaultDateFormat)
		cachingTimeFormat[unix.Unix()] = result
		return result
	}
	if typed, isTyped := v.(float64); isTyped {
		unix := time.Unix(0, int64(typed))
		if ccache, ok := cachingTimeFormat[unix.Unix()]; ok {
			return ccache
		}
		result := unix.Format(DefaultDateFormat)
		cachingTimeFormat[unix.Unix()] = result
		return result
	}
	return ""
}

// TimeHourValueFormatter is a ValueFormatter for timestamps.
func TimeHourValueFormatter(v interface{}) string {
	return formatTime(v, DefaultDateHourFormat)
}

// TimeMinuteValueFormatter is a ValueFormatter for timestamps.
func TimeMinuteValueFormatter(v interface{}) string {
	return formatTime(v, DefaultDateMinuteFormat)
}

// TimeDateValueFormatter is a ValueFormatter for timestamps.
func TimeDateValueFormatter(v interface{}) string {
	return formatTime(v, "2006-01-02")
}

// TimeValueFormatterWithFormat returns a time formatter with a given format.
func TimeValueFormatterWithFormat(format string) ValueFormatter {
	return func(v interface{}) string {
		return formatTime(v, format)
	}
}

// TimeValueFormatterWithFormat is a ValueFormatter for timestamps with a given format.
func formatTime(v interface{}, dateFormat string) string {
	if typed, isTyped := v.(time.Time); isTyped {
		return typed.Format(dateFormat)
	}
	if typed, isTyped := v.(int64); isTyped {
		return time.Unix(0, typed).Format(dateFormat)
	}
	if typed, isTyped := v.(float64); isTyped {
		return time.Unix(0, int64(typed)).Format(dateFormat)
	}
	return ""
}

// IntValueFormatter is a ValueFormatter for float64.
func IntValueFormatter(v interface{}) string {
	switch v := v.(type) {
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case float32:
		return strconv.FormatInt(int64(v), 10)
	case float64:
		return strconv.FormatInt(int64(v), 10)
	default:
		return ""
	}
}

// FloatValueFormatter is a ValueFormatter for float64.
func FloatValueFormatter(v interface{}) string {
	if typed, isTyped := v.(int); isTyped {
		return ftoa2(float64(typed))
	}
	if typed, isTyped := v.(int64); isTyped {
		return ftoa2(float64(typed))
	}
	if typed, isTyped := v.(float32); isTyped {
		return ftoa2(float64(typed))
	}
	if typed, isTyped := v.(float64); isTyped {
		return ftoa2(typed)
	}
	return ""
}

// PercentValueFormatter is a formatter for percent values.
// NOTE: it normalizes the values, i.e. multiplies by 100.0.
func PercentValueFormatter(v interface{}) string {
	if typed, isTyped := v.(float64); isTyped {
		return FloatValueFormatterWithFormat(typed*100.0, DefaultPercentValueFormat)
	}
	return ""
}

// FloatValueFormatterWithFormat is a ValueFormatter for float64 with a given format.
func FloatValueFormatterWithFormat(v interface{}, floatFormat string) string {
	if typed, isTyped := v.(int); isTyped {
		return fmt.Sprintf(floatFormat, float64(typed))
	}
	if typed, isTyped := v.(int64); isTyped {
		return fmt.Sprintf(floatFormat, float64(typed))
	}
	if typed, isTyped := v.(float32); isTyped {
		return fmt.Sprintf(floatFormat, typed)
	}
	if typed, isTyped := v.(float64); isTyped {
		return fmt.Sprintf(floatFormat, typed)
	}
	return ""
}

// KValueFormatter is a formatter for K values.
func KValueFormatter(k float64, vf ValueFormatter) ValueFormatter {
	return func(v interface{}) string {
		return fmt.Sprintf("%0.0fσ %s", k, vf(v))
	}
}
