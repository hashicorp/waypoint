// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package partial

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

// Eval parses code as HCL and then performs a perform evaluation on the result.
// If the result could be fully evaluated, the value is returned. Otherwise
// a string with the partial evaluated results is returned.
func Eval(ctx *hcl.EvalContext, code string) (cty.Value, string, error) {
	expr, diags := hclsyntax.ParseExpression([]byte(code), "frag.hcl", hcl.Pos{Line: 1, Column: 1})
	if diags != nil {
		return cty.Value{}, "", diags
	}

	res, ok, err := eval(ctx, expr)
	if err != nil {
		return cty.Value{}, "", err
	}

	if ok {
		return res.(*hclsyntax.LiteralValueExpr).Val, "", nil
	}

	str, err := render(ctx, res)
	if err != nil {
		return cty.Value{}, "", err
	}

	return cty.Value{}, str, nil
}

// EvalExpression takes an existing expression and partial evaluates it, returning
// a string which can be used as the value of the expression. If the expression
// is fully evaluated, the literal is converted back into a string and returned.
// If the result is still an HCL construct, it is rendered out and returned. One
// cavaet that a toplevel TemplateExpr will be rendered without surounding quotes
// when returned.
func EvalExpression(ctx *hcl.EvalContext, expr hcl.Expression) (string, error) {
	res, _, err := eval(ctx, expr)
	if err != nil {
		return "", err
	}

	str, err := renderTop(ctx, res)
	if err != nil {
		return "", err
	}

	return str, err
}

func renderTop(ctx *hcl.EvalContext, expr hcl.Expression) (string, error) {
	switch expr := expr.(type) {
	case *hclsyntax.TemplateExpr:
		return renderTemplate(ctx, expr)
	default:
		return render(ctx, expr)
	}
}

func renderTemplate(ctx *hcl.EvalContext, expr *hclsyntax.TemplateExpr) (string, error) {
	var parts []string

	for _, part := range expr.Parts {
		switch p := part.(type) {
		case *hclsyntax.LiteralValueExpr:
			parts = append(parts, p.Val.AsString())
		default:
			str, err := render(ctx, part)
			if err != nil {
				return "", err
			}

			parts = append(parts, "${"+str+"}")
		}
	}

	return strings.Join(parts, ""), nil
}

