package convertago

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestUnwrapValue(t *testing.T) {
	type TestCase struct {
		Input    any
		IsStruct bool
	}

	testCases := map[string]TestCase{
		"bool": {
			Input:    true,
			IsStruct: false,
		},
		"int": {
			Input:    32,
			IsStruct: false,
		},
		"uint": {
			Input:    uint(32),
			IsStruct: false,
		},
		"uintptr": {
			Input:    uintptr(32),
			IsStruct: false,
		},
		"float": {
			Input:    32.0,
			IsStruct: false,
		},
		"complex": {
			Input:    1.2 + 3.4i,
			IsStruct: false,
		},
		"array": {
			Input:    [3]int{1, 2, 3},
			IsStruct: false,
		},
		"channel": {
			Input:    make(chan int),
			IsStruct: false,
		},
		"function": {
			Input:    func() {},
			IsStruct: false,
		},
		"interface": {
			Input:    new(error),
			IsStruct: false,
		},
		"map": {
			Input:    map[int]int{1: 1, 2: 2, 3: 3},
			IsStruct: false,
		},
		"slice": {
			Input:    []int{1, 2, 3},
			IsStruct: false,
		},
		"string": {
			Input:    "hello world",
			IsStruct: false,
		},
		"struct": {
			Input:    struct{}{},
			IsStruct: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			value := unwrapValue(tc.Input)
			Type := value.Type()

			if tc.IsStruct {
				assert.Equal(t, reflect.Struct, Type.Kind())
			} else {
				assert.NotEqual(t, reflect.Struct, Type.Kind())
			}
		})

		t.Run("pointer of "+name, func(t *testing.T) {
			value := unwrapValue(&tc.Input)
			Type := value.Type()

			if tc.IsStruct {
				assert.Equal(t, reflect.Struct, Type.Kind())
			} else {
				assert.NotEqual(t, reflect.Struct, Type.Kind())
			}
		})
	}
}
