package crdt

import (
	"reflect"
	"testing"
)

type decreasingInt int

func (i *decreasingInt) Merge(other interface{}) bool {
	otherInt := other.(decreasingInt)
	if otherInt < *i {
		*i = otherInt
		return true
	}
	return false
}

func TestMergeMerger(t *testing.T) {
	value := decreasingInt(0)
	testMerge := func(other decreasingInt, expectedChanged bool, expectedResult decreasingInt) {
		changed := Merge(&value, other)
		if changed != expectedChanged {
			t.Errorf("Merge(a, %#v) = %v, expected %v", other, changed, expectedChanged)
		}
		if value != expectedResult {
			t.Fatalf("After merge was %#v, expected %#v", value, expectedResult)
		}
	}
	testMerge(0, false, 0)
	testMerge(1, false, 0)
	testMerge(-1, true, -1)
	testMerge(2, false, -1)
	testMerge(-2, true, -2)
}

func TestMergeStruct(t *testing.T) {
	type A struct {
		I int
		J string
	}
	value := A{}
	testMerge := func(other A, expectedChanged bool, expectedResult A) {
		changed := Merge(&value, other)
		if changed != expectedChanged {
			t.Errorf("Merge(a, %#v) = %v, expected %v", other, changed, expectedChanged)
		}
		if value != expectedResult {
			t.Fatalf("After merge was %#v, expected %#v", value, expectedResult)
		}
	}
	testMerge(A{}, false, A{})
	testMerge(A{1, ""}, true, A{1, ""})
	testMerge(A{0, "a"}, true, A{1, "a"})
	testMerge(A{1, "a"}, false, A{1, "a"})
	testMerge(A{-1, "b"}, true, A{1, "b"})
	testMerge(A{2, "aa"}, true, A{2, "b"})
}

func TestMergeMap(t *testing.T) {
	type A map[int]int
	var value A
	testMerge := func(other A, expectedChanged bool, expectedResult A) {
		changed := Merge(&value, other)
		if changed != expectedChanged {
			t.Errorf("Merge(a, %#v) = %v, expected %v", other, changed, expectedChanged)
		}
		if !reflect.DeepEqual(value, expectedResult) {
			t.Fatalf("After merge was %#v, expected %#v", value, expectedResult)
		}
	}
	testMerge(A{}, false, A{})
	testMerge(A{1: 0}, true, A{1: 0})
	testMerge(A{1: 0}, false, A{1: 0})
	testMerge(A{2: 0}, true, A{1: 0, 2: 0})
	testMerge(A{2: 0}, false, A{1: 0, 2: 0})
	testMerge(A{1: 0, 2: 1}, true, A{1: 0, 2: 1})
	testMerge(A{1: 0, 2: 1}, false, A{1: 0, 2: 1})
	testMerge(A{1: 1, 2: 0}, true, A{1: 1, 2: 1})
	testMerge(A{1: 1, 2: 0}, false, A{1: 1, 2: 1})
}

func TestJoin(t *testing.T) {
	testJoin := func(a, b, expected interface{}) {
		if result := Join(a, b); !reflect.DeepEqual(result, expected) {
			t.Errorf("Join(%#v, %#v) = %#v, expected %#v", a, b, result, expected)
		}
	}
	// bool
	testJoin(false, false, false)
	testJoin(false, true, true)
	testJoin(true, false, true)
	testJoin(true, true, true)

	// int
	testJoin(int(0), int(1), int(1))
	testJoin(int(1), int(0), int(1))
	// int8
	testJoin(int8(0), int8(1), int8(1))
	testJoin(int8(1), int8(0), int8(1))
	// int16
	testJoin(int16(0), int16(1), int16(1))
	testJoin(int16(1), int16(0), int16(1))
	// int32
	testJoin(int32(0), int32(1), int32(1))
	testJoin(int32(1), int32(0), int32(1))
	// int64
	testJoin(int64(0), int64(1), int64(1))
	testJoin(int64(1), int64(0), int64(1))

	// uint
	testJoin(uint(0), uint(1), uint(1))
	testJoin(uint(1), uint(0), uint(1))
	// uint8
	testJoin(uint8(0), uint8(1), uint8(1))
	testJoin(uint8(1), uint8(0), uint8(1))
	// uint16
	testJoin(uint16(0), uint16(1), uint16(1))
	testJoin(uint16(1), uint16(0), uint16(1))
	// uint32
	testJoin(uint32(0), uint32(1), uint32(1))
	testJoin(uint32(1), uint32(0), uint32(1))
	// uint64
	testJoin(uint64(0), uint64(1), uint64(1))
	testJoin(uint64(1), uint64(0), uint64(1))

	// float32
	testJoin(float32(0), float32(1), float32(1))
	testJoin(float32(1), float32(0), float32(1))
	// float64
	testJoin(float64(0), float64(1), float64(1))
	testJoin(float64(1), float64(0), float64(1))

	// string
	testJoin("foo", "bar", "foo")
	testJoin("bar", "foo", "foo")
}
