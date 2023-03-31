package pagination

import (
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	publicpb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// These constants are used to validate the per field order in the request
const (
	ASC  = "asc"
	DESC = "desc"
)

// SortingConfig configures how a request should be paginated
// from a request with allowed sorting fields.
// If an API doesn't want to allow sorting from request, paginator.Config should be used instead.
type SortingConfig struct {
	// AllowedSortFields is a map of fields that clients are allowed to use to sort
	// the result. This is a map of api fields to their names in the Go model.
	// If not set, clients won't be able to set sorting from the request.
	// Key -> field API name
	// Value -> field model name
	//
	// Ex: {"total_cost": "Total"} - 'total_cost' is the known name for the client
	// and 'Total' is the fields name in the model.
	AllowedSortFields map[string]string

	Config
}

// Validate validates the sorted pagination configuration is valid.
func (c *SortingConfig) Validate() error {
	if c == nil {
		return errors.New("Config must not be nil")
	}

	if err := c.Config.Validate(); err != nil {
		return err
	}

	return validation.ValidateStruct(c,
		validation.Field(&c.AllowedSortFields, validation.Required),
	)
}

// RequestWithSort is an interface which can be used to retrieve a paginated and sorted request.
type RequestWithSort interface {
	GetPagination() *publicpb.PaginationRequest
	GetSorting() *publicpb.SortingRequest
}

// RequestSortingValidation returns an ozzo-validation.Rule type to validate the
// a sorting request. An example usage is:
//
//	err := svcvalidation.Error(validation.ValidateStruct(req,
//		// Validate the sorting request
//		svcvalidation.StructField(&req.Sorting, func() []*validation.FieldRules {
//			return listPaginationConfig.RequestSortingValidation(req.Sorting)
//		}),
//	))
//	if err != nil { ... }
//
// Notes:
// - The OrderByValidationRule is also called within FromRequest.
// - To validate the pagination request use listPaginationConfig.RequestValidation(*pb.PaginationRequest).
func (c *SortingConfig) RequestSortingValidation(req *publicpb.SortingRequest) []*validation.FieldRules {
	if req != nil {
		return []*validation.FieldRules{validation.Field(&req.OrderBy, &OrderByValidationRule{c})}
	}
	return nil
}

// FromRequest reads the pagination and sorting values set in a request and returns an
// associated request context. An error is returned if the request holds
// malformed page tokens or has page tokens from an unsupported cursor type.
func (c *SortingConfig) FromRequest(req RequestWithSort) (*RequestContext, error) {
	result, err := c.Config.FromRequest(req)
	if err != nil {
		return nil, err
	}

	s := req.GetSorting()

	// No sorting specified so don't override the default sorting order
	if len(s.GetOrderBy()) == 0 {
		return result, nil
	}

	// Validate order_by parameter
	rule := &OrderByValidationRule{c}
	if err := rule.Validate(req.GetSorting().OrderBy); err != nil {
		return nil, err
	}

	// Parse sorting configuration
	var fields []SortedField
	for _, reqF := range s.GetOrderBy() {
		// One parameter can defined the whole list of fields
		// so we split the list here, if any
		reqFields := strings.Split(reqF, ",")

		for _, f := range reqFields {
			sortedField := SortedField{}

			f = strings.TrimSpace(f)
			// Field might have ordering
			orderedField := strings.Split(f, " ")

			// Defaults per field oder to ascending
			sortedField.Field = c.AllowedSortFields[orderedField[0]]
			sortedField.Order = Ascending

			// In case ordering is present we attach it back to the database field name
			if len(orderedField) > 1 {
				order := strings.ToLower(orderedField[1])
				switch order {
				case ASC:
					sortedField.Order = Ascending
				case DESC:
					sortedField.Order = Descending
				}
			}

			fields = append(fields, sortedField)
		}
	}
	result.SortFields = fields

	return result, nil
}

// GormPaginator is a helper that wraps FromRequest and returns a paginator that
// implements the GormPaginator interface. If the paginator selected by
// FromRequest doesn't implement GormPaginator an error is returned.
func (c *SortingConfig) GormPaginator(req RequestWithSort) (GormPaginator, error) {
	rc, err := c.FromRequest(req)
	if err != nil {
		return nil, err
	}
	return c.newGormPaginator(rc)
}

// GormPaginator is a helper that wraps FromRequest and returns a paginator that
// implements the GormPaginator interface. If the paginator selected by
// FromRequest doesn't implement GormPaginator an error is returned.
func (c *SortingConfig) GormV2Paginator(req RequestWithSort) (GormV2Paginator, error) {
	rc, err := c.FromRequest(req)
	if err != nil {
		return nil, err
	}
	switch rc.Paginator {
	case PaginatorGormCursor:
		return NewGormV2CursorPaginator(rc)
	default:
		return nil, status.Errorf(codes.Internal, "selected paginator %q doesn't support Gorm based pagination", rc.Paginator)
	}
}