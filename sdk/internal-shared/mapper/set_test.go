package mapper

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetConvert_primitive(t *testing.T) {
	set := Set([]*Func{
		TestFunc(t, func(v int) string {
			return strconv.Itoa(v)
		}),
	})

	var strVal string
	require.NoError(t, set.Convert(int(12), &strVal))
	require.Equal(t, "12", strVal)
}

func TestSetConvertSlice(t *testing.T) {
	set := Set([]*Func{
		TestFunc(t, func(v int) string {
			return strconv.Itoa(v)
		}),
	})

	var out []string
	require.NoError(t, set.ConvertSlice([]int{1, 2, 3}, &out))
	require.Equal(t, []string{"1", "2", "3"}, out)
}
