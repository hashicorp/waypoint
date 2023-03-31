package pagination

import (
	"github.com/hashicorp/waypoint/pkg/pagination/database/testsql"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/stretchr/testify/require"
)

type requestWithSorting struct {
	*pb.PaginationRequest
	*pb.SortingRequest
}

func (r *requestWithSorting) GetPagination() *pb.PaginationRequest {
	return r.PaginationRequest
}

func (r *requestWithSorting) GetSorting() *pb.SortingRequest {
	return r.SortingRequest
}

func newRequestWithSort(pageSize uint32, next, previous string, fields []string) *requestWithSorting {
	return &requestWithSorting{
		PaginationRequest: &pb.PaginationRequest{
			PageSize:          pageSize,
			NextPageToken:     next,
			PreviousPageToken: previous,
		},
		SortingRequest: &pb.SortingRequest{
			OrderBy: fields,
		},
	}
}

func mockSortingConfig() *SortingConfig {
	return &SortingConfig{
		AllowedSortFields: map[string]string{"created_at": "CreatedAt", "uuid": "UUID"},
		Config: Config{
			DefaultSortedFields: []SortedField{
				{Field: "CreatedAt"},
				{Field: "UUID"},
			},
			PageSizeDefault:     20,
			PageSizeLimit:       100,
			NewPaginator:        PaginatorGormCursor,
			SupportedPaginators: []PaginatorType{PaginatorGormCursor},
		},
	}
}

func TestSortingConfig_Validation(t *testing.T) {
	cases := []struct {
		Name          string
		SortingConfig *SortingConfig
		Fail          bool
		ErrMessage    string
	}{
		{
			Name:          "nil",
			SortingConfig: nil,
			Fail:          true,
			ErrMessage:    "not be nil",
		},
		{
			Name:          "empty",
			SortingConfig: &SortingConfig{},
			Fail:          true,
			ErrMessage:    "SortFields: either SortFields or DefaultSortedFields is required",
		},
		{
			Name: "PageSizeDefault",
			SortingConfig: &SortingConfig{
				AllowedSortFields: map[string]string{"test": "test"},
				Config: Config{
					DefaultSortedFields: []SortedField{{Field: "test"}},
					PageSizeDefault:     0,
					NewPaginator:        PaginatorGormCursor,
					SupportedPaginators: []PaginatorType{PaginatorGormCursor},
				},
			},
			Fail:       true,
			ErrMessage: "cannot be blank",
		},
		{
			Name: "SupportedPaginators missing NewPaginator",
			SortingConfig: &SortingConfig{
				AllowedSortFields: map[string]string{"test": "test"},
				Config: Config{
					DefaultSortedFields: []SortedField{{Field: "test"}},
					PageSizeDefault:     10,
					NewPaginator:        PaginatorGormCursor,
					SupportedPaginators: []PaginatorType{PaginatorNone},
				},
			},
			Fail:       true,
			ErrMessage: "must be in list of supported paginators",
		},
		{
			Name: "valid",
			SortingConfig: &SortingConfig{
				AllowedSortFields: map[string]string{"test": "test"},
				Config: Config{
					DefaultSortedFields: []SortedField{{Field: "test"}},
					PageSizeDefault:     5,
					NewPaginator:        PaginatorGormCursor,
					SupportedPaginators: []PaginatorType{PaginatorGormCursor},
				},
			},
			Fail: false,
		},
		{
			Name: "valid sorting defaults",
			SortingConfig: &SortingConfig{
				AllowedSortFields: map[string]string{"test": "test"},
				Config: Config{
					DefaultSortedFields: []SortedField{{Field: "test"}},
					PageSizeDefault:     5,
					NewPaginator:        PaginatorGormCursor,
					SupportedPaginators: []PaginatorType{PaginatorGormCursor},
				},
			},
			Fail: false,
		},
		{
			Name: "valid deprecated sorting defaults",
			SortingConfig: &SortingConfig{
				AllowedSortFields: map[string]string{"test": "test"},
				Config: Config{
					SortFields:          []string{"test"},
					Order:               Ascending,
					PageSizeDefault:     5,
					NewPaginator:        PaginatorGormCursor,
					SupportedPaginators: []PaginatorType{PaginatorGormCursor},
				},
			},
			Fail: false,
		},
		{
			Name: "unset AllowedSortFields",
			SortingConfig: &SortingConfig{
				Config: Config{
					SortFields:          []string{"test"},
					Order:               Ascending,
					PageSizeDefault:     5,
					NewPaginator:        PaginatorGormCursor,
					SupportedPaginators: []PaginatorType{PaginatorGormCursor},
				},
			},
			Fail:       true,
			ErrMessage: "AllowedSortFields: cannot be blank",
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			require := require.New(t)

			err := c.SortingConfig.Validate()
			if c.Fail {
				require.Error(err)
				require.Contains(err.Error(), c.ErrMessage)
			} else {
				require.NoError(err)
			}
		})
	}
}

