package meilisearch

import (
	"encoding/json"
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

func TestIntOpt(t *testing.T) {
	v := 5
	o := Int(v)
	assert.True(t, o.Valid(), "Int() should produce a valid Opt")
	assert.False(t, o.Null(), "Int() Opt should not be null")
	assert.Equal(t, v, o.Value, "Int() Opt has wrong value")
	b, err := json.Marshal(o)
	assert.NoError(t, err)
	assert.Equal(t, "5", string(b))
}

func TestInt64Opt(t *testing.T) {
	v := int64(123456789)
	o := Int64(v)
	assert.True(t, o.Valid())
	assert.False(t, o.Null())
	assert.Equal(t, v, o.Value)
	b, err := json.Marshal(o)
	assert.NoError(t, err)
	assert.Equal(t, "123456789", string(b))
}

func TestBoolOpt(t *testing.T) {
	v := true
	o := Bool(v)
	assert.True(t, o.Valid())
	assert.False(t, o.Null())
	assert.Equal(t, v, o.Value)
	b, err := json.Marshal(o)
	assert.NoError(t, err)
	assert.Equal(t, "true", string(b))
}

func TestFloatOpt(t *testing.T) {
	v := 3.14159
	o := Float(v)
	assert.True(t, o.Valid())
	assert.False(t, o.Null())
	assert.Equal(t, v, o.Value)
	b, err := json.Marshal(o)
	assert.NoError(t, err)
	assert.Contains(t, string(b), "3.14159") // float formatting
}

func TestStringOpt(t *testing.T) {
	v := "hello"
	o := String(v)
	assert.True(t, o.Valid())
	assert.False(t, o.Null())
	assert.Equal(t, v, o.Value)
	b, err := json.Marshal(o)
	assert.NoError(t, err)
	assert.Equal(t, "\"hello\"", string(b))
}

func TestNullOpt(t *testing.T) {
	o := Null[int]()
	assert.False(t, o.Valid(), "Null Opt should not be valid")
	assert.True(t, o.Null(), "Null Opt should report Null() = true")
	b, err := json.Marshal(o)
	assert.NoError(t, err)
	assert.Equal(t, "null", string(b))
}

func TestDefaultOptIsOmitted(t *testing.T) {
	var o Opt[int]
	assert.False(t, o.Valid(), "default Opt should not be valid (omitted)")
	assert.False(t, o.Null(), "default Opt should not be null; it's omitted")
	b, err := json.Marshal(o)
	assert.NoError(t, err)
	assert.Equal(t, "null", string(b), "omitted should marshal as null")
}

func TestUnmarshalIncluded(t *testing.T) {
	var o Opt[int]
	err := json.Unmarshal([]byte("42"), &o)
	assert.NoError(t, err)
	assert.True(t, o.Valid())
	assert.False(t, o.Null())
	assert.Equal(t, 42, o.Value)
}

func TestUnmarshalNull(t *testing.T) {
	var o Opt[string]
	o = String("preset") // start included with non-zero value to ensure reset
	assert.True(t, o.Valid())
	err := json.Unmarshal([]byte("null"), &o)
	assert.NoError(t, err)
	assert.False(t, o.Valid())
	assert.True(t, o.Null())
	assert.Equal(t, "", o.Value, "value should be zero value after null unmarshal")
}

func TestReUnmarshalOverwrite(t *testing.T) {
	var o Opt[int]
	err := json.Unmarshal([]byte("10"), &o)
	assert.NoError(t, err)
	assert.Equal(t, 10, o.Value)
	assert.True(t, o.Valid())
	// Overwrite with a new value
	err = json.Unmarshal([]byte("20"), &o)
	assert.NoError(t, err)
	assert.Equal(t, 20, o.Value)
	assert.True(t, o.Valid())
	// Overwrite with null
	err = json.Unmarshal([]byte("null"), &o)
	assert.NoError(t, err)
	assert.False(t, o.Valid())
	assert.True(t, o.Null())
}
