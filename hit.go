package meilisearch

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sync"
)

type (
	Hit  map[string]json.RawMessage // Hit is a map of key and value raw buffer
	Hits []Hit                      // Hits is an alias for a slice of Hit.
)

// Deprecated: Decode decodes a single Hit into the provided struct.
//
// Please use DecodeInto for better performance without intermediate marshaling.
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

// DecodeInto decodes a single Hit into the provided struct without intermediate marshaling.
func (h Hit) DecodeInto(out any) error {
	if out == nil {
		return errors.New("out must be a non-nil pointer")
	}
	rv := reflect.ValueOf(out)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("out must be a non-nil pointer")
	}
	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("out must point to a struct, got %T", out)
	}

	ti := getTypeInfo(rv.Type())

	// iterate only present json keys in the hit
	for key, raw := range h {
		if len(raw) == 0 || isJSONNull(raw) {
			continue
		}
		idx, ok := ti.byNameIndex[key]
		if !ok {
			// unknown field: ignore (or collect for a Strict mode)
			continue
		}
		f := ti.fields[idx]
		fv := rv.FieldByIndex(f.indexPath)
		if !fv.CanAddr() {
			continue
		}
		if err := json.Unmarshal(raw, fv.Addr().Interface()); err != nil {
			return fmt.Errorf("decode field %q: %w", key, err)
		}
	}
	return nil
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

// Deprecated: Decode decodes the Hits into the provided target slice.
//
// Please use DecodeInto for better performance without intermediate marshaling.
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

// DecodeInto decodes hs into the provided slice pointer without re-marshal.
// vSlicePtr must be a non-nil pointer to a slice whose element type is a struct or *struct.
// Example:
//
//	var out []exampleBookForTest
//	if err := hits.DecodeInto(&out); err != nil { ... }
//
//	var outPtr []*exampleBookForTest
//	if err := hits.DecodeInto(&outPtr); err != nil { ... }
func (h Hits) DecodeInto(vSlicePtr interface{}) error {
	if vSlicePtr == nil {
		return fmt.Errorf("vSlicePtr must be a non-nil pointer to a slice")
	}

	rv := reflect.ValueOf(vSlicePtr)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("vSlicePtr must be a non-nil pointer, got %T", vSlicePtr)
	}

	sv := rv.Elem()
	if sv.Kind() != reflect.Slice {
		return fmt.Errorf("vSlicePtr must point to a slice, got %s", sv.Kind())
	}

	elemType := sv.Type().Elem()
	out := reflect.MakeSlice(sv.Type(), 0, len(h))

	switch elemType.Kind() {
	case reflect.Struct:
		// Target is []S
		for i := range h {
			elemPtr := reflect.New(elemType) // *S
			if err := h[i].DecodeInto(elemPtr.Interface()); err != nil {
				return fmt.Errorf("decode hits[%d]: %w", i, err)
			}
			out = reflect.Append(out, elemPtr.Elem()) // S
		}

	case reflect.Ptr:
		// Target is []*S
		if elemType.Elem().Kind() != reflect.Struct {
			return fmt.Errorf("slice element must be struct or *struct, got %s", elemType)
		}
		for i := range h {
			elemPtr := reflect.New(elemType.Elem()) // *S
			if err := h[i].DecodeInto(elemPtr.Interface()); err != nil {
				return fmt.Errorf("decode hits[%d]: %w", i, err)
			}
			out = reflect.Append(out, elemPtr.Convert(elemType)) // *S
		}

	default:
		return fmt.Errorf("slice element must be struct or *struct, got %s", elemType)
	}

	sv.Set(out)
	return nil
}

type fieldMeta struct {
	jsonName  string
	indexPath []int
}

type typeInfo struct {
	fields      []fieldMeta
	byNameIndex map[string]int // json name -> index in fields
}

var (
	typeInfoCache sync.Map // reflect.Type -> *typeInfo
)

func getTypeInfo(t reflect.Type) *typeInfo {
	if v, ok := typeInfoCache.Load(t); ok {
		return v.(*typeInfo)
	}
	f := collectFields(t, nil)
	ti := &typeInfo{
		fields:      f,
		byNameIndex: make(map[string]int, len(f)),
	}
	for i := range f {
		ti.byNameIndex[f[i].jsonName] = i
	}
	typeInfoCache.Store(t, ti)
	return ti
}

func collectFields(t reflect.Type, prefix []int) []fieldMeta {
	var out []fieldMeta

	// Walk struct fields, including anonymous/embedded.
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		// Skip unexported (PkgPath != "", except for embedded anonymous exported types)
		if sf.PkgPath != "" && !sf.Anonymous {
			continue
		}
		tag := sf.Tag.Get("json")
		if tag == "-" {
			continue
		}
		name := sf.Name
		if tag != "" {
			// take first token before comma
			if c := indexByte(tag, ','); c >= 0 {
				tag = tag[:c]
			}
			if tag != "" {
				name = tag
			}
		}
		idx := append(append([]int(nil), prefix...), i)

		if sf.Anonymous && sf.Type.Kind() == reflect.Struct && name == sf.Name {
			// inline embedded struct fields (respect tags if present)
			out = append(out, collectFields(sf.Type, idx)...)
			continue
		}

		out = append(out, fieldMeta{
			jsonName:  name,
			indexPath: idx,
		})
	}
	return out
}

// small helper: faster than strings.IndexByte here, avoids import
func indexByte(s string, c byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return i
		}
	}
	return -1
}

func isJSONNull(b []byte) bool {
	return len(b) == 4 && b[0] == 'n' && b[1] == 'u' && b[2] == 'l' && b[3] == 'l'
}