func TestSortingConfig_RequestValidation(t *testing.T) {
	cases := []struct {
		Name          string
		SortingConfig *SortingConfig
		Req           *requestWithSorting
		Fail          bool
		ErrMessage    string
	}{
		{
			Name: "empty",
			SortingConfig: &SortingConfig{
				Config: Config{PageSizeLimit: 10}},
			Req:  &requestWithSorting{},
			Fail: false,
		},
		{
			Name: "invalid field order",
			SortingConfig: &SortingConfig{
				AllowedSortFields: map[string]string{"name": "name"},
				Config:            Config{PageSizeLimit: 10}},
			Req: &requestWithSorting{
				PaginationRequest: &pb.PaginationRequest{},
				SortingRequest: &pb.SortingRequest{
					OrderBy: []string{"name invalid"},
				},
			},
			Fail:       true,
			ErrMessage: "field 'name' has invalid order. valid orders: 'asc' or 'desc'",
		},
		{
			Name: "invalid field",
			SortingConfig: &SortingConfig{
				AllowedSortFields: map[string]string{"name": "name"},
				Config:            Config{PageSizeLimit: 10}},
			Req: &requestWithSorting{
				PaginationRequest: &pb.PaginationRequest{},
				SortingRequest: &pb.SortingRequest{
					OrderBy: []string{"title asc"},
				},
			},
			Fail:       true,
			ErrMessage: "field 'title' is invalid. valid fields are: name",
		},
		{
			Name: "sorted field with typo",
			SortingConfig: &SortingConfig{
				AllowedSortFields: map[string]string{"name": "name"},
				Config:            Config{PageSizeLimit: 10}},
			Req: &requestWithSorting{
				PaginationRequest: &pb.PaginationRequest{},
				SortingRequest: &pb.SortingRequest{
					OrderBy: []string{"nme asc"},
				},
			},
			Fail:       true,
			ErrMessage: "field 'nme' is invalid. did you mean 'name'?",
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			require := require.New(t)
			err := validation.ValidateStruct(c.Req, c.SortingConfig.RequestSortingValidation(c.Req.SortingRequest)...)
			if c.Fail {
				require.Error(err)
				require.Contains(err.Error(), c.ErrMessage)
			} else {
				require.NoError(err)
			}
		})
	}
}