// render is called recursively, rendering out the givin expression as a string
func render(ctx *hcl.EvalContext, expr hcl.Expression) (string, error) {
	switch expr := expr.(type) {
	case *hclsyntax.LiteralValueExpr:
		// TokensForValue automatically escapes any HCL present in expr.Val
		// so we don't need to do it ourselves.
		tokens := hclwrite.TokensForValue(expr.Val)
		var buf bytes.Buffer

		_, err := tokens.WriteTo(&buf)
		if err != nil {
			return "", err
		}

		return buf.String(), nil

	case *hclsyntax.TemplateExpr:
		inner, err := renderTemplate(ctx, expr)
		if err != nil {
			return "", err
		}

		return `"` + inner + `"`, nil
	case *hclsyntax.TemplateWrapExpr:

		str, err := render(ctx, expr.Wrapped)
		if err != nil {
			return "", err
		}

		return `"${` + str + `}"`, nil
	case *hclsyntax.ScopeTraversalExpr:
		var parts []string

		for _, t := range expr.Traversal {
			switch t := t.(type) {
			case hcl.TraverseRoot:
				parts = append(parts, t.Name)
			case hcl.TraverseAttr:
				parts = append(parts, "."+t.Name)
			case hcl.TraverseIndex:
				parts = append(parts, "["+t.Key.AsString()+"]")
			default:
				panic("unknown traversal type")
			}
		}

		return strings.Join(parts, ""), nil
	case *hclsyntax.RelativeTraversalExpr:
		src, err := render(ctx, expr.Source)
		if err != nil {
			return "", err
		}

		var parts []string

		for _, t := range expr.Traversal {
			switch t := t.(type) {
			case hcl.TraverseAttr:
				parts = append(parts, "."+t.Name)
			case hcl.TraverseIndex:
				var s string

				switch t.Key.Type() {
				case cty.String:
					s = t.Key.AsString()
				case cty.Number:
					s = t.Key.AsBigFloat().String()
				default:
					return "", fmt.Errorf("unknown index type: %T", t.Key)
				}

				parts = append(parts, "["+s+"]")
			default:
				panic("unknown traversal type")
			}
		}

		return fmt.Sprintf("%s%s", src, strings.Join(parts, "")), nil
	case *hclsyntax.FunctionCallExpr:
		var parts []string

		for _, part := range expr.Args {
			str, err := render(ctx, part)
			if err != nil {
				return "", err
			}

			parts = append(parts, str)
		}

		return expr.Name + "(" + strings.Join(parts, ", ") + ")", nil
	case *hclsyntax.ConditionalExpr:
		c, err := render(ctx, expr.Condition)
		if err != nil {
			return "", err
		}

		t, err := render(ctx, expr.TrueResult)
		if err != nil {
			return "", err
		}

		f, err := render(ctx, expr.FalseResult)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("%s ? %s : %s", c, t, f), nil
	case *hclsyntax.IndexExpr:
		c, err := render(ctx, expr.Collection)
		if err != nil {
			return "", err
		}

		k, err := render(ctx, expr.Key)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("(%s)[%s]", c, k), nil
	case *hclsyntax.ParenthesesExpr:
		c, err := render(ctx, expr.Expression)
		if err != nil {
			return "", err
		}

		return "(" + c + ")", nil
	case *hclsyntax.ObjectConsKeyExpr:
		c, err := render(ctx, expr.Wrapped)
		return c, err
	case *hclsyntax.ObjectConsExpr:
		var parts []string

		for _, item := range expr.Items {
			k, err := render(ctx, item.KeyExpr)
			if err != nil {
				return "", err
			}

			v, err := render(ctx, item.ValueExpr)
			if err != nil {
				return "", err
			}

			parts = append(parts, fmt.Sprintf("%s = %s", k, v))
		}

		return "{" + strings.Join(parts, "\n") + "}", nil
	case *hclsyntax.ForExpr:

		coll, err := render(ctx, expr.CollExpr)
		if err != nil {
			return "", err
		}

		val, err := render(ctx, expr.ValExpr)
		if err != nil {
			return "", err
		}

		var key, cond string

		if expr.KeyExpr != nil {
			key, err = render(ctx, expr.KeyExpr)
			if err != nil {
				return "", err
			}
		}

		if expr.CondExpr != nil {
			cond, err = render(ctx, expr.CondExpr)
			if err != nil {
				return "", err
			}
		}

		if key != "" {
			if cond != "" {
				return fmt.Sprintf("{for %s, %s in %s: %s, %s if %s}", expr.KeyVar, expr.ValVar, coll, key, val, cond), nil
			}

			return fmt.Sprintf("{for %s, %s in %s: %s, %s}", expr.KeyVar, expr.ValVar, coll, key, val), nil
		} else {
			if cond != "" {
				return fmt.Sprintf("[for %s in %s: %s if %s]", expr.ValVar, coll, val, cond), nil
			}

			return fmt.Sprintf("[for %s in %s: %s]", expr.ValVar, coll, val), nil
		}
	case *hclsyntax.BinaryOpExpr:
		lhs, err := render(ctx, expr.LHS)
		if err != nil {
			return "", err
		}

		rhs, err := render(ctx, expr.RHS)
		if err != nil {
			return "", err
		}

		return lhs + " " + operator(expr.Op) + " " + rhs, nil
	default:
		return "", fmt.Errorf("Unknown type in render: %T", expr)
	}
}

// Weirdly there is no api or map in hcl that provides this, so we need to.
func operator(op *hclsyntax.Operation) string {
	switch op {
	case hclsyntax.OpLogicalOr:
		return string(hclsyntax.TokenOr)
	case hclsyntax.OpLogicalAnd:
		return string(hclsyntax.TokenAnd)
	case hclsyntax.OpLogicalNot:
		return string(hclsyntax.TokenBang)
	case hclsyntax.OpEqual:
		return string(hclsyntax.TokenEqualOp)
	case hclsyntax.OpNotEqual:
		return string(hclsyntax.TokenNotEqual)
	case hclsyntax.OpGreaterThan:
		return string(hclsyntax.TokenGreaterThan)
	case hclsyntax.OpGreaterThanOrEqual:
		return string(hclsyntax.TokenGreaterThanEq)
	case hclsyntax.OpLessThan:
		return string(hclsyntax.TokenLessThan)
	case hclsyntax.OpLessThanOrEqual:
		return string(hclsyntax.TokenLessThanEq)
	case hclsyntax.OpAdd:
		return string(hclsyntax.TokenPlus)
	case hclsyntax.OpSubtract:
		return string(hclsyntax.TokenMinus)
	case hclsyntax.OpMultiply:
		return string(hclsyntax.TokenStar)
	case hclsyntax.OpDivide:
		return string(hclsyntax.TokenSlash)
	case hclsyntax.OpModulo:
		return string(hclsyntax.TokenPercent)
	case hclsyntax.OpNegate:
		return string(hclsyntax.TokenMinus)
	default:
		panic("unknown operator")
	}
}

