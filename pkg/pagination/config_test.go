package pagination

import (
	"testing"

	validation "github.com/go-ozzo/ozzo-validation"
	publicpb "github.com/hashicorp/cloud-api-grpc-go/hashicorp/cloud/common"
	pb "github.com/hashicorp/cloud-sdk/api/pagination/proto/go"
	"github.com/stretchr/testify/require"
)

func TestPaginationConfig_Validation(t *testing.T) {
	cases := []struct {
		Name       string
		Config     *Config
		Fail       bool
		ErrMessage string
	}{
		{
			Name:       "nil",
			Config:     nil,
			Fail:       true,
			ErrMessage: "not be nil",
		},
		{
			Name:       "empty",
			Config:     &Config{},
			Fail:       true,
			ErrMessage: "SortFields: either SortFields or DefaultSortedFields is required",
		},
		{
			Name: "PageSizeDefault",
			Config: &Config{
				DefaultSortedFields: []SortedField{{Field: "test"}},
				PageSizeDefault:     0,
				NewPaginator:        PaginatorGormCursor,
				SupportedPaginators: []PaginatorType{PaginatorGormCursor},
			},
			Fail:       true,
			ErrMessage: "cannot be blank",
		},
		{
			Name: "SupportedPaginators missing NewPaginator",
			Config: &Config{
				DefaultSortedFields: []SortedField{{Field: "test"}},
				PageSizeDefault:     10,
				NewPaginator:        PaginatorGormCursor,
				SupportedPaginators: []PaginatorType{PaginatorNone},
			},
			Fail:       true,
			ErrMessage: "must be in list of supported paginators",
		},
		{
			Name: "valid",
			Config: &Config{
				DefaultSortedFields: []SortedField{{Field: "test"}},
				PageSizeDefault:     5,
				NewPaginator:        PaginatorGormCursor,
				SupportedPaginators: []PaginatorType{PaginatorGormCursor},
			},
			Fail: false,
		},
		{
			Name: "valid sorting defaults",
			Config: &Config{
				DefaultSortedFields: []SortedField{{Field: "test"}},
				PageSizeDefault:     5,
				NewPaginator:        PaginatorGormCursor,
				SupportedPaginators: []PaginatorType{PaginatorGormCursor},
			},
			Fail: false,
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			require := require.New(t)

			err := c.Config.Validate()
			if c.Fail {
				require.Error(err)
				require.Contains(err.Error(), c.ErrMessage)
			} else {
				require.NoError(err)
			}
		})
	}
}

func TestPaginationConfig_RequestValidation(t *testing.T) {
	cases := []struct {
		Name       string
		Config     *Config
		Req        *publicpb.PaginationRequest
		Fail       bool
		ErrMessage string
	}{
		{
			Name:   "empty",
			Config: &Config{PageSizeLimit: 10},
			Req:    &publicpb.PaginationRequest{},
			Fail:   false,
		},
		{
			Name:   "out of bounds",
			Config: &Config{PageSizeLimit: 10},
			Req: &publicpb.PaginationRequest{
				PageSize: 20,
			},
			Fail:       true,
			ErrMessage: "must be no greater than 10",
		},
		{
			Name:   "not base64",
			Config: &Config{PageSizeLimit: 10},
			Req: &publicpb.PaginationRequest{
				PageSize:      2,
				NextPageToken: "lol%",
			},
			Fail:       true,
			ErrMessage: "must be encoded in Base64",
		},
		{
			Name:   "one of page token",
			Config: &Config{PageSizeLimit: 10},
			Req: &publicpb.PaginationRequest{
				PageSize:          4,
				NextPageToken:     "abc",
				PreviousPageToken: "abc",
			},
			Fail:       true,
			ErrMessage: "only one of next and previous",
		},
		{
			Name:   "valid",
			Config: &Config{PageSizeLimit: 10},
			Req: &publicpb.PaginationRequest{
				PageSize:      2,
				NextPageToken: "aGVsbG8K",
			},
			Fail: false,
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			require := require.New(t)
			err := validation.ValidateStruct(c.Req, c.Config.RequestValidation(c.Req)...)
			if c.Fail {
				require.Error(err)
				require.Contains(err.Error(), c.ErrMessage)
			} else {
				require.NoError(err)
			}
		})
	}
}

