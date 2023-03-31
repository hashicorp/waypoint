package pagination

import (
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

// Paginate paginates data
//func (p *Paginator) Paginate(db *gorm.DB, dest interface{}) (result *gorm.DB, c Cursor, err error) {
//	if err = p.validate(dest); err != nil {
//		return
//	}
//	p.setup(db, dest)
//	fields, err := p.decodeCursor(dest)
//	if err != nil {
//		return
//	}
//	if result = p.appendPagingQuery(db, fields).Find(dest); result.Error != nil {
//		return
//	}
//	// dest must be a pointer type or gorm will panic above
//	elems := reflect.ValueOf(dest).Elem()
//	// only encode next cursor when elems is not empty slice
//	if elems.Kind() == reflect.Slice && elems.Len() > 0 {
//		hasMore := elems.Len() > p.limit
//		if hasMore {
//			elems.Set(elems.Slice(0, elems.Len()-1))
//		}
//		if p.isBackward() {
//			elems.Set(reverse(elems))
//		}
//		if c, err = p.encodeCursor(elems, hasMore); err != nil {
//			return
//		}
//	}
//	return
//}