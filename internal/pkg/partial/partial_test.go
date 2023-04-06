// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package partial

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/require"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/function/stdlib"
)

func TestPartial(t *testing.T) {
	t.Run("static values", func(t *testing.T) {
		var ctx hcl.EvalContext
		val, _, err := Eval(&ctx, `"hello"`)
		require.NoError(t, err)

		require.Equal(t, "hello", val.AsString())
	})

	t.Run("known templates", func(t *testing.T) {
		var ctx hcl.EvalContext
		ctx.Variables = map[string]cty.Value{
			"foo": cty.MapVal(map[string]cty.Value{
				"bar": cty.StringVal("qux"),
			}),
		}

		val, _, err := Eval(&ctx, `"hello: ${foo.bar}"`)
		require.NoError(t, err)

		require.Equal(t, "hello: qux", val.AsString())
	})

	t.Run("unknown templates", func(t *testing.T) {
		var ctx hcl.EvalContext
		ctx.Variables = map[string]cty.Value{
			"foo": cty.MapVal(map[string]cty.Value{
				"nope": cty.BoolVal(false),
			}),
		}

		_, str, err := Eval(&ctx, `"hello: ${foo.bar}"`)
		require.NoError(t, err)

		require.Equal(t, `"hello: ${foo.bar}"`, str)
	})

	t.Run("mix of known and unknown templates", func(t *testing.T) {
		var ctx hcl.EvalContext
		ctx.Variables = map[string]cty.Value{
			"foo": cty.MapVal(map[string]cty.Value{
				"bar": cty.StringVal("qux"),
			}),
		}

		_, str, err := Eval(&ctx, `"hello: ${foo.bar} ${remote.name}"`)
		require.NoError(t, err)

		require.Equal(t, `"hello: qux ${remote.name}"`, str)
	})

	t.Run("static function args", func(t *testing.T) {
		var ctx hcl.EvalContext

		ctx.Functions = map[string]function.Function{
			"upper": stdlib.UpperFunc,
		}

		val, _, err := Eval(&ctx, `"hello: ${upper("blah")}"`)
		require.NoError(t, err)

		require.Equal(t, "hello: BLAH", val.AsString())
	})

	t.Run("unknown function args", func(t *testing.T) {
		var ctx hcl.EvalContext

		ctx.Functions = map[string]function.Function{
			"upper": stdlib.UpperFunc,
		}

		_, str, err := Eval(&ctx, `"hello: ${upper(local.name)}"`)
		require.NoError(t, err)

		require.Equal(t, `"hello: ${upper(local.name)}"`, str)
	})

	t.Run("known function args", func(t *testing.T) {
		var ctx hcl.EvalContext

		ctx.Variables = map[string]cty.Value{
			"foo": cty.MapVal(map[string]cty.Value{
				"bar": cty.StringVal("qux"),
			}),
		}

		ctx.Functions = map[string]function.Function{
			"upper": stdlib.UpperFunc,
		}

		val, _, err := Eval(&ctx, `"hello: ${upper(foo.bar)}"`)
		require.NoError(t, err)

		require.Equal(t, "hello: QUX", val.AsString())
	})

	t.Run("mixed function args", func(t *testing.T) {
		var ctx hcl.EvalContext

		ctx.Variables = map[string]cty.Value{
			"foo": cty.MapVal(map[string]cty.Value{
				"bar": cty.StringVal("qux"),
			}),
		}

		ctx.Functions = map[string]function.Function{
			"upper": stdlib.UpperFunc,
		}

		_, str, err := Eval(&ctx, `"hello: ${upper("a", foo.bar, local.name)}"`)
		require.NoError(t, err)

		require.Equal(t, `"hello: ${upper("a", "qux", local.name)}"`, str)
	})

	t.Run("static conditions", func(t *testing.T) {
		var ctx hcl.EvalContext

		ctx.Variables = map[string]cty.Value{
			"foo": cty.MapVal(map[string]cty.Value{
				"bar": cty.BoolVal(true),
			}),
		}

		val, _, err := Eval(&ctx, `"hello: ${true ? "yup" : "nope"}"`)
		require.NoError(t, err)

		require.Equal(t, `hello: yup`, val.AsString())
	})

	t.Run("cond: known conditions", func(t *testing.T) {
		var ctx hcl.EvalContext

		ctx.Variables = map[string]cty.Value{
			"foo": cty.MapVal(map[string]cty.Value{
				"bar": cty.BoolVal(true),
			}),
		}

		val, _, err := Eval(&ctx, `"hello: ${foo.bar ? "yup" : "nope"}"`)
		require.NoError(t, err)

		require.Equal(t, `hello: yup`, val.AsString())
	})

	t.Run("cond: unknown conditions", func(t *testing.T) {
		var ctx hcl.EvalContext

		ctx.Variables = map[string]cty.Value{
			"foo": cty.MapVal(map[string]cty.Value{
				"bar": cty.BoolVal(true),
			}),
		}

		_, str, err := Eval(&ctx, `"hello: ${local.name ? "yup" : "nope"}"`)
		require.NoError(t, err)

		require.Equal(t, `"hello: ${local.name ? "yup" : "nope"}"`, str)
	})

	t.Run("cond: unknown conditions, known branches", func(t *testing.T) {
		var ctx hcl.EvalContext

		ctx.Variables = map[string]cty.Value{
			"foo": cty.MapVal(map[string]cty.Value{
				"bar": cty.BoolVal(true),
			}),
		}

		_, str, err := Eval(&ctx, `"hello: ${local.name ? foo.bar : "nope"}"`)
		require.NoError(t, err)

		require.Equal(t, `"hello: ${local.name ? true : "nope"}"`, str)
	})

	t.Run("cond: known true, unknown false", func(t *testing.T) {
		var ctx hcl.EvalContext

		ctx.Variables = map[string]cty.Value{
			"foo": cty.MapVal(map[string]cty.Value{
				"name": cty.StringVal("ok"),
			}),
			"bar": cty.BoolVal(true),
		}

		val, _, err := Eval(&ctx, `"hello: ${bar ? foo.name : local.nope}"`)
		require.NoError(t, err)

		require.Equal(t, `hello: ok`, val.AsString())
	})

	t.Run("cond: unknown true, known false", func(t *testing.T) {
		var ctx hcl.EvalContext

		ctx.Variables = map[string]cty.Value{
			"foo": cty.MapVal(map[string]cty.Value{
				"name": cty.StringVal("ok"),
			}),

			"bar": cty.BoolVal(true),
		}

		_, str, err := Eval(&ctx, `"hello: ${bar ? local.nope : foo.name}"`)
		require.NoError(t, err)

		require.Equal(t, `"hello: ${local.nope}"`, str)
	})

	t.Run("index known", func(t *testing.T) {
		var ctx hcl.EvalContext

		ctx.Variables = map[string]cty.Value{
			"foo": cty.MapVal(map[string]cty.Value{
				"bar": cty.ListVal([]cty.Value{cty.StringVal("ok")}),
			}),

			"idx": cty.NumberIntVal(0),
		}

		val, _, err := Eval(&ctx, `"hello: ${foo.bar[idx]}"`)
		require.NoError(t, err)

		require.Equal(t, `hello: ok`, val.AsString())
	})

	t.Run("index unknown", func(t *testing.T) {
		var ctx hcl.EvalContext

		ctx.Variables = map[string]cty.Value{
			"foo": cty.MapVal(map[string]cty.Value{
				"bar": cty.ListVal([]cty.Value{cty.StringVal("ok")}),
			}),
		}

		_, str, err := Eval(&ctx, `"hello: ${foo.bar[local.name]}"`)
		require.NoError(t, err)

		require.Equal(t, `"hello: ${(["ok"])[local.name]}"`, str)
	})

	t.Run("index static", func(t *testing.T) {
		var ctx hcl.EvalContext

		ctx.Variables = map[string]cty.Value{
			"foo": cty.MapVal(map[string]cty.Value{
				"bar": cty.ListVal([]cty.Value{cty.StringVal("ok")}),
			}),
		}

		val, _, err := Eval(&ctx, `"hello: ${foo.bar[0]}"`)
		require.NoError(t, err)

		require.Equal(t, `hello: ok`, val.AsString())
	})

	t.Run("index through relative", func(t *testing.T) {
		var ctx hcl.EvalContext

		ctx.Variables = map[string]cty.Value{
			"foo": cty.MapVal(map[string]cty.Value{
				"bar": cty.ListVal([]cty.Value{cty.StringVal("ok")}),
			}),
		}

		val, _, err := Eval(&ctx, `"hello: ${({x = ["ok"]}).x[0]}"`)
		require.NoError(t, err)

		require.Equal(t, `hello: ok`, val.AsString())
	})

	t.Run("for known", func(t *testing.T) {
		var ctx hcl.EvalContext

		ctx.Variables = map[string]cty.Value{
			"foo": cty.MapVal(map[string]cty.Value{
				"bar": cty.ListVal([]cty.Value{cty.StringVal("ok")}),
			}),
		}

		ctx.Functions = map[string]function.Function{
			"upper": stdlib.UpperFunc,
		}

		val, _, err := Eval(&ctx, `"hello: ${[for v in foo.bar: upper(v)][0]}"`)
		require.NoError(t, err)

		require.Equal(t, `hello: OK`, val.AsString())
	})

	t.Run("for unknown", func(t *testing.T) {
		var ctx hcl.EvalContext

		ctx.Variables = map[string]cty.Value{
			"foo": cty.MapVal(map[string]cty.Value{
				"bar": cty.ListVal([]cty.Value{cty.StringVal("ok")}),
			}),
		}

		ctx.Functions = map[string]function.Function{
			"upper": stdlib.UpperFunc,
		}

		_, str, err := Eval(&ctx, `"hello: ${[for v in qux: upper(v)][0]}"`)
		require.NoError(t, err)

		require.Equal(t, `"hello: ${[for v in qux: upper(v)][0]}"`, str)
	})

	t.Run("for scope var in val", func(t *testing.T) {
		var ctx hcl.EvalContext

		ctx.Variables = map[string]cty.Value{
			"foo": cty.MapVal(map[string]cty.Value{
				"bar": cty.ListVal([]cty.Value{cty.StringVal("ok")}),
			}),
			"prefix": cty.StringVal("ST"),
		}

		ctx.Functions = map[string]function.Function{
			"upper": stdlib.UpperFunc,
		}

		val, _, err := Eval(&ctx, `"hello: ${[for v in foo.bar: prefix][0]}"`)
		require.NoError(t, err)

		require.Equal(t, `hello: ST`, val.AsString())
	})

	t.Run("for mixed", func(t *testing.T) {
		var ctx hcl.EvalContext

		ctx.Variables = map[string]cty.Value{
			"foo": cty.MapVal(map[string]cty.Value{
				"bar": cty.ListVal([]cty.Value{cty.StringVal("ok")}),
			}),
			"prefix": cty.StringVal("ST"),
		}

		ctx.Functions = map[string]function.Function{
			"upper": stdlib.UpperFunc,
		}

		_, str, err := Eval(&ctx, `"hello: ${[for v in list: prefix][0]}"`)
		require.NoError(t, err)

		require.Equal(t, `"hello: ${[for v in list: "ST"][0]}"`, str)
	})

	t.Run("for mixed 2", func(t *testing.T) {
		var ctx hcl.EvalContext

		ctx.Variables = map[string]cty.Value{
			"foo": cty.MapVal(map[string]cty.Value{
				"bar": cty.ListVal([]cty.Value{cty.StringVal("ok")}),
			}),
			"prefix": cty.StringVal("ST"),
		}

		ctx.Functions = map[string]function.Function{
			"upper": stdlib.UpperFunc,
		}

		_, str, err := Eval(&ctx, `"hello: ${[for v in list: prefix + v + later.known][0]}"`)
		require.NoError(t, err)

		require.Equal(t, `"hello: ${[for v in list: "ST" + v + later.known][0]}"`, str)
	})

}