func TestSortingConfig_FromRequestWithSorting(t *testing.T) {
	t.Run("empty request", func(t *testing.T) {
		r := require.New(t)
		c := mockSortingConfig()
		req := newRequestWithSort(0, "", "", []string{})

		rc, err := c.FromRequest(req)
		r.NoError(err)
		r.Equal(c.NewPaginator, rc.Paginator)
		r.Nil(rc.Cursor.Next)
		r.Nil(rc.Cursor.Previous)
		r.Equal(c.PageSizeDefault, rc.Limit)
		r.Equal(c.DefaultSortedFields, rc.SortFields)
	})

	t.Run("limit request", func(t *testing.T) {
		r := require.New(t)
		c := mockSortingConfig()
		req := newRequestWithSort(2, "", "", []string{})

		rc, err := c.FromRequest(req)
		r.NoError(err)
		r.EqualValues(2, rc.Limit)
	})
	t.Run("unsupported paginator", func(t *testing.T) {
		r := require.New(t)
		c := mockSortingConfig()
		c.NewPaginator = PaginatorNone
		c.SupportedPaginators = nil

		cursor, err := encodeCursor(&pb.PaginationCursor_Cursor{
			Value: &pb.PaginationCursor_Cursor_GormPagination{
				GormPagination: "[]",
			},
		})
		r.NoError(err)

		req := newRequestWithSort(0, cursor, "", []string{})
		_, err = c.FromRequest(req)
		r.Error(err)
		r.Contains(err.Error(), "cursor is from an unsupported paginator")
	})

	t.Run("invalid sort_field", func(t *testing.T) {
		r := require.New(t)
		c := mockSortingConfig()

		req := newRequestWithSort(0, "", "", []string{"name"})
		_, err := c.FromRequest(req)
		r.Error(err)
		r.Contains(err.Error(), "field 'name' is invalid")
	})
	t.Run("allowed sorting", func(t *testing.T) {
		r := require.New(t)
		c := mockSortingConfig()

		req := newRequestWithSort(0, "", "", []string{"created_at desc"})
		reqCtx, err := c.FromRequest(req)
		r.NoError(err)
		r.Equal([]SortedField{{Field: "CreatedAt", Order: Descending}}, reqCtx.SortFields)
	})
	t.Run("default sorting", func(t *testing.T) {
		r := require.New(t)
		c := mockSortingConfig()

		req := &requestWithSorting{
			SortingRequest: &pb.SortingRequest{},
		}
		reqCtx, err := c.FromRequest(req)
		r.NoError(err)
		r.Equal(c.DefaultSortedFields, reqCtx.SortFields)

		// Make sure still works with deprecated sorting fields
		c.SortFields = []string{"CreatedAt", "UUID"}
		c.Order = Descending
		c.DefaultSortedFields = nil

		req = &requestWithSorting{
			SortingRequest: &pb.SortingRequest{},
		}
		reqCtx, err = c.FromRequest(req)
		r.NoError(err)
		r.Equal([]SortedField{
			{Field: "CreatedAt", Order: Descending},
			{Field: "UUID", Order: Descending},
		}, reqCtx.SortFields)
	})
	t.Run("list of fields as a single parameter field", func(t *testing.T) {
		r := require.New(t)
		c := mockSortingConfig()

		sortFields := [][]string{
			{"created_at desc,uuid asc"},
			// With spaces
			{" created_at desc , uuid asc "},
		}
		for _, fields := range sortFields {
			req := newRequestWithSort(0, "", "", fields)
			reqCtx, err := c.FromRequest(req)
			r.NoError(err)
			r.Equal([]SortedField{
				{Field: "CreatedAt", Order: Descending},
				{Field: "UUID", Order: Ascending},
			}, reqCtx.SortFields)
		}
	})

	t.Run("ordering per field", func(t *testing.T) {
		r := require.New(t)
		c := mockSortingConfig()
		c.AllowedSortFields["total_cost"] = "Total"

		sortFields := []struct {
			req      []string
			expected []SortedField
		}{
			{
				req:      []string{"total_cost asc"},
				expected: []SortedField{{Field: "Total", Order: Ascending}},
			},
			{
				req:      []string{"total_cost"},
				expected: []SortedField{{Field: "Total", Order: Ascending}},
			},
			{
				req: []string{"created_at", "uuid desc"},
				expected: []SortedField{
					{Field: "CreatedAt", Order: Ascending},
					{Field: "UUID", Order: Descending},
				},
			},
		}

		for _, fields := range sortFields {
			req := newRequestWithSort(0, "", "", fields.req)
			reqCtx, err := c.FromRequest(req)
			r.NoError(err)
			r.Equal(fields.expected, reqCtx.SortFields)
		}
	})
}

