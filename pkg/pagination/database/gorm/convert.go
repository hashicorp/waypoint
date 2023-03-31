package gorm

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
	mysqlDriver "gorm.io/driver/mysql"
	pgDriver "gorm.io/driver/postgres"
	gormv2 "gorm.io/gorm"
)

// ConvertToV2 converts gormV1 db to a gormV2 version basically "on-demand".
// Note: the gormV1 db passed in cannot have an inflight transaction.  Currently
// supports both mysql and postgres dialects.
//
// This function allows teams to incrementally/slowly migrate their code while
// still taking advantage of the cloud-sdk to setup their connections and
// continuing to support all the existing gormV1 code within their applications.
//
// The V2 migration can be daunting to take on all at once
// (https://gorm.io/docs/v2_release_note.html).  ConvertToV2 allows teams to
// just migrate one "database function" at a time within their code base and
// slowly progress to using a maintained version of the gorm package.
func ConvertToV2(v1db *gorm.DB) (v2db *gormv2.DB, err error) {
	if v1db == nil {
		return nil, errors.New("missing database")
	}
	// unfortunately, the call to gormv2.Open can/will panic at times, so we
	// need to handle that.
	defer func() {
		if deferMsg := recover(); deferMsg != nil {
			err = fmt.Errorf("%s", deferMsg)
		}
	}()
	const (
		postgresDialect = "postgres"
		mysqlDialect    = "mysql"
	)
	v1dialect := v1db.Dialect()
	if v1dialect == nil {
		return nil, errors.New("unable to determine dialect")
	}
	switch v1dialect.GetName() {
	case postgresDialect:
		return gormv2.Open(pgDriver.New(pgDriver.Config{
			// turn-off prepared statements so we can collect explain plans in
			// DD (see: https://tinyurl.com/yc47pp63)
			Conn:                 v1db.DB(),
			PreferSimpleProtocol: true,
		}), &gormv2.Config{
			PrepareStmt: false,
		})
	case mysqlDialect:
		return gormv2.Open(mysqlDriver.New(mysqlDriver.Config{
			Conn: v1db.DB(),
		}))
	default:
		return nil, fmt.Errorf("unsupported dialect %q", v1db.Dialect().GetName())
	}
}
