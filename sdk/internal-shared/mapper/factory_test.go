package mapper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFactory(t *testing.T) {
	require := require.New(t)

	factory, err := NewFactory((*adder)(nil))
	require.NoError(err)
	require.NoError(factory.Register("two", func(a int) *adderTwo {
		return &adderTwo{From: a}
	}))

	// Get a valid mapper with satisfied types
	{
		fn := factory.Func("two")
		require.NotNil(fn)
		impl, err := fn.Call("two", 42)
		require.NoError(err)
		adder := impl.(adder)
		require.Equal(adder.Add(), 44)
	}

	// Unregistered
	{
		fn := factory.Func("three")
		require.Nil(fn)
	}
}

// Test that our function can return an interface{} type and still implement
// the factory interface.
func TestFactory_interface(t *testing.T) {
	require := require.New(t)

	factory, err := NewFactory((*adder)(nil))
	require.NoError(err)
	require.NoError(factory.Register("two", func(a int) interface{} {
		return &adderTwo{From: a}
	}))

	fn := factory.Func("two")
	require.NotNil(fn)
	impl, err := fn.Call("two", 42)
	require.NoError(err)
	adder := impl.(adder)
	require.Equal(adder.Add(), 44)
}

type adder interface {
	Add() int
}

type adderTwo struct{ From int }

func (a *adderTwo) Add() int { return a.From + 2 }
