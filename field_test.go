package meilisearch

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInt(t *testing.T) {
	v := 42
	ptr := Int(v)
	assert.NotNil(t, ptr, "Int returned nil pointer")
	assert.Equal(t, v, *ptr, "Int returned wrong value")
}

func TestInt64(t *testing.T) {
	v := int64(42)
	ptr := Int64(v)
	assert.NotNil(t, ptr, "Int64 returned nil pointer")
	assert.Equal(t, v, *ptr, "Int64 returned wrong value")
}

func TestBool(t *testing.T) {
	v := true
	ptr := Bool(v)
	assert.NotNil(t, ptr, "Bool returned nil pointer")
	assert.Equal(t, v, *ptr, "Bool returned wrong value")
}

func TestFloat(t *testing.T) {
	v := 3.14
	ptr := Float(v)
	assert.NotNil(t, ptr, "Float returned nil pointer")
	assert.Equal(t, v, *ptr, "Float returned wrong value")
}

func TestString(t *testing.T) {
	v := "hello"
	ptr := String(v)
	assert.NotNil(t, ptr, "String returned nil pointer")
	assert.Equal(t, v, *ptr, "String returned wrong value")
}

func TestTime(t *testing.T) {
	v := time.Date(2020, 1, 2, 3, 4, 5, 6, time.UTC)
	ptr := Time(v)
	assert.NotNil(t, ptr, "Time returned nil pointer")
	assert.True(t, ptr.Equal(v), "Time returned wrong value")
}
