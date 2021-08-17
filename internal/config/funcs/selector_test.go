package funcs

import (
	"fmt"
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func TestSelectorMatch(t *testing.T) {
	tests := []struct {
		Map map[string]string
		Selector string
		Want   cty.Value
		Err    bool
	}{
		{
			map[string]string{"env": "production"},
			"env == production",
			cty.BoolVal(true),
			false,
		},

		{
			map[string]string{"env": "production"},
			"env != production",
			cty.BoolVal(false),
			false,
		},

		/*
		This fails right now pending a  go-bexpr discussion.
		{
			map[string]string{"waypoint/workspace": "foo"},
			"waypoint.workspace == foo",
			cty.BoolVal(true),
			false,
		},
		*/
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("selectormatch(%#v, %#v)", test.Map, test.Selector), func(t *testing.T) {
			// Build our map val
			mapValues := map[string]cty.Value{}
			for k, v := range test.Map {
				mapValues[k] = cty.StringVal(v)
			}

			got, err := SelectorMatch(cty.MapVal(mapValues), cty.StringVal(test.Selector))

			if test.Err {
				if err == nil {
					t.Fatal("succeeded; want error")
				}
				return
			} else if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if !got.RawEquals(test.Want) {
				t.Errorf("wrong result\ngot:  %#v\nwant: %#v", got, test.Want)
			}
		})
	}
}
