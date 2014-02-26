// Package crdt implements utilities for working with CRDTs.
//
// Formally, a CRDT is a type for which the set of all values forms a join-semilattice.
// Less formally, a CRDT is a data type for which there exists a Join operation
// that merges information from two values in a commutative, idempotent way.
//
// This package implements Join in terms of an in-place Merge operation, such that
// Join(a, b) is equivalent to the value of result after Merge(&result, a); Merge(&result, b).
//
// Merges are done as follows:
//   * If the type implements Merger, Merge(&a, b) simply calls (&a).Merge(b).
//   * If the type is a struct, merges are done recursively fieldwise.
//   * If the type is a map, merges are done recursively keywise.
//   * If the type has a total ordering (bool, string, u?int{,8,16,32,64}, float{32,64}),
//     Merge(&a, b) sets a to the greater of (a, b).
//   * Otherwise, Merge panics.
//
// The zero value of any type is special: any non-zero value is considered to be greater than it.
// As a result, Join(a, zero) == a for any value a.
package crdt

import (
	"fmt"
	"reflect"
)

// Merger is an interface to a value that can be merged with another in place.
type Merger interface {
	// Merge joins this value with other, in place.
	// In other words, it sets this value to the least upper bound of this and other.
	// Other must be the same type as this.
	Merge(other interface{}) bool
}

// isOrdered returns true if the given kind of value has a total ordering.
func isOrdered(kind reflect.Kind) bool {
	switch kind {
	case reflect.Bool:
		return true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	case reflect.Float32, reflect.Float64:
		return true
	case reflect.String:
		return true
	default:
		return false
	}
}

// greater returns true if a > b.
func greater(a, b interface{}) bool {
	switch a.(type) {
	case bool:
		return a.(bool) && !b.(bool)
	case int:
		return a.(int) > b.(int)
	case int8:
		return a.(int8) > b.(int8)
	case int16:
		return a.(int16) > b.(int16)
	case int32:
		return a.(int32) > b.(int32)
	case int64:
		return a.(int64) > b.(int64)
	case uint:
		return a.(uint) > b.(uint)
	case uint8:
		return a.(uint8) > b.(uint8)
	case uint16:
		return a.(uint16) > b.(uint16)
	case uint32:
		return a.(uint32) > b.(uint32)
	case uint64:
		return a.(uint64) > b.(uint64)
	case float32:
		return a.(float32) > b.(float32)
	case float64:
		return a.(float64) > b.(float64)
	case string:
		return a.(string) > b.(string)
	default:
		panic("don't know how to handle type: " + reflect.TypeOf(a).String())
	}
}

// merge sets the value of a to the least upper bound of (a, b).
// It returns true if the value of a was modified.
// Both a and b must be mergeable values, and a must be addressable.
func merge(a, b reflect.Value) bool {
	var changed bool
	if merger, ok := a.Addr().Interface().(Merger); ok {
		changed = merger.Merge(b.Interface())
	} else if a.Kind() == reflect.Struct {
		for i := 0; i < a.NumField(); i++ {
			field := a.Type().Field(i)
			if field.PkgPath != "" {
				panic(fmt.Errorf("field %s (%s) is unexported", field.Name, field.PkgPath))
			}
			if merge(a.Field(i), b.Field(i)) {
				changed = true
			}
		}
	} else if a.Kind() == reflect.Map {
		if a.IsNil() && !b.IsNil() {
			a.Set(reflect.MakeMap(a.Type()))
		}
		for _, key := range b.MapKeys() {
			aValue := a.MapIndex(key)
			bValue := b.MapIndex(key)
			if aValue.IsValid() {
				newValue := reflect.New(aValue.Type()).Elem()
				merge(newValue, aValue)
				if merge(newValue, bValue) {
					a.SetMapIndex(key, newValue)
					changed = true
				}
			} else {
				a.SetMapIndex(key, bValue)
				changed = true
			}
		}
	} else if isOrdered(a.Kind()) {
		if greater(b.Interface(), a.Interface()) {
			a.Set(b)
			changed = true
		}
	} else {
		panic("don't know how to merge type " + a.Type().String())
	}
	return changed
}

// Merge sets the value of a to the least upper bound of (a, b).
// It returns true if the value of a was modified.
// a must be a pointer to a mergeable type, and b must be a non-pointer value of the same type.
func Merge(a, b interface{}) bool {
	aVal := reflect.ValueOf(a)
	bVal := reflect.ValueOf(b)
	if aVal.Kind() != reflect.Ptr {
		panic("a must be a pointer")
	}
	if aVal.Elem().Type() != bVal.Type() {
		panic("a and &b must be the same type")
	}
	return merge(aVal.Elem(), bVal)
}

func join(a, b reflect.Value) reflect.Value {
	value := reflect.New(a.Type()).Elem()
	merge(value, a)
	merge(value, b)
	return value
}

// Join returns the least upper bound of (a, b).
// Both a and b must be mergeable values of the same type.
func Join(a, b interface{}) interface{} {
	aVal := reflect.ValueOf(a)
	bVal := reflect.ValueOf(b)
	if aVal.Type() != bVal.Type() {
		panic("a and b must be the same type")
	}
	return join(aVal, bVal).Interface()
}