func mockConfig() *Config {
	return &Config{
		DefaultSortedFields: []SortedField{
			{Field: "CreatedAt"},
			{Field: "UUID"},
		},
		PageSizeDefault:     20,
		PageSizeLimit:       100,
		NewPaginator:        PaginatorGormCursor,
		SupportedPaginators: []PaginatorType{PaginatorGormCursor},
	}
}

type request struct {
	*publicpb.PaginationRequest
}

func newRequest(pageSize uint32, next, previous string) *request {
	return &request{
		PaginationRequest: &publicpb.PaginationRequest{
			PageSize:          pageSize,
			NextPageToken:     next,
			PreviousPageToken: previous,
		},
	}
}

func (r *request) GetPagination() *publicpb.PaginationRequest {
	return r.PaginationRequest
}

func TestPaginationConfig_FromRequest(t *testing.T) {
	t.Run("empty request", func(t *testing.T) {
		r := require.New(t)
		c := mockConfig()
		req := newRequest(0, "", "")

		rc, err := c.FromRequest(req)
		r.NoError(err)
		r.Equal(c.NewPaginator, rc.Paginator)
		r.Nil(rc.Cursor.Next)
		r.Nil(rc.Cursor.Previous)
		r.Equal(c.PageSizeDefault, rc.Limit)
	})

	t.Run("limit request", func(t *testing.T) {
		r := require.New(t)
		c := mockConfig()
		req := newRequest(2, "", "")

		rc, err := c.FromRequest(req)
		r.NoError(err)
		r.EqualValues(2, rc.Limit)
	})

	t.Run("both next and previous cursor", func(t *testing.T) {
		r := require.New(t)
		c := mockConfig()
		req := newRequest(0, "aGVsbG8K", "aGVsbG8K")

		_, err := c.FromRequest(req)
		r.Error(err)
		r.Contains(err.Error(), "both previous and next page token cannot be set")
	})

	t.Run("invalid next cursor", func(t *testing.T) {
		r := require.New(t)
		c := mockConfig()
		req := newRequest(0, "abclol", "")

		_, err := c.FromRequest(req)
		r.Error(err)
		r.Contains(err.Error(), "failed to decode next page token")
	})
	t.Run("invalid previous cursor", func(t *testing.T) {
		r := require.New(t)
		c := mockConfig()
		req := newRequest(0, "", "abclol")

		_, err := c.FromRequest(req)
		r.Error(err)
		r.Contains(err.Error(), "failed to decode previous page token")
	})
	t.Run("unsupported paginator", func(t *testing.T) {
		r := require.New(t)
		c := mockConfig()
		c.NewPaginator = PaginatorNone
		c.SupportedPaginators = nil

		cursor, err := encodeCursor(&pb.PaginationCursor_Cursor{
			Value: &pb.PaginationCursor_Cursor_GormPagination{
				GormPagination: "[]",
			},
		})
		r.NoError(err)

		req := newRequest(0, cursor, "")
		_, err = c.FromRequest(req)
		r.Error(err)
		r.Contains(err.Error(), "cursor is from an unsupported paginator")
	})
}

func TestPaginationConfig_GormPaginator(t *testing.T) {
	r := require.New(t)
	c := mockConfig()
	req := newRequest(0, "", "")
	_, err := c.GormPaginator(req)
	r.NoError(err)
}

func TestPaginationConfig_GormV2Paginator(t *testing.T) {
	r := require.New(t)
	c := mockConfig()
	req := newRequest(0, "", "")
	_, err := c.GormV2Paginator(req)
	r.NoError(err)
}
