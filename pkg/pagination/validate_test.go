package pagination

import "testing"

func TestFieldSuggestion(t *testing.T) {
	tc := []struct {
		given     string
		fields    []string
		suggested string
	}{
		{
			given:     "urname",
			fields:    []string{"name", "surname"},
			suggested: "surname",
		},
		{
			given:     "ubtitle",
			fields:    []string{"title", "subtitle"},
			suggested: "subtitle",
		},
	}

	for _, c := range tc {
		t.Run(c.given, func(t *testing.T) {
			suggested := FieldSuggestion(c.given, c.fields)
			if suggested != c.suggested {
				t.Fatalf("wrong suggestion: %s", suggested)
			}
		})
	}
}
