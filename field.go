package meilisearch

import "time"

// Int returns a pointer to the given int value.
func Int(v int) *int { return &v }

// Int64 returns a pointer to the given int64 value.
func Int64(v int64) *int64 { return &v }

// Bool returns a pointer to the given bool value.
func Bool(v bool) *bool { return &v }

// Float returns a pointer to the given float64 value.
func Float(v float64) *float64 { return &v }

// String returns a pointer to the given string value.
func String(v string) *string { return &v }

// Time returns a pointer to the given time.Time value.
func Time(v time.Time) *time.Time { return &v }
