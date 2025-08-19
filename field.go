package meilisearch

import (
	"encoding/json"
	"time"
)

// Int creates an Opt[int] with the given int value and sets its status to included.
func Int(v int) Opt[int] { return NewOpt(v) }

// Int64 creates an Opt[int64] with the given int64 value and sets its status to included.
func Int64(v int64) Opt[int64] { return NewOpt(v) }

// Bool creates an Opt[bool] with the given bool value and sets its status to included.
func Bool(v bool) Opt[bool] { return NewOpt(v) }

// Float creates an Opt[float64] with the given float64 value and sets its status to included.
func Float(v float64) Opt[float64] { return NewOpt(v) }

// String creates an Opt[string] with the given string value and sets its status to included.
func String(v string) Opt[string] { return NewOpt(v) }

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

type status int8

const (
	omitted status = iota
	null
	included
)

type Opt[T any] struct {
	Value  T
	status status
}

// NewOpt creates an Opt[T] with the given value and sets its status to included.
func NewOpt[T any](v T) Opt[T] { return Opt[T]{Value: v, status: included} }

// Null creates an Opt[T] with a null status.
func Null[T any]() Opt[T] { return Opt[T]{status: null} }

// Valid returns true if the Opt[T] is included (i.e., contains a value).
func (o Opt[T]) Valid() bool { return o.status == included }

// Null returns true if the Opt[T] has a null status.
func (o Opt[T]) Null() bool { return o.status == null }

func (o Opt[T]) MarshalJSON() ([]byte, error) {
	if o.status != included {
		return []byte("null"), nil
	}
	return json.Marshal(o.Value)
}

func (o *Opt[T]) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		o.status = null
		var z T
		o.Value = z
		return nil
	}
	o.status = included
	return json.Unmarshal(b, &o.Value)
}
