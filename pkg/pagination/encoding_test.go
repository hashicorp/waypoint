package pagination

import (
	"testing"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/stretchr/testify/require"
)

func TestEncoding_DecodeToken(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		r := require.New(t)
		c, p, err := decodeToken("")
		r.NoError(err)
		r.NotNil(c)
		r.Equal(PaginatorNone, p)
	})

	t.Run("malformed", func(t *testing.T) {
		r := require.New(t)
		_, _, err := decodeToken("abclol")
		r.Error(err)
	})

	t.Run("valid", func(t *testing.T) {
		r := require.New(t)

		cursor, err := encodeCursor(&pb.PaginationCursor_Cursor{
			Value: &pb.PaginationCursor_Cursor_GormPagination{
				GormPagination: "[]",
			},
		})
		r.NoError(err)

		c, p, err := decodeToken(cursor)
		r.NoError(err)
		r.NotNil(c)
		r.EqualValues(PaginatorGormCursor, p)
	})
}

func TestEncoding_EncodeCursor(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		r := require.New(t)
		c, err := encodeCursor(nil)
		r.NoError(err)
		r.Empty(c)
	})

	t.Run("valid", func(t *testing.T) {
		r := require.New(t)
		c, err := encodeCursor(&pb.PaginationCursor_Cursor{
			Value: &pb.PaginationCursor_Cursor_GormPagination{
				GormPagination: "[a]",
			},
		})
		r.NoError(err)
		r.NotEmpty(c)

	})
}

func mockCursor() *pb.PaginationCursor_Cursor {
	return &pb.PaginationCursor_Cursor{
		Value: &pb.PaginationCursor_Cursor_GormPagination{
			GormPagination: "[a]",
		},
	}
}

func TestEncoding_CreateResponse(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		r := require.New(t)
		resp, err := createResponse(nil)
		r.NoError(err)
		r.NotNil(resp)
	})

	t.Run("next", func(t *testing.T) {
		r := require.New(t)
		resp, err := createResponse(&pb.PaginationCursor{
			Next: mockCursor(),
		})
		r.NoError(err)
		r.NotNil(resp)
		r.NotEmpty(resp.NextPageToken)
		r.Empty(resp.PreviousPageToken)
	})
	t.Run("previous", func(t *testing.T) {
		r := require.New(t)
		resp, err := createResponse(&pb.PaginationCursor{
			Previous: mockCursor(),
		})
		r.NoError(err)
		r.NotNil(resp)
		r.NotEmpty(resp.PreviousPageToken)
		r.Empty(resp.NextPageToken)
	})
	t.Run("both", func(t *testing.T) {
		r := require.New(t)
		resp, err := createResponse(&pb.PaginationCursor{
			Next:     mockCursor(),
			Previous: mockCursor(),
		})
		r.NoError(err)
		r.NotNil(resp)
		r.NotEmpty(resp.NextPageToken)
		r.NotEmpty(resp.PreviousPageToken)
	})
}
