package meilisearch

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInt(t *testing.T) {
	v := 42
	ptr := IntPtr(v)
	assert.NotNil(t, ptr, "IntPtr returned nil pointer")
	assert.Equal(t, v, *ptr, "IntPtr returned wrong value")
}

func TestInt64(t *testing.T) {
	v := int64(42)
	ptr := Int64Ptr(v)
	assert.NotNil(t, ptr, "Int64Ptr returned nil pointer")
	assert.Equal(t, v, *ptr, "Int64Ptr returned wrong value")
}

func TestBool(t *testing.T) {
	v := true
	ptr := BoolPtr(v)
	assert.NotNil(t, ptr, "BoolPtr returned nil pointer")
	assert.Equal(t, v, *ptr, "BoolPtr returned wrong value")
}

func TestFloat(t *testing.T) {
	v := 3.14
	ptr := FloatPtr(v)
	assert.NotNil(t, ptr, "FloatPtr returned nil pointer")
	assert.Equal(t, v, *ptr, "FloatPtr returned wrong value")
}

func TestString(t *testing.T) {
	v := "hello"
	ptr := StringPtr(v)
	assert.NotNil(t, ptr, "StringPtr returned nil pointer")
	assert.Equal(t, v, *ptr, "StringPtr returned wrong value")
}

func TestTime(t *testing.T) {
	v := time.Date(2020, 1, 2, 3, 4, 5, 6, time.UTC)
	ptr := TimePtr(v)
	assert.NotNil(t, ptr, "TimePtr returned nil pointer")
	assert.True(t, ptr.Equal(v), "TimePtr returned wrong value")
}
