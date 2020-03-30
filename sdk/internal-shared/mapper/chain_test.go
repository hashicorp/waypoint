package mapper

import (
	"reflect"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChainTarget(t *testing.T) {
	t.Run("direct", func(t *testing.T) {
		require := require.New(t)

		chain := ChainTarget(checkMatchType(int(42)), []*Func{
			mustFunc(t, func() int { return 12 }),
			mustFunc(t, func() int32 { return 24 }),
		})
		require.NotNil(chain)

		result, err := chain.Call()
		require.NoError(err)
		require.Equal(int(12), result)
	})

	t.Run("function chain no args", func(t *testing.T) {
		require := require.New(t)

		chain := ChainTarget(checkMatchType(int(42)), []*Func{
			mustFunc(t, func(bool) int { return 12 }),
			mustFunc(t, func(string) bool { return false }),
			mustFunc(t, func() string { return "" }),
		})
		require.NotNil(chain)

		result, err := chain.Call()
		require.NoError(err)
		require.Equal(int(12), result)
	})

	t.Run("function chain with args", func(t *testing.T) {
		require := require.New(t)

		chain := ChainTarget(checkMatchType(int(42)), []*Func{
			mustFunc(t, func(bool) int { return 12 }),
			mustFunc(t, func(string) bool { return false }),
		}, "hello")
		require.NotNil(chain)

		result, err := chain.Call()
		require.NoError(err)
		require.Equal(int(12), result)
	})

	t.Run("function chain with args (unsatisfied)", func(t *testing.T) {
		require := require.New(t)

		chain := ChainTarget(checkMatchType(int(42)), []*Func{
			mustFunc(t, func(bool) int { return 12 }),
			mustFunc(t, func(string) bool { return false }),
		})
		require.Nil(chain)
	})
}

func TestFuncChain(t *testing.T) {
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

	t.Run("two funcs with input via extra args", func(t *testing.T) {
		require := require.New(t)

		adderFunc := TestFunc(t, adder, WithValues(intC(5)))
		chain, err := adderFunc.Chain([]*Func{
			produceAfromC, produceBfromC,
		})
		require.NoError(err)
		result, err := chain.Call()
		require.NoError(err)
		require.Equal(15, result)
	})
}

func TestFuncChainInputSet(t *testing.T) {
	checkOr := func(fs ...func(Type) bool) func(Type) bool {
		return func(t Type) bool {
			for _, f := range fs {
				if f(t) {
					return true
				}
			}

			return false
		}
	}

	checkMatchBool := checkMatchType(false)
	checkMatchInt := checkMatchType(int(12))

	t.Run("already satisfied (no args)", func(t *testing.T) {
		require := require.New(t)

		f, err := NewFunc(func() int { return 42 })
		require.NoError(err)
		result := f.ChainInputSet([]*Func{}, func(Type) bool { return false })
		require.NotNil(result)
		require.Len(result, 0)
	})

	t.Run("depth 0", func(t *testing.T) {
		require := require.New(t)

		f := mustFunc(t, func(v int) int { return v + 1 })
		result := f.ChainInputSet([]*Func{}, checkMatchInt)
		require.NotNil(result)

		// We expect 1 result because we want the one "int"
		require.Len(result, 1)
	})

	t.Run("depth 0,1", func(t *testing.T) {
		require := require.New(t)

		type intA int
		type intB int

		// The path here:
		//   - f requires intA, intB
		//   - a mapper provides intB from int
		//   - a mapper provides intA from int
		//   > can solve with int
		f := mustFunc(t, func(intA, intB) int { return 0 })
		result := f.ChainInputSet([]*Func{
			mustFunc(t, func(int) intA { return 0 }),
			mustFunc(t, func(int) intB { return 0 }),
		}, checkMatchInt)
		require.NotNil(result)

		// We expect 1 result because we want the one "int"
		require.Len(result, 1)
	})

	t.Run("depth 1,2", func(t *testing.T) {
		require := require.New(t)

		type intA int
		type intB int
		type intC int

		// The path here:
		//   - f requires intA, intB
		//   - a mapper provides intA from int
		//   - a mapper provides intC from inB
		//   - a mapper provides intB from int
		//   > can solve with int
		f := mustFunc(t, func(intA, intB) int { return 0 })
		result := f.ChainInputSet([]*Func{
			mustFunc(t, func(int) intA { return 0 }),
			mustFunc(t, func(intC) intB { return 0 }),
			mustFunc(t, func(int) intC { return 0 }),
		}, checkMatchInt)
		require.NotNil(result)

		// We expect 1 result because we want the one "int"
		require.Len(result, 1)
	})

	t.Run("cycle", func(t *testing.T) {
		require := require.New(t)

		type intA int
		type intB int
		type intC int

		f := mustFunc(t, func(intA, intB) int { return 0 })
		result := f.ChainInputSet([]*Func{
			mustFunc(t, func(intB) intA { return 0 }),
			mustFunc(t, func(intA) intB { return 0 }),
		}, checkMatchInt)
		require.Nil(result)
	})

	t.Run("depth 0 multiple", func(t *testing.T) {
		require := require.New(t)

		f := mustFunc(t, func(int, bool) int { return 0 })
		result := f.ChainInputSet([]*Func{}, checkOr(checkMatchInt, checkMatchBool))
		require.NotNil(result)

		// We expect 2 results: bool and int
		require.Len(result, 2)
	})

	t.Run("depth 1,2 multiple", func(t *testing.T) {
		require := require.New(t)

		type intA int
		type intB int
		type intC int

		// The path here:
		//   - f requires intA, intB
		//   - a mapper provides intA from int
		//   - a mapper provides intC from inB
		//   - a mapper provides intB from bool
		//   > can solve with int and bool
		f := mustFunc(t, func(intA, intB) int { return 0 })
		result := f.ChainInputSet([]*Func{
			mustFunc(t, func(int) intA { return 0 }),
			mustFunc(t, func(intC) intB { return 0 }),
			mustFunc(t, func(bool) intC { return 0 }),
		}, checkOr(checkMatchInt, checkMatchBool))
		require.NotNil(result)

		// We expect 2 results: bool and int
		require.Len(result, 2)
	})
}

func checkMatchType(v interface{}) func(Type) bool {
	return func(t Type) bool {
		return t.(*ReflectType).Type == reflect.TypeOf(v)
	}
}