func TestSortingConfig_GormPaginator(t *testing.T) {
	r := require.New(t)
	c := mockSortingConfig()
	req := newRequestWithSort(0, "", "", []string{})
	_, err := c.GormPaginator(req)
	r.NoError(err)
}

func TestSortedGormPaginator_Valid(t *testing.T) {
	db := testsql.TestPostgresDBWithOpts(t, "paginator", &testsql.TestDBOptions{
		SkipMigration: true,
	})
	db.AutoMigrate(&order{})
	// Insert some orders
	newOrders(t, db, 22)

	r := require.New(t)
	config := mockSortingConfig()
	config.AllowedSortFields = map[string]string{"created_at": "CreatedAt", "id": "ID"}
	config.DefaultSortedFields = []SortedField{
		{Field: "CreatedAt"}, // Defaults to ascending
		{Field: "UUID"},
	}
	config.PageSizeDefault = 10
	req := newRequestWithSort(9, "", "", []string{"created_at desc", "id"})

	p, err := config.GormPaginator(req)
	r.NoError(err)
	r.NotNil(p)
	r.EqualValues(PaginatorGormCursor, p.Type())

	// Paginate
	var out []*order
	result := p.Paginate(db, &out)
	r.NoError(result.Error)
	r.Len(out, 9)
	r.True(out[0].CreatedAt.After(out[1].CreatedAt))

	// Get the cursor
	c := p.Cursor()
	r.NotNil(c)
	r.NotNil(c.Next)
	r.Nil(c.Previous)

	// Get the response
	resp := p.PaginationResponse()
	r.NotNil(resp)
	r.NotEmpty(resp.NextPageToken)
	r.Empty(resp.PreviousPageToken)

	// Page again
	req = newRequestWithSort(9, resp.NextPageToken, "", []string{"created_at desc", "id"})
	p, err = config.GormPaginator(req)
	r.NoError(err)

	var out2 []*order
	result = p.Paginate(db, &out2)
	r.NoError(result.Error)
	r.Len(out2, 9)
	r.True(out2[0].CreatedAt.After(out2[1].CreatedAt))

	// Get the cursor
	c2 := p.Cursor()
	r.NotNil(c2)
	r.NotNil(c2.Next)
	r.NotNil(c2.Previous)

	// Get the response
	resp2 := p.PaginationResponse()
	r.NotNil(resp2)
	r.NotEmpty(resp2.NextPageToken)
	r.NotEmpty(resp2.PreviousPageToken)

	// Page again and clear the previous token
	req = newRequestWithSort(9, resp2.NextPageToken, "", []string{"created_at desc", "id"})
	p, err = config.GormPaginator(req)
	r.NoError(err)

	var out3 []*order
	result = p.Paginate(db, &out3)
	r.NoError(result.Error)
	r.Len(out3, 4)
	r.True(out3[0].CreatedAt.After(out3[1].CreatedAt))

	// Get the cursor
	c3 := p.Cursor()
	r.NotNil(c3)
	r.Nil(c3.Next)
	r.NotNil(c3.Previous)

	// Get the response
	resp3 := p.PaginationResponse()
	r.NotNil(resp3)
	r.Empty(resp3.NextPageToken)
	r.NotEmpty(resp3.PreviousPageToken)

	// Page back
	req = newRequestWithSort(9, "", resp3.PreviousPageToken, []string{"created_at desc", "id"})
	p, err = config.GormPaginator(req)
	r.NoError(err)

	var out4 []*order
	result = p.Paginate(db, &out4)
	r.NoError(result.Error)
	r.Len(out4, 9)
	r.True(out4[0].CreatedAt.After(out4[1].CreatedAt))

	// Get the cursor
	c4 := p.Cursor()
	r.NotNil(c4)
	r.NotNil(c4.Next)
	r.NotNil(c4.Previous)

	// Get the response
	resp4 := p.PaginationResponse()
	r.NotNil(resp4)
	r.NotEmpty(resp4.NextPageToken)
	r.NotEmpty(resp4.PreviousPageToken)
}

