package pagination

import (
	"fmt"
	"strings"

	"github.com/agext/levenshtein"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// OrderByValidationRule is a custom validation.Rule to validate
// the SortingRequest.OrderBy field
type OrderByValidationRule struct {
	config *SortingConfig
}

// Validate validates if the SortingRequest.OrderBy contains allowed fields
// and valid per field order (asc or desc).
func (v *OrderByValidationRule) Validate(value interface{}) error {
	orderBy := value.([]string)

	for _, reqF := range orderBy {
		// One parameter can defined the whole list of fields
		// so we split the list here, if any
		reqFields := strings.Split(reqF, ",")

		for _, reqField := range reqFields {
			reqField = strings.TrimSpace(reqField)
			orderedField := strings.Split(reqField, " ")
			field := orderedField[0]

			if _, ok := v.config.AllowedSortFields[field]; !ok {
				var fields []string
				for f := range v.config.AllowedSortFields {
					fields = append(fields, f)
				}

				if suggestion := FieldSuggestion(field, fields); suggestion != "" {
					return status.Errorf(codes.InvalidArgument,
						fmt.Sprintf("field '%s' is invalid. did you mean '%s'?", field, suggestion))
				}

				return status.Errorf(codes.InvalidArgument,
					fmt.Sprintf("field '%s' is invalid. valid fields are: %s", field, strings.Join(fields, ", ")))
			}

			// In case ordering is present we check it
			if len(orderedField) > 1 {
				order := strings.ToLower(orderedField[1])
				switch order {
				case ASC, DESC:
					continue
				default:
					return status.Errorf(codes.InvalidArgument,
						fmt.Sprintf("field '%s' has invalid order. valid orders: 'asc' or 'desc'", field))
				}
			}
		}
	}

	return nil
}

// FieldSuggestion tries to find a given allowed field that is close
// to the request field and returns it if found.
// If no suggestion is close enough, returns the empty string.
//
// The suggestions are tried in order, so earlier suggestions take precedence if
// the given string is similar to two or more suggestions.
//
// This function is intended to be used with a relatively-small number of
// suggestions. It's not optimized for hundreds or thousands of them.
func FieldSuggestion(given string, fields []string) string {
	closestField := struct {
		field string
		dist  int
	}{
		dist: 3,
	}
	for _, field := range fields {
		dist := levenshtein.Distance(given, field, nil)
		if dist < closestField.dist {
			closestField.field = field
			closestField.dist = dist
		}
	}

	return closestField.field
}

// SortFieldsValidationRule is a custom validation.Rule to validate
// the if either config.SortFields or config.DefaultSortedFields is set.
type SortFieldsValidationRule struct {
	config *Config
}

func (s *SortFieldsValidationRule) Validate(value interface{}) error {
	sortFields := value.([]string)

	if len(sortFields) == 0 && len(s.config.DefaultSortedFields) == 0 {
		return fmt.Errorf("either SortFields or DefaultSortedFields is required")
	}
	return nil
}