// These are some helpers because we're wrapping and unwrapping literals a lot
// in this code.

func isLit(x hclsyntax.Expression) bool {
	_, ok := x.(*hclsyntax.LiteralValueExpr)
	return ok
}

func litVal(x hclsyntax.Expression) cty.Value {
	return x.(*hclsyntax.LiteralValueExpr).Val
}

func asLit(v cty.Value) *hclsyntax.LiteralValueExpr {
	return &hclsyntax.LiteralValueExpr{Val: v}
}

// eval is called recursively to attempt to transform expr into a LiteralValueExpr.
// if not possible, it returns an expression that is as fully evaluted as possible.
func eval(ctx *hcl.EvalContext, expr hcl.Expression) (hclsyntax.Expression, bool, error) {
	switch expr := expr.(type) {
	case *hclsyntax.LiteralValueExpr:
		return expr, true, nil
	case *hclsyntax.TemplateExpr:
		top := true

		for i, part := range expr.Parts {
			val, ok, err := eval(ctx, part)
			if err != nil {
				return nil, false, err
			}

			expr.Parts[i] = val
			if !ok {
				top = false
			}
		}

		if top {
			val, diags := expr.Value(ctx)
			if diags != nil {
				return nil, true, diags
			}

			return asLit(val), true, nil
		}

		return expr, false, nil
	case *hclsyntax.TemplateWrapExpr:
		n, ok, err := eval(ctx, expr.Wrapped)
		if err != nil {
			return nil, false, err
		}

		if ok {
			return n, true, nil
		}

		expr.Wrapped = n

		return expr, false, nil
	case *hclsyntax.ScopeTraversalExpr:
		return traverseAbs(expr.Traversal, ctx)
	case *hclsyntax.RelativeTraversalExpr:
		v, ok, err := eval(ctx, expr.Source)
		if err != nil {
			return nil, false, err
		}

		expr.Source = v

		if ok {
			return traverseRel(expr.Traversal, litVal(v))
		}

		return expr, false, nil
	case *hclsyntax.FunctionCallExpr:
		top := true

		for i, part := range expr.Args {
			val, ok, err := eval(ctx, part)
			if err != nil {
				return nil, false, err
			}

			if !ok {
				top = false
			} else {
				expr.Args[i] = val
			}
		}

		if top {
			val, diags := expr.Value(ctx)
			if diags != nil {
				return nil, true, diags
			}

			return asLit(val), true, nil
		}

		return expr, false, nil
	case *hclsyntax.ConditionalExpr:
		n, ok, err := eval(ctx, expr.Condition)
		if err != nil {
			return n, ok, err
		}

		expr.Condition = n

		if ok {
			v := litVal(n)

			if v.Equals(cty.True).True() {
				return eval(ctx, expr.TrueResult)
			} else {
				return eval(ctx, expr.FalseResult)
			}
		}

		t, ok, err := eval(ctx, expr.TrueResult)
		if err != nil {
			return t, ok, err
		}

		expr.TrueResult = t

		f, ok, err := eval(ctx, expr.FalseResult)
		if err != nil {
			return f, ok, err
		}

		expr.FalseResult = f

		return expr, false, nil
	case *hclsyntax.IndexExpr:
		c, ok, err := eval(ctx, expr.Collection)
		if !ok || err != nil {
			return expr, false, err
		}

		expr.Collection = c

		idx, ok, err := eval(ctx, expr.Key)
		if !ok || err != nil {
			return expr, false, err
		}

		expr.Key = idx

		val, diags := expr.Value(ctx)
		if diags != nil {
			return nil, true, diags
		}

		return asLit(val), true, nil
	case *hclsyntax.ParenthesesExpr:
		c, ok, err := eval(ctx, expr.Expression)
		if err != nil {
			return nil, false, err
		}

		if ok {
			return c, true, nil
		}

		expr.Expression = c

		return expr, false, nil
	case *hclsyntax.ObjectConsKeyExpr:
		str := hcl.ExprAsKeyword(expr.Wrapped)
		if str != "" {
			return asLit(cty.StringVal(str)), true, nil
		}

		c, ok, err := eval(ctx, expr.Wrapped)
		if err != nil {
			return nil, false, err
		}

		expr.Wrapped = c

		return expr, ok, nil

	case *hclsyntax.ObjectConsExpr:
		try := true
		for i, item := range expr.Items {
			k, ok, err := eval(ctx, item.KeyExpr)
			if err != nil {
				return nil, false, err
			}

			if !ok {
				try = false
			}

			v, ok, err := eval(ctx, item.ValueExpr)
			if err != nil {
				return nil, false, err
			}

			if !ok {
				try = false
			}

			expr.Items[i] = hclsyntax.ObjectConsItem{KeyExpr: k, ValueExpr: v}
		}

		if try {
			val, err := expr.Value(ctx)
			if err != nil {
				return nil, false, err
			}

			return asLit(val), true, nil
		}

		return expr, false, nil
	case *hclsyntax.TupleConsExpr:
		try := true

		for i, e := range expr.Exprs {
			n, ok, err := eval(ctx, e)
			if err != nil {
				return nil, false, err
			}

			if !ok {
				try = false
			}

			expr.Exprs[i] = n
		}

		if try {
			v, err := expr.Value(ctx)
			if err != nil {
				return nil, false, err
			}

			return asLit(v), true, nil
		}

		return expr, false, nil
	case *hclsyntax.BinaryOpExpr:
		lhs, lok, err := eval(ctx, expr.LHS)
		if err != nil {
			return nil, false, err
		}

		expr.LHS = lhs

		rhs, rok, err := eval(ctx, expr.RHS)
		if err != nil {
			return nil, false, err
		}

		expr.RHS = rhs

		if lok && rok {
			v, err := expr.Value(ctx)
			if err != nil {
				return nil, false, err
			}

			return asLit(v), true, nil
		}

		return expr, false, nil
	case *hclsyntax.ForExpr:
		c, ok, err := eval(ctx, expr.CollExpr)
		if err != nil {
			return nil, false, err
		}

		try := ok

		expr.CollExpr = c

		v, _, err := eval(ctx, expr.ValExpr)
		if err != nil {
			return nil, false, err
		}

		expr.ValExpr = v

		if try {
			// TODO(evanphx) We're not doing partial eval on the ValExpr, KeyExpr or CondExpr because we'll have
			// take into account the for loop vars. That should happen in the future.
			v, err := expr.Value(ctx)
			if err != nil {
				// If there is an en error, we just return the whole for instead of actually
				// erroring out due to our slightly lossy handling of Val and Cond
				return expr, false, nil
			}

			return asLit(v), true, nil
		}

		return expr, false, nil
	default:
		return nil, false, fmt.Errorf("Unknown type in eval: %T", expr)
	}
}

