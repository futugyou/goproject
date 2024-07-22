package convert

import (
	"errors"
	"math"
	"strconv"
)

type StrTo string

func (s StrTo) String() string {
	return string(s)
}

func (s StrTo) Int() (int, error) {
	i64, err := strconv.ParseInt(s.String(), 10, 0) // 0 specifies that it will infer the size from int
	if err != nil {
		return 0, err
	}

	// Determine the range for int based on the system architecture
	var intMin, intMax int64
	if strconv.IntSize == 32 {
		intMin = math.MinInt32
		intMax = math.MaxInt32
	} else {
		intMin = math.MinInt64
		intMax = math.MaxInt64
	}

	if i64 < intMin || i64 > intMax {
		return 0, errors.New("value out of int range")
	}
	return int(i64), nil
}

func (s StrTo) MustInt() int {
	v, _ := s.Int()
	return v
}

func (s StrTo) UInt32() (uint32, error) {
	u64, err := strconv.ParseUint(s.String(), 10, 32)
	if err != nil {
		return 0, err
	}
	if u64 > uint64(math.MaxUint32) {
		return 0, errors.New("value out of uint32 range")
	}
	return uint32(u64), nil
}

func (s StrTo) MustUInt32() uint32 {
	v, _ := s.UInt32()
	return v
}
