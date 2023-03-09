package pagination

import (
	// publicpb "github.com/hashicorp/cloud-api-grpc-go/hashicorp/cloud/common"
	pb "github.com/hashicorp/cloud-sdk/api/pagination/proto/go"
	publicpb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/jinzhu/gorm"
	gormV2 "gorm.io/gorm"
)

// Paginator is the interface for a generic paginator.
type Paginator interface {
	// Type returns the paginator type.
	Type() PaginatorType

	// Cursor returns the internal pagination cursor. Behavior is only defined
	// once a pagination has occurred.
	Cursor() *pb.PaginationCursor

	// PaginationResponse returns the public pagination response. Behavior is
	// only defined once a pagination has occurred.
	PaginationResponse() *publicpb.PaginationResponse
}

// GormPaginator is the interface for a paginator that supports operating on a
// gorm connection.
type GormPaginator interface {
	Paginator

	// Paginate is used to paginate a GORM query. The results are returned via
	// the out interface. An example usage is as follows:
	//
	//   paginator, err := paginationConfig.GormPaginator(req)
	//   if err != nil { ... }
	//
	//   var outModels []*Model
	//   query := s.DB(ctx).Where("...")
	//   result := paginator.Paginate(query, &outModels)
	//   if result.Error != nil { ... } // Or result.GetErrors() for an error list
	Paginate(stmt *gorm.DB, out interface{}) *gorm.DB
}

// PaginatorType is the type of paginator.
type PaginatorType string

const (
	PaginatorNone       PaginatorType = "none"
	PaginatorGormCursor PaginatorType = "gorm-cursor-paginator"
)

// GormPaginator is the interface for a paginator that supports operating on a
// gorm connection.
type GormV2Paginator interface {
	Paginator

	// Paginate is used to paginate a GORM query. The results are returned via
	// the out interface. An example usage is as follows:
	//
	//   paginator, err := paginationConfig.GormPaginator(req)
	//   if err != nil { ... }
	//
	//   var outModels []*Model
	//   query := s.DB(ctx).Where("...")
	//   result := paginator.Paginate(query, &outModels)
	//   if result.Error != nil { ... } // Or result.GetErrors() for an error list
	Paginate(stmt *gormV2.DB, out interface{}) *gormV2.DB
}
