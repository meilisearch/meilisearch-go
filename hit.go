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

	switch rv.Kind() {
	case reflect.Struct:
		// existing struct path (keep your current implementation) ...
		ti := getTypeInfo(rv.Type())
		for key, raw := range h {
			if len(raw) == 0 {
				continue
			}
			idx, ok := ti.byNameIndex[key]
			if !ok {
				continue
			}
			f := ti.fields[idx]
			fv, ok := fieldByIndexPathAlloc(rv, f.indexPath)
			if !ok {
				continue
			}
			if isJSONNull(raw) {
				if fv.CanSet() {
					fv.Set(reflect.Zero(fv.Type()))
				}
				continue
			}
			if f.hasString {
				if err := unmarshalSingleField(rv.Addr().Interface(), f.jsonName, raw); err != nil {
					return fmt.Errorf("decode field %q: %w", key, err)
				}
				continue
			}
			if !fv.CanAddr() {
				continue
			}
			if err := json.Unmarshal(raw, fv.Addr().Interface()); err != nil {
				return fmt.Errorf("decode field %q: %w", key, err)
			}
		}
		return nil

	case reflect.Map:
		// NEW: map[string]T support
		if rv.Type().Key().Kind() != reflect.String {
			return fmt.Errorf("map key must be string, got %s", rv.Type().Key())
		}
		if rv.IsNil() {
			rv.Set(reflect.MakeMapWithSize(rv.Type(), len(h)))
		}
		elemT := rv.Type().Elem()
		for k, raw := range h {
			// For null values, set zero (nil for pointer/slice/map/interface, 0 for numbers/bool)
			if isJSONNull(raw) {
				rv.SetMapIndex(reflect.ValueOf(k), reflect.Zero(elemT))
				continue
			}
			elemV := reflect.New(elemT) // *elemT
			// If elemT is interface{}, this is *interface{} and is fine.
			if err := json.Unmarshal(raw, elemV.Interface()); err != nil {
				return fmt.Errorf("decode map value for key %q: %w", k, err)
			}
			// Store the concrete elem (dereference)
			rv.SetMapIndex(reflect.ValueOf(k), elemV.Elem())
		}
		return nil

	default:
		return fmt.Errorf("out must point to a struct or map, got %T", out)
	}
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
		for i := range h {
			elemPtr := reflect.New(elemType)
			if err := h[i].DecodeInto(elemPtr.Interface()); err != nil {
				return fmt.Errorf("decode hits[%d]: %w", i, err)
			}
			out = reflect.Append(out, elemPtr.Elem())
		}

	case reflect.Ptr:
		et := elemType.Elem()
		switch et.Kind() {
		case reflect.Struct:
			for i := range h {
				elemPtr := reflect.New(et)
				if err := h[i].DecodeInto(elemPtr.Interface()); err != nil {
					return fmt.Errorf("decode hits[%d]: %w", i, err)
				}
				out = reflect.Append(out, elemPtr.Convert(elemType))
			}
		case reflect.Map:
			// *** FIX 1: use et (map type), not elemType (pointer type)
			if et.Key().Kind() != reflect.String {
				return fmt.Errorf("slice element must be map with string key, got %s", et)
			}
			for i := range h {
				// allocate *map[string]V and initialize the map
				mPtr := reflect.New(et)              // *map[string]V
				mPtr.Elem().Set(reflect.MakeMap(et)) // init
				// decode into *map
				if err := h[i].DecodeInto(mPtr.Interface()); err != nil {
					return fmt.Errorf("decode hits[%d]: %w", i, err)
				}
				// append the pointer as required by [] *map[string]V
				out = reflect.Append(out, mPtr.Convert(elemType))
			}
		default:
			return fmt.Errorf("slice element must be struct, *struct, or *map[string]T, got %s", elemType)
		}

	case reflect.Map:
		if elemType.Key().Kind() != reflect.String {
			return fmt.Errorf("slice element must be map with string key, got %s", elemType)
		}
		for i := range h {
			// allocate *map[string]V and initialize
			mPtr := reflect.New(elemType)              // *map[string]V
			mPtr.Elem().Set(reflect.MakeMap(elemType)) // init
			// decode into *map
			if err := h[i].DecodeInto(mPtr.Interface()); err != nil {
				return fmt.Errorf("decode hits[%d]: %w", i, err)
			}
			// append the map value (dereferenced)
			out = reflect.Append(out, mPtr.Elem())
		}

	default:
		return fmt.Errorf("slice element must be struct, *struct, map[string]T, or *map[string]T, got %s", elemType)
	}

	sv.Set(out)
	return nil
}

type fieldMeta struct {
	jsonName  string
	indexPath []int
	hasString bool
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
		hasString := false
		if tag != "" {
			// split first token (name) and options
			if c := indexByte(tag, ','); c >= 0 {
				nameToken := tag[:c]
				if nameToken != "" {
					name = nameToken
				}
				if hasJSONTagOption(tag[c+1:], "string") {
					hasString = true
				}
			} else {
				// tag without options
				name = tag
			}
		}
		idx := append(append([]int(nil), prefix...), i)

		// Inline embedded struct or *struct (pointer-embedded) when not renamed.
		if sf.Anonymous && name == sf.Name {
			u := sf.Type
			if u.Kind() == reflect.Ptr {
				u = u.Elem()
			}
			if u.Kind() == reflect.Struct {
				out = append(out, collectFields(u, idx)...)
				continue
			}
		}

		out = append(out, fieldMeta{
			jsonName:  name,
			indexPath: idx,
			hasString: hasString,
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

// hasJSONTagOption reports whether the comma-separated tag options contain opt.
func hasJSONTagOption(opts, opt string) bool {
	// opts like: "omitempty,string"
	start := 0
	for i := 0; i <= len(opts); i++ {
		if i == len(opts) || opts[i] == ',' {
			if opts[start:i] == opt {
				return true
			}
			start = i + 1
		}
	}
	return false
}

// fieldByIndexPathAlloc walks from the addressable struct value rv to the field
// indicated by indexPath, allocating intermediate *struct fields if needed.
// Returns the leaf field value and whether it is usable (addressable/settable).
func fieldByIndexPathAlloc(rv reflect.Value, indexPath []int) (reflect.Value, bool) {
	cur := rv
	for _, idx := range indexPath {
		if cur.Kind() == reflect.Ptr {
			if cur.IsNil() {
				if !cur.CanSet() {
					return reflect.Value{}, false
				}
				if cur.Type().Elem().Kind() != reflect.Struct {
					return reflect.Value{}, false
				}
				cur.Set(reflect.New(cur.Type().Elem()))
			}
			cur = cur.Elem()
		}
		if cur.Kind() != reflect.Struct {
			return reflect.Value{}, false
		}

		if idx < 0 || idx >= cur.NumField() {
			return reflect.Value{}, false
		}

		cur = cur.Field(idx)
	}
	if cur.Kind() == reflect.Ptr && cur.IsNil() {
		if cur.CanSet() && cur.Type().Elem().Kind() == reflect.Struct {
			cur.Set(reflect.New(cur.Type().Elem()))
		}
	}
	return cur, cur.CanAddr() || cur.CanSet()
}

// unmarshalSingleField unmarshals a single field into dst (pointer to struct)
// by constructing a minimal {"name": raw} object and delegating to encoding/json.
// This preserves stdlib behaviors such as json:",string".
func unmarshalSingleField(dst any, name string, raw []byte) error {
	tmp := map[string]json.RawMessage{name: json.RawMessage(raw)}
	b, err := json.Marshal(tmp)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, dst)
}

func isJSONNull(b []byte) bool {
	return len(b) == 4 && b[0] == 'n' && b[1] == 'u' && b[2] == 'l' && b[3] == 'l'
}