func TestSortingConfig_GormV2Paginator(t *testing.T) {
	r := require.New(t)
	c := mockSortingConfig()
	req := newRequestWithSort(0, "", "", []string{})
	_, err := c.GormV2Paginator(req)
	r.NoError(err)
}

func TestSortedGormV2Paginator_Valid(t *testing.T) {
	db := testsql.TestPostgresDBWithOptsGormV2(t, "paginator", &testsql.TestDBOptions{
		SkipMigration: true,
	})
	r := require.New(t)
	r.NoError(db.AutoMigrate(&order{}))
	// Insert some orders
	newOrdersGormV2(t, db, 22)

	// r := require.New(t)
	config := mockSortingConfig()
	config.AllowedSortFields = map[string]string{"created_at": "CreatedAt", "id": "ID"}
	config.DefaultSortedFields = []SortedField{
		{Field: "CreatedAt"}, // Defaults to ascending
		{Field: "UUID"},
	}
	config.PageSizeDefault = 10
	req := newRequestWithSort(9, "", "", []string{"created_at desc", "id"})

	p, err := config.GormV2Paginator(req)
	r.NoError(err)
	r.NotNil(p)
	r.EqualValues(PaginatorGormCursor, p.Type())

	// Paginate
	var out []*order
	result := p.Paginate(db, &out)
	r.NoError(result.Error)
	r.Len(out, 9)
	r.True(out[0].CreatedAt.After(out[1].CreatedAt))

	// Get the cursor
	c := p.Cursor()
	r.NotNil(c)
	r.NotNil(c.Next)
	r.Nil(c.Previous)

	// Get the response
	resp := p.PaginationResponse()
	r.NotNil(resp)
	r.NotEmpty(resp.NextPageToken)
	r.Empty(resp.PreviousPageToken)

	// Page again
	req = newRequestWithSort(9, resp.NextPageToken, "", []string{"created_at desc", "id"})
	p, err = config.GormV2Paginator(req)
	r.NoError(err)

	var out2 []*order
	result = p.Paginate(db, &out2)
	r.NoError(result.Error)
	r.Len(out2, 9)
	r.True(out2[0].CreatedAt.After(out2[1].CreatedAt))

	// Get the cursor
	c2 := p.Cursor()
	r.NotNil(c2)
	r.NotNil(c2.Next)
	r.NotNil(c2.Previous)

	// Get the response
	resp2 := p.PaginationResponse()
	r.NotNil(resp2)
	r.NotEmpty(resp2.NextPageToken)
	r.NotEmpty(resp2.PreviousPageToken)

	// Page again and clear the previous token
	req = newRequestWithSort(9, resp2.NextPageToken, "", []string{"created_at desc", "id"})
	p, err = config.GormV2Paginator(req)
	r.NoError(err)

	var out3 []*order
	result = p.Paginate(db, &out3)
	r.NoError(result.Error)
	r.Len(out3, 4)
	r.True(out3[0].CreatedAt.After(out3[1].CreatedAt))

	// Get the cursor
	c3 := p.Cursor()
	r.NotNil(c3)
	r.Nil(c3.Next)
	r.NotNil(c3.Previous)

	// Get the response
	resp3 := p.PaginationResponse()
	r.NotNil(resp3)
	r.Empty(resp3.NextPageToken)
	r.NotEmpty(resp3.PreviousPageToken)

	// Page back
	req = newRequestWithSort(9, "", resp3.PreviousPageToken, []string{"created_at desc", "id"})
	p, err = config.GormV2Paginator(req)
	r.NoError(err)

	var out4 []*order
	result = p.Paginate(db, &out4)
	r.NoError(result.Error)
	r.Len(out4, 9)
	r.True(out4[0].CreatedAt.After(out4[1].CreatedAt))

	// Get the cursor
	c4 := p.Cursor()
	r.NotNil(c4)
	r.NotNil(c4.Next)
	r.NotNil(c4.Previous)

	// Get the response
	resp4 := p.PaginationResponse()
	r.NotNil(resp4)
	r.NotEmpty(resp4.NextPageToken)
	r.NotEmpty(resp4.PreviousPageToken)
}
