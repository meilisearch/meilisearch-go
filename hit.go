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
	raw, err := json.Marshal(h)
	if err != nil {
		return fmt.Errorf("failed to marshal hit: %w", err)
	}

	return json.Unmarshal(raw, vPtr)
}

func (h Hits) Len() int {
	return len(h)
}

// Decode decodes the Hits into the provided target slice.
func (h Hits) Decode(vPtr interface{}) error {
	v := reflect.ValueOf(vPtr)

	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Slice {
		return errors.New("v must be a pointer to a slice")
	}

	sliceVal := v.Elem()
	elemType := sliceVal.Type().Elem()

	// Pre-allocate slice capacity
	sliceVal.Set(reflect.MakeSlice(sliceVal.Type(), 0, len(h)))

	for i, hit := range h {
		elemPtr := reflect.New(elemType).Interface()

		if err := hit.Decode(elemPtr); err != nil {
			return fmt.Errorf("failed to decode hit at index %d: %w", i, err)
		}

		sliceVal.Set(reflect.Append(sliceVal, reflect.ValueOf(elemPtr).Elem()))
	}
	return nil
}
