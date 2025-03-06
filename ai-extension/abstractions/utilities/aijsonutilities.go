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

func HashDataToString(values ...interface{}) string {
	h := fnv.New64a()

	for _, v := range values {
		writeHash(h, v)
	}

	return strconv.FormatUint(h.Sum64(), 16)
}

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
