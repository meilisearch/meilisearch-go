package meilisearch

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

type (
	Hit  map[string]json.RawMessage // Hit is a map of key and value raw buffer
	Hits []Hit                      // Hits is an alias for a slice of Hit.
)

// Decode decodes a single Hit into the provided struct.
func (h Hit) Decode(vPtr interface{}) error {
	if vPtr == nil || reflect.ValueOf(vPtr).Kind() != reflect.Ptr {
		return errors.New("vPtr must be a non-nil pointer")
	}

	raw, err := json.Marshal(h)
	if err != nil {
		return fmt.Errorf("failed to marshal hit: %w", err)
	}

	return json.Unmarshal(raw, vPtr)
}

// DecodeWith decodes a Hit into the provided struct using the provided marshal and unmarshal functions.
func (h Hit) DecodeWith(vPtr interface{}, marshal JSONMarshal, unmarshal JSONUnmarshal) error {
	if vPtr == nil || reflect.ValueOf(vPtr).Kind() != reflect.Ptr {
		return errors.New("vPtr must be a non-nil pointer")
	}

	raw, err := marshal(h)
	if err != nil {
		return fmt.Errorf("failed to marshal hit: %w", err)
	}

	return unmarshal(raw, vPtr)
}

func (h Hits) Len() int {
	return len(h)
}

// Decode decodes the Hits into the provided target slice.
func (h Hits) Decode(vSlicePtr interface{}) error {
	v := reflect.ValueOf(vSlicePtr)

	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("v must be a pointer to a slice, got %T", vSlicePtr)
	}

	raw := []Hit(h)
	buf, err := json.Marshal(raw)
	if err != nil {
		return fmt.Errorf("failed to marshal hits: %w", err)
	}

	return json.Unmarshal(buf, vSlicePtr)
}

// DecodeWith decodes a Hits into the provided struct using the provided marshal and unmarshal functions.
func (h Hits) DecodeWith(vSlicePtr interface{}, marshal JSONMarshal, unmarshal JSONUnmarshal) error {
	v := reflect.ValueOf(vSlicePtr)

	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Slice {
		return errors.New("v must be a pointer to a slice")
	}

	raw := []Hit(h)
	buf, err := marshal(raw)
	if err != nil {
		return fmt.Errorf("failed to marshal hit: %w", err)
	}

	return unmarshal(buf, vSlicePtr)
}
