package mapper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMapper(t *testing.T) {
	require := require.New(t)

	m := NewM((*adder)(nil))
	require.NoError(m.RegisterImpl("two", (*adderTwo)(nil)))
	require.NoError(m.RegisterMapper("two", func(a int) *adderTwo {
		return &adderTwo{From: a}
	}))

	// Get a valid mapper with satisfied types
	{
		fn := m.Mapper("two", 42)
		require.NotNil(fn)
		impl, err := fn()
		require.NoError(err)
		adder := impl.(adder)
		require.Equal(adder.Add(), 44)
	}

	// Mapper with unsatisfied types
	{
		fn := m.Mapper("two", "hello")
		require.Nil(fn)
	}
}

type adder interface {
	Add() int
}

type adderTwo struct{ From int }

func (a *adderTwo) Add() int { return a.From + 2 }
