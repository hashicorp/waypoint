package pagination

import (
	"encoding/base64"
	"fmt"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	publicpb "github.com/hashicorp/waypoint/pkg/server/gen"
	// publicpb "github.com/hashicorp/cloud-api-grpc-go/hashicorp/cloud/common"
	pb "github.com/hashicorp/cloud-sdk/api/pagination/proto/go"
)

// decodeToken takes a page_token and returns the internal pagination
// cursor value and the cursor type that set it.
func decodeToken(token string) (*pb.PaginationCursor_Cursor, PaginatorType, error) {
	if token == "" {
		return &pb.PaginationCursor_Cursor{}, PaginatorNone, nil
	}

	// Base64 decode the token
	data, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return nil, PaginatorNone, err
	}

	// Unmarshal the cursor
	c := &pb.PaginationCursor_Cursor{}
	err = proto.Unmarshal(data, c)
	if err != nil {
		return nil, PaginatorNone, err
	}

	// NOTE: for now we only support gorm pagination (v1 and v2 are compatible
	// for cursor information).  More switch cases will need to be added when
	// more types of pagination are supported/added.
	switch c.GetValue().(type) {
	case *pb.PaginationCursor_Cursor_GormPagination:
		return c, PaginatorGormCursor, nil
	default:
		return nil, PaginatorNone, status.Error(codes.Internal, "token couldn't be decoded")
	}
}

// createResponse takes an internal cursor and converts it to our public
// pagination response type.
func createResponse(c *pb.PaginationCursor) (*publicpb.PaginationResponse, error) {
	resp := &publicpb.PaginationResponse{}

	// Nothing to do
	if c == nil || c.Next == nil && c.Previous == nil {
		return resp, nil
	}

	var err error
	if c.Next != nil {
		resp.NextPageToken, err = encodeCursor(c.Next)
		if err != nil {
			return nil, err
		}
	}

	if c.Previous != nil {
		resp.PreviousPageToken, err = encodeCursor(c.Previous)
		if err != nil {
			return nil, err
		}
	}

	return resp, nil
}

// encodeCursor takes a cursor and encodes it to a URL safe, base64, marshaled
// protobuf value.
func encodeCursor(c *pb.PaginationCursor_Cursor) (string, error) {
	// Nothing to do
	if c == nil {
		return "", nil
	}

	data, err := proto.Marshal(c)
	if err != nil {
		return "", fmt.Errorf("failed to encode cursor: %w", err)
	}

	// Base64 encode the token
	return base64.URLEncoding.EncodeToString(data), nil
}
