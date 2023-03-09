package pagination

import (
	"errors"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/v4/is"

	publicpb "github.com/hashicorp/waypoint/pkg/server/gen"
	// publicpb "github.com/hashicorp/cloud-api-grpc-go/hashicorp/cloud/common"
	pb "github.com/hashicorp/cloud-sdk/api/pagination/proto/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Order defines the sort order of the SQL query.
type Order uint8

const (
	// unsetOrder is unset order. We define it so we can detect if the user did
	// not set a order.
	unsetOrder Order = iota

	// Ascending is sorted in ascending order.
	Ascending

	// Descending is sorted in descending order.
	Descending
)

// SortedField defines a field and its order to be use in the per field ordering rule
// at the paginator.
type SortedField struct {
	// Field name as in the model.
	Field string
	// Order is the field ordering. If left unset, it will default to Ascending.
	Order Order
}

// Config configures how a request should be paginated.
type Config struct {
	// SortFields are the fields of the model that should be used for sorting. The
	// order matters as rows are sorted in order by fields and when the field
	// matches, the next field is used to tie break the ordering.
	//
	// The fields should be immutable, unique, and orderable. If the field is
	// not unique, more than one sort fields should be passed.
	//
	// Deprecated: Prefer DefaultSortedFields
	SortFields []string

	// DefaultSortedFields are the per field ordering of the model that should be used as default for sorting.
	// If a field order is left unset, it will default to Ascending.
	DefaultSortedFields []SortedField

	// Order specifies the order by which to sort.
	//
	// Deprecated: Prefer DefaultSortedFields
	Order Order

	// PageSizeDefault is the default page size to use if the request did not
	// set a page size.
	PageSizeDefault uint32

	// PageSizeLimit is the page size limit.
	PageSizeLimit uint32

	// NewPaginator is the paginator to use if this is the first request to
	// paginate. This value must also be set in supported paginators.
	NewPaginator PaginatorType

	// SupportedPaginators are all the paginators that are supported.
	SupportedPaginators []PaginatorType
}

// Validate validates the pagination configuration is valid.
func (c *Config) Validate() error {
	if c == nil {
		return errors.New("Config must not be nil")
	}

	return validation.ValidateStruct(c,
		validation.Field(&c.SortFields, &SortFieldsValidationRule{c}),
		validation.Field(&c.PageSizeDefault, validation.Required, validation.Min(uint32(1))),
		validation.Field(&c.NewPaginator, validation.Required, validation.NotIn(PaginatorNone)),
		validation.Field(&c.SupportedPaginators, validation.Required, validation.By(func(value interface{}) error {
			paginators, _ := value.([]PaginatorType)
			for _, p := range paginators {
				if p == c.NewPaginator {
					return nil
				}
			}
			return fmt.Errorf("NewPaginator %q must be in list of supported paginators", c.NewPaginator)
		})),
	)
}

// RequestValidation returns an ozzo-validation.Rule type to validate the
// pagination request. An example usage is:
//
//	err := svcvalidation.Error(validation.ValidateStruct(req,
//		// Validate the pagination request
//		svcvalidation.StructField(&req.Pagination, func() []*validation.FieldRules {
//			return listPaginationConfig.RequestValidation(req.Pagination)
//		}),
//	))
//	if err != nil { ... }
func (c *Config) RequestValidation(req *publicpb.PaginationRequest) []*validation.FieldRules {
	oneOf := func(v interface{}) error {
		return fmt.Errorf("only one of next and previous page token may be set")
	}

	rules := []*validation.FieldRules{
		// Check that the set page limit is within bounds
		validation.Field(&req.PageSize, validation.Min(uint32(0)), validation.Max(c.PageSizeLimit)),

		// Check the tokens are base64
		validation.Field(&req.NextPageToken, is.Base64),
		validation.Field(&req.PreviousPageToken, is.Base64),
	}

	// Check only one token is set
	if req.PreviousPageToken != "" && req.NextPageToken != "" {
		rules = append(rules,
			validation.Field(&req.NextPageToken, validation.By(oneOf)),
			validation.Field(&req.PreviousPageToken, validation.By(oneOf)),
		)
	}

	return rules
}

// Request is an interface which can be used to retrieve a paginated request.
type Request interface {
	GetPagination() *publicpb.PaginationRequest
}

// RequestContext holds values related to pagination for a given request.
type RequestContext struct {
	// Config embeds the pagination config.
	// Deprecated: not used anymore
	Config Config

	// Cursor is the decoded pagination cursor from the requests page tokens.
	Cursor *pb.PaginationCursor

	// Paginator is the paginator type that was selected based on the page
	// tokens being set by a particular cursor or if it is a new pagination
	// request, by the desired NewPaginator from the config.
	Paginator PaginatorType

	// SortFields are the fields of the model that should be used for sorting. The
	// order matters as rows are sorted in order by fields and when the field
	// matches, the next field is used to tie break the ordering.
	SortFields []SortedField

	// Limit is the page size limit to be enforced.
	Limit uint32
}

// FromRequest reads the pagination values set in a request and returns an
// associated request context. An error is returned if the request holds
// malformed page tokens or has page tokens from an unsupported cursor type.
func (c *Config) FromRequest(req Request) (*RequestContext, error) {
	result := RequestContext{
		Paginator: c.NewPaginator,
		Cursor:    &pb.PaginationCursor{},
	}

	// Parse the page size
	p := req.GetPagination()
	if p.GetPageSize() == 0 {
		// Default the page size
		result.Limit = c.PageSizeDefault
	} else {
		result.Limit = p.PageSize
	}

	result.SortFields = c.DefaultSortedFields
	if len(result.SortFields) == 0 {
		// Must be using deprecated fields then
		for _, field := range c.SortFields {
			result.SortFields = append(result.SortFields, SortedField{Field: field, Order: c.Order})
		}
	}

	// Check if we have a pagination request with any page tokens set.
	if p == nil || p.NextPageToken == "" && p.PreviousPageToken == "" {
		return &result, nil
	}

	// Parse the page tokens
	if p.NextPageToken != "" && p.PreviousPageToken != "" {
		return nil, status.Error(codes.InvalidArgument, "both previous and next page token cannot be set")
	} else if p.NextPageToken != "" {
		nc, p, err := decodeToken(p.NextPageToken)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "failed to decode next page token: %v", err)
		}

		result.Paginator = p
		result.Cursor.Next = nc
	} else if p.PreviousPageToken != "" {
		// Parse the previous page cursor
		pc, p, err := decodeToken(p.PreviousPageToken)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "failed to decode previous page token: %v", err)
		}

		result.Paginator = p
		result.Cursor.Previous = pc
	}

	// Check that the paginator is one of the supported paginators
	match := false
	for _, supported := range c.SupportedPaginators {
		if supported == result.Paginator {
			match = true
			break
		}
	}
	if !match {
		return nil, status.Errorf(codes.InvalidArgument, "cursor is from an unsupported paginator: %q", result.Paginator)
	}

	return &result, nil
}

// GormPaginator is a helper that wraps FromRequest and returns a paginator that
// implements the GormPaginator interface. If the paginator selected by
// FromRequest doesn't implement GormPaginator an error is returned.
func (c *Config) GormPaginator(req Request) (GormPaginator, error) {
	rc, err := c.FromRequest(req)
	if err != nil {
		return nil, err
	}
	return c.newGormPaginator(rc)
}

func (c *Config) newGormPaginator(rc *RequestContext) (GormPaginator, error) {
	switch rc.Paginator {
	case PaginatorGormCursor:
		return NewGormCursorPaginator(rc)
	default:
		return nil, status.Errorf(codes.Internal, "selected paginator %q doesn't support Gorm based pagination", rc.Paginator)
	}
}

// GormPaginator is a helper that wraps FromRequest and returns a paginator that
// implements the GormPaginator interface. If the paginator selected by
// FromRequest doesn't implement GormPaginator an error is returned.
func (c *Config) GormV2Paginator(req Request) (GormV2Paginator, error) {
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
