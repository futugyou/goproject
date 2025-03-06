package utilities

import (
	"encoding/json"
	"fmt"
	"hash"
	"hash/fnv"
	"reflect"
	"sort"
	"strconv"
)

// HashDataToString generates a hash string from the given values using FNV-1a algorithm.
// It accepts variable number of interface{} arguments and returns their combined hash as a hexadecimal string.
func HashDataToString(values ...interface{}) string {
	h := fnv.New64a()

	for _, v := range values {
		writeHash(h, v)
	}

	return strconv.FormatUint(h.Sum64(), 16)
}

// writeHash writes the hash representation of a value v into the provided hash.Hash64.
// It handles various basic types (string, numeric types, bool) and fmt.Stringer interface.
// For complex types, it delegates to encodeComplex.
// If the input value is nil, it writes "nil" to the hash.
func writeHash(h hash.Hash64, v interface{}) {
	if v == nil {
		h.Write([]byte("nil"))
		return
	}

	switch val := v.(type) {
	case string:
		h.Write([]byte(val))
	case int, int8, int16, int32, int64:
		h.Write([]byte(strconv.FormatInt(reflect.ValueOf(val).Int(), 10)))
	case uint, uint8, uint16, uint32, uint64:
		h.Write([]byte(strconv.FormatUint(reflect.ValueOf(val).Uint(), 10)))
	case float32, float64:
		h.Write([]byte(strconv.FormatFloat(reflect.ValueOf(val).Float(), 'f', -1, 64)))
	case bool:
		if val {
			h.Write([]byte("true"))
		} else {
			h.Write([]byte("false"))
		}
	case fmt.Stringer:
		h.Write([]byte(val.String()))
	default:
		encodeComplex(h, v)
	}
}

// encodeComplex writes a hash representation of a complex data structure into the provided hash.Hash64.
// It handles various types including pointers, slices, arrays, maps, structs, and primitive types.
// For maps, it ensures consistent hashing by sorting keys before processing.
// For structs, it uses JSON marshaling to create a consistent string representation.
func encodeComplex(h hash.Hash64, v interface{}) {
	rv := reflect.ValueOf(v)

	switch rv.Kind() {
	case reflect.Ptr:
		if rv.IsNil() {
			h.Write([]byte("nil"))
		} else {
			writeHash(h, rv.Elem().Interface())
		}
	case reflect.Slice, reflect.Array:
		h.Write([]byte("["))
		for i := 0; i < rv.Len(); i++ {
			writeHash(h, rv.Index(i).Interface())
			h.Write([]byte(","))
		}
		h.Write([]byte("]"))
	case reflect.Map:
		keys := rv.MapKeys()
		sort.Slice(keys, func(i, j int) bool {
			return fmt.Sprint(keys[i].Interface()) < fmt.Sprint(keys[j].Interface())
		})

		h.Write([]byte("{"))
		for _, k := range keys {
			writeHash(h, k.Interface())
			h.Write([]byte(":"))
			writeHash(h, rv.MapIndex(k).Interface())
			h.Write([]byte(","))
		}
		h.Write([]byte("}"))
	case reflect.Struct:
		jsonBytes, _ := json.Marshal(v)
		h.Write(jsonBytes)
	default:
		h.Write([]byte(fmt.Sprintf("%v", v)))
	}
}
