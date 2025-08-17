package meilisearch

import (
	"testing"
	"time"
)

func TestInt(t *testing.T) {
	v := int(42)
	ptr := Int(v)
	if ptr == nil {
		t.Fatal("Int64 returned nil pointer")
	}
	if *ptr != v {
		t.Errorf("Int64 returned pointer to %d, want %d", *ptr, v)
	}
}

func TestInt64(t *testing.T) {
	v := int64(42)
	ptr := Int64(v)
	if ptr == nil {
		t.Fatal("Int64 returned nil pointer")
	}
	if *ptr != v {
		t.Errorf("Int64 returned pointer to %d, want %d", *ptr, v)
	}
}

func TestBool(t *testing.T) {
	v := true
	ptr := Bool(v)
	if ptr == nil {
		t.Fatal("Bool returned nil pointer")
	}
	if *ptr != v {
		t.Errorf("Bool returned pointer to %v, want %v", *ptr, v)
	}
}

func TestFloat(t *testing.T) {
	v := 3.14
	ptr := Float(v)
	if ptr == nil {
		t.Fatal("Float returned nil pointer")
	}
	if *ptr != v {
		t.Errorf("Float returned pointer to %f, want %f", *ptr, v)
	}
}

func TestString(t *testing.T) {
	v := "hello"
	ptr := String(v)
	if ptr == nil {
		t.Fatal("String returned nil pointer")
	}
	if *ptr != v {
		t.Errorf("String returned pointer to %s, want %s", *ptr, v)
	}
}

func TestTime(t *testing.T) {
	v := time.Date(2020, 1, 2, 3, 4, 5, 6, time.UTC)
	ptr := Time(v)
	if ptr == nil {
		t.Fatal("Time returned nil pointer")
	}
	if !ptr.Equal(v) {
		t.Errorf("Time returned pointer to %v, want %v", *ptr, v)
	}
}
