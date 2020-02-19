package mapper

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChain_simple(t *testing.T) {
	type intA int
	type intB int
	type intC int
	adder := func(a intA, b intB) int { return int(a) + int(b) }
	adderFunc, err := NewFunc(adder)
	require.NoError(t, err)

	mustFunc := func(f *Func, err error) *Func {
		require.NoError(t, err)
		return f
	}

	produceA := mustFunc(NewFunc(func() intA { return intA(12) }))
	produceB := mustFunc(NewFunc(func() intB { return intB(10) }))
	produceAfromC := mustFunc(NewFunc(func(c intC) intA { return intA(c) }))
	produceBfromC := mustFunc(NewFunc(func(c intC) intB { return intB(c) * 2 }))

	var produceClock sync.Mutex
	produceConce_called := false
	produceConce := mustFunc(NewFunc(func() intC {
		produceClock.Lock()
		defer produceClock.Unlock()

		if produceConce_called {
			panic("fail")
		}
		produceConce_called = true
		return intC(5)
	}))

	noop := func() error { return nil }
	noopFunc, err := NewFunc(noop)
	require.NoError(t, err)

	t.Run("satisfied", func(t *testing.T) {
		require := require.New(t)

		chain, err := adderFunc.Chain([]*Func{}, intA(1), intB(12))
		require.NoError(err)
		result, err := chain.Call()
		require.NoError(err)
		require.Equal(result, 13)
	})

	t.Run("unsatisfied directly", func(t *testing.T) {
		require := require.New(t)

		chain, err := adderFunc.Chain([]*Func{}, intA(1))
		require.Error(err)
		require.Nil(chain)
	})

	t.Run("one func", func(t *testing.T) {
		require := require.New(t)

		chain, err := adderFunc.Chain([]*Func{
			produceA,
		}, intB(10))
		require.NoError(err)
		result, err := chain.Call()
		require.NoError(err)
		require.Equal(22, result)
	})

	t.Run("two funcs", func(t *testing.T) {
		require := require.New(t)

		chain, err := adderFunc.Chain([]*Func{
			produceA, produceB,
		})
		require.NoError(err)
		result, err := chain.Call()
		require.NoError(err)
		require.Equal(22, result)
	})

	t.Run("two funcs with input", func(t *testing.T) {
		require := require.New(t)

		chain, err := adderFunc.Chain([]*Func{
			produceAfromC, produceBfromC,
		}, intC(5))
		require.NoError(err)
		result, err := chain.Call()
		require.NoError(err)
		require.Equal(15, result)
	})

	t.Run("two levels", func(t *testing.T) {
		require := require.New(t)

		chain, err := adderFunc.Chain([]*Func{
			produceAfromC, produceBfromC, produceConce,
		})
		require.NoError(err)
		result, err := chain.Call()
		require.NoError(err)
		require.Equal(15, result)
	})

	t.Run("unsatisfied indirectly", func(t *testing.T) {
		require := require.New(t)

		chain, err := adderFunc.Chain([]*Func{
			produceAfromC, produceBfromC,
		})
		require.Error(err)
		require.Nil(chain)
		t.Log(err.Error())
	})

	t.Run("only one return value as nil", func(t *testing.T) {
		require := require.New(t)

		chain, err := noopFunc.Chain([]*Func{})
		require.NoError(err)
		_, err = chain.Call()
		require.NoError(err)
	})

}
