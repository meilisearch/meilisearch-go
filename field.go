package meilisearch

import "time"

// IntPtr returns a pointer to the given int value.
func IntPtr(v int) *int { return &v }

// Int64Ptr returns a pointer to the given int64 value.
func Int64Ptr(v int64) *int64 { return &v }

// BoolPtr returns a pointer to the given bool value.
func BoolPtr(v bool) *bool { return &v }

// FloatPtr returns a pointer to the given float64 value.
func FloatPtr(v float64) *float64 { return &v }

// StringPtr returns a pointer to the given string value.
func StringPtr(v string) *string { return &v }

// TimePtr returns a pointer to the given time.Time value.
func TimePtr(v time.Time) *time.Time { return &v }
