package mapper

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReflectTypeMatch(t *testing.T) {
	t.Run("direct match", func(t *testing.T) {
		rt := &ReflectType{Type: reflect.TypeOf("string!")}

		value := "hello"
		require.Equal(t, value, rt.Match(value))
	})

	t.Run("direct mismatch", func(t *testing.T) {
		rt := &ReflectType{Type: reflect.TypeOf("string!")}

		value := 42
		require.Nil(t, rt.Match(value))
	})

	t.Run("interface implementation", func(t *testing.T) {
		rt := &ReflectType{Type: reflect.TypeOf((*fmt.Stringer)(nil)).Elem()}

		{
			// Doesn't implement stringer
			value := struct{}{}
			require.Nil(t, rt.Match(value))
		}

		{
			// Implements stringer
			value := &testStringer{}
			require.Equal(t, value, rt.Match(value))
		}
	})
}

type testStringer struct{}

func (v *testStringer) String() string { return "hello" }