// and will panic if applied to an absolute traversal.
func traverseRel(t hcl.Traversal, val cty.Value) (hclsyntax.Expression, bool, error) {
	if !t.IsRelative() {
		panic("can't use TraverseRel on an absolute traversal")
	}

	current := val
	var diags hcl.Diagnostics
	for i, tr := range t {
		var newDiags hcl.Diagnostics
		next, newDiags := tr.TraversalStep(current)
		diags = append(diags, newDiags...)
		if newDiags.HasErrors() {
			return &hclsyntax.RelativeTraversalExpr{
				Source:    asLit(current),
				Traversal: t[i:],
			}, false, nil
		}

		current = next
	}

	return asLit(current), true, nil
}

// TraverseAbs applies the receiving traversal to the given eval context,
// returning the resulting value. This is supported only for absolute
// traversals, and will panic if applied to a relative traversal.
func traverseAbs(t hcl.Traversal, ctx *hcl.EvalContext) (hclsyntax.Expression, bool, error) {
	if t.IsRelative() {
		panic("can't use TraverseAbs on a relative traversal")
	}

	split := t.SimpleSplit()
	root := split.Abs[0].(hcl.TraverseRoot)
	name := root.Name

	thisCtx := ctx
	for thisCtx != nil {
		if thisCtx.Variables == nil {
			thisCtx = thisCtx.Parent()
			continue
		}

		val, exists := thisCtx.Variables[name]
		if exists {
			expr, ok, err := traverseRel(split.Rel, val)
			if err != nil {
				return nil, false, err
			}

			if !ok {
				return &hclsyntax.ScopeTraversalExpr{Traversal: t}, false, nil
			}

			return expr, true, nil
		}
		thisCtx = thisCtx.Parent()
	}

	return &hclsyntax.ScopeTraversalExpr{Traversal: t}, false, nil
}
