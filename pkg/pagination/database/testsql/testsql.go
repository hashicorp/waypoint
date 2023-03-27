package testsql

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	migratedb "github.com/golang-migrate/migrate/v4/database"
	migrateMysql "github.com/golang-migrate/migrate/v4/database/mysql"
	migratePostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/hashicorp/go-multierror"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/mitchellh/go-testing-interface"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	gormV2 "gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Package testsql provides helpers for working with MySQL databases in
// unit tests, from creating per-test DBs to also running the migrations.
//
// This package should be used for all unit tests to safely create isolated
// database environments for your tests.
//
// When this package is inserted, it will introduce a `-gorm-debug` flag to
// the global flags package. This flag will turn on verbose output that logs
// all SQL statements.


// These variables control the database created for tests. They should not
// be modified while any active tests are running. These generally don't need
// to be modified at all.
var (
	FamilyPostgreSQL = "postgres"
	FamilyMySQL = "mysql"
	// DefaultDBName is the default name of the database to create for tests.
	DefaultDBName = "hcp_test"

	// UserName and UserPassword are the username and password, respectively,
	// of the test user to create with access to the test database.
	UserName     = "testuser"
	UserPassword = "5b5f4ba1a428498b526f799ae4ec3e59e"

	// MigrationsDir is the directory path that is looked for for migrations.
	// If this is relative, then TestDB will walk parent directories until
	// this path is found and load migrations from there.
	MigrationsDir = filepath.Join("models", "migrations")

	// mysqlDBInitialized is used to indicate whether the MySQL database was
	// created and migrated at least once. This is used in combination with the
	// ReuseDB option so that databases are only reused if they are known to
	// exist and migrated before.
	//
	// This value is stored for all invocations within this process. When
	// running tests, Go splits tests into multiple binaries (i.e. one binary
	// per package). Therefore, a database is effectively reused only across
	// tests from the same package.
	mysqlDBInitialized = false

	// postgresDBInitialized is used to indicate whether the Postgres database
	// was created and migrated at least once. This is used in combination with
	// the ReuseDB option so that databases are only reused if they are known to
	// exist and migrated before.
	//
	// This value is stored for all invocations within this process. When
	// running tests, Go splits tests into multiple binaries (i.e. one binary
	// per package). Therefore, a database is effectively reused only across
	// tests from the same package.
	postgresDBInitialized = false
)

var gormDebug = flag.Bool("gorm-debug", false, "set to true to have Gorm log all generated SQL.")

// TestDB is the legacy function to setup a test DB.
// It simply delegates to TestMySQLDB to preserve its original behavior, defined
// when cloud-sdk only supported MySQL.
//
// It's still implemented in this package in order to facilitate updating cloud-sdk
// in the HCP services.
//
// Deprecated: do not use.
//func TestDB(t testing.T, dbName string) *gorm.DB {
//	return TestMySQLDB(t, dbName)
//}

// TestDBOptions collects options that customize the test databases.
type TestDBOptions struct {
	// SkipMigration allows skipping over the migration of the database.
	SkipMigration bool

	// ReuseDB indicates whether the potentially existing test database can be
	// reused. If set, the database is created and migrated at least once, but
	// won't be destroyed and recreated every time one of the `Test<Type>DB`
	// functions is called.
	ReuseDB bool

	// DBKeepTables is used in conjunction with ReuseDB=true to
	// exclude specific relations from being truncated between tests.
	DBKeepTables []string
}

// TestPostgresDBWithOpts sets up the test DB to use, including running any migrations.
// In case the ReuseDB option is set to true, this function might not create a
// new database.
//
// This expects a local Postgres to be running with default "postgres/postgres"
// superuser credentials.
//
// This also expects that at this or some parent directory, there is the
// path represented by MigrationsDir where the migration files can be found.
// If no migrations are found, then an error is raised.
func TestPostgresDBWithOpts(t testing.T, dbName string, opts *TestDBOptions) *gorm.DB {
	t.Helper()

	// Services setting up their tests with this helper are expected to provide
	// a database name. If they don't, use a default.
	if dbName == "" {
		dbName = DefaultDBName
	}

	// Create the DB. We first drop the existing DB. The complex SQL
	// statement below evicts any connections to that database so we can
	// drop it.
	db := testDBConnectWithUser(t, "postgres", "", "postgres", "postgres")
	if *gormDebug {
		db.LogMode(true)
	}

	// If the database shouldn't be reused or if it wasn't yet initialized, drop
	// a potentially existing database, create a new one, and migrate it to the
	// latest version.
	if !opts.ReuseDB || !postgresDBInitialized {
		// Sometimes a Postgres database can't be dropped because of some internal
		// Postgres housekeeping (Postgres runs as a collection of collaborating
		// OS processes). Before trying to drop it, terminate all other connections.
		db.Exec(`SELECT pid, pg_terminate_backend(pid)
		FROM pg_stat_activity
		WHERE datname = '` + dbName + `'
		AND pid != pg_backend_pid();`)

		db.Exec("DROP DATABASE IF EXISTS " + dbName + ";")
		db.Exec("CREATE DATABASE " + dbName + ";")

		db.Exec(fmt.Sprintf(`
			DO $$
				BEGIN
					IF NOT EXISTS (
						SELECT FROM pg_catalog.pg_roles WHERE rolname = '%s'
					) THEN
						CREATE USER %s WITH PASSWORD '%s';
					END IF;
				END
			$$;`,
			UserName, UserName, UserPassword),
		)

		db.Close()

		if !opts.SkipMigration {
			// Migrate using our migrations
			testMigrate(t, FamilyPostgreSQL, dbName)
		}

		db = testDBConnectWithUser(t, FamilyPostgreSQL, dbName, "postgres", "postgres")
		db.Exec("GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO " + UserName + ";")
		db.Exec("GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO " + UserName + ";")

		db.Close()
		postgresDBInitialized = true
	} else {
		// Build a map for those tables that should be kept if reusing the database.
		tablesToKeep := make(map[string]struct{})
		for _, k := range opts.DBKeepTables {
			tablesToKeep[k] = struct{}{}
		}

		// If the database should be reused and already exists, we truncate
		// tables.
		var tablesToTruncate []string

		db = testDBConnectWithUser(t, FamilyPostgreSQL, dbName, "postgres", "postgres")

		// Find all user tables except the schema migrations table.
		rows, err := db.Table("pg_stat_user_tables").
			Where("relname != 'schema_migrations'").
			Rows()
		if err != nil {
			t.Errorf("unable to determine tables to truncate: %v", err)
		}
		defer rows.Close()

		var table struct {
			Relname    string
			Schemaname string
		}

		for rows.Next() {
			if err := db.ScanRows(rows, &table); err != nil {
				t.Errorf("unable to scan rows: %v", err)
			}
			if _, ok := tablesToKeep[table.Relname]; ok {
				continue
			}
			// We truncate the user tables from all schemas, so prepend the
			// table name with the schema name and quote both, e.g.
			//		"public"."operations"
			tablesToTruncate = append(tablesToTruncate,
				fmt.Sprintf("%q.%q", table.Schemaname, table.Relname),
			)
		}

		if len(tablesToTruncate) > 0 {
			// By truncating all tables within the same query, foreign key
			// constraints don't break.
			db = db.Exec(fmt.Sprintf("TRUNCATE %s;", strings.Join(tablesToTruncate, ",")))
			if errs := db.GetErrors(); len(errs) > 0 {
				t.Errorf("failed to truncate tables: %v",
					multierror.Append(nil, errs...),
				)
			}
		}
		db.Close()
	}

	return TestDBConnect(t, FamilyPostgreSQL, dbName)
}

//TestPostgresDBWithOptsGormV2 sets up the test DB and returns a GormV2 connection.
func TestPostgresDBWithOptsGormV2(t testing.T, dbName string, opts *TestDBOptions) *gormV2.DB {
	t.Helper()
	// wrapping TestPostgresDBWithOpts to reuse the logic to set up the test DB using Gorm V1.
	_ = TestPostgresDBWithOpts(t, dbName, opts)

	return TestDBConnectGormV2(t, FamilyPostgreSQL, dbName)
}

//TestPostgresDB sets up the test DB to use, including running any migrations.
//
//This expects a local Postgres to be running with default "postgres/postgres"
//superuser credentials.
//
//This also expects that at this or some parent directory, there is the
//path represented by MigrationsDir where the migration files can be found.
//If no migrations are found, then an error is raised.
func TestPostgresDB(t testing.T, dbName string) *gorm.DB {
	return TestPostgresDBWithOpts(t, dbName, &TestDBOptions{})
}

// TestDBConnect connects to the local test database but does not recreate it.
func TestDBConnect(t testing.T, family, dbName string) *gorm.DB {
	return testDBConnectWithUser(t, family, dbName, UserName, UserPassword)
}

// TestDBConnect connects to the local test database but does not recreate it.
func TestDBConnectGormV2(t testing.T, family, dbName string) *gormV2.DB {
	return testDBConnectWithUserGormV2(t, family, dbName, UserName, UserPassword)
}

// TestDBConnectSuper connects to the local database as a super user but does
// not recreate it.
func TestDBConnectSuper(t testing.T, family, dbName string) *gorm.DB {
	return testDBConnectWithUser(t, family, dbName, "root", "root")
}

func testDBConnectWithUser(t testing.T, family, database, user, pass string) *gorm.DB {
	t.Helper()

	var db *gorm.DB
	var err error

	if family == FamilyMySQL {
		userPass := user
		if pass != "" {
			userPass += ":" + pass
		}

		db, err = gorm.Open(FamilyMySQL,
			fmt.Sprintf("%s@tcp(127.0.0.1:3306)/%s?parseTime=true&multiStatements=true", userPass, database))
	} else {
		db, err = gorm.Open(FamilyPostgreSQL,
			fmt.Sprintf("host=127.0.0.1 port=5432 sslmode=disable user=%s password=%s dbname=%s", user, pass, database))
	}

	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// BlockGlobalUpdate - Error on update/delete without where clause
	// since it's usually a bug to not have any filters when querying models.
	//
	// By default, GORM will not include zero values in where clauses.
	// This setting may help prevent bugs or missing validation causing major data corruption.
	//
	// For example:
	//
	//	var result models.Deployment
	//	db.Model(models.Deployment{}).
	//		Where(models.Deployment{ClusterID: sqluuid.UUID{}}).
	//		First(&result)
	//
	//	Results in this query:
	//	SELECT * FROM `consul_deployments` ORDER BY `consul_deployments`.`number` ASC LIMIT 1
	//
	//  Which effectively picks a random deployment.
	db.BlockGlobalUpdate(true)

	if *gormDebug {
		db.LogMode(true)
	}

	return db
}

func testDBConnectWithUserGormV2(t testing.T, family, database, user, pass string) *gormV2.DB {
	t.Helper()

	var db *gormV2.DB
	var err error

	if family == FamilyMySQL {
		userPass := user
		if pass != "" {
			userPass += ":" + pass
		}

		db, err = gormV2.Open(mysql.New(mysql.Config{
			DriverName: FamilyMySQL,
			DSN:        fmt.Sprintf("%s@tcp(localhost:9910)/%s?charset=utf8&parseTime=True&loc=Local", userPass, database),
		}), &gormV2.Config{})
	} else {
		db, err = gormV2.Open(postgres.New(postgres.Config{
			DriverName: FamilyPostgreSQL,
			DSN:        fmt.Sprintf("host=127.0.0.1 port=5432 sslmode=disable user=%s password=%s dbname=%s", user, pass, database),
		}))
	}

	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// BlockGlobalUpdate is true by default in V2 - https://gorm.io/docs/v2_release_note.html#BlockGlobalUpdate

	if *gormDebug {
		db.Logger.LogMode(logger.Info)
	}

	return db
}

// testMigrate migrates the current database.
func testMigrate(t testing.T, family, dbName string) {
	t.Helper()

	// Find the path to the migrations. We do this using a heuristic
	// of just searching up directories until we find
	// "models/migrations". This assumes any tests run
	// will be a child of the root folder.
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("err getting working dir: %s", err)
	}
	for {
		current := filepath.Join(dir, MigrationsDir)
		_, err := os.Stat(current)
		if err == nil {
			// Found it!
			dir = current
			break
		}
		if err != nil && !os.IsNotExist(err) {
			t.Fatalf("error at %s: %s", dir, err)
		}

		// Traverse to parent
		next := filepath.Dir(dir)
		if dir == next {
			t.Fatal("cannot use DB helpers outside of folder with models/migrations")
		}
		dir = next
	}

	var driver migratedb.Driver
	if family == FamilyMySQL {
		// Connect as super user (enabling extensions) and wrap the existing DB
		// connection
		db := testDBConnectWithUser(t, FamilyMySQL, dbName, "root", "root")
		defer db.Close()
		driver, err = migrateMysql.WithInstance(db.DB(), &migrateMysql.Config{})
	} else {
		db := testDBConnectWithUser(t, FamilyPostgreSQL, dbName, "postgres", "postgres")
		defer db.Close()
		driver, err = migratePostgres.WithInstance(db.DB(), &migratePostgres.Config{})
	}
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Creator the migrator
	migrator, err := migrate.NewWithDatabaseInstance(
		"file://"+dir, family, driver)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer migrator.Close()

	// Enable logging
	if *gormDebug {
		migrator.Log = &migrateLogger{t: t}
	}

	// Migrate
	if err := migrator.Up(); err != nil {
		t.Fatalf("err migrating: %s", err)
	}
}

// migrateLogger implements migrate.Logger so that we can have logging
// on migrations when requested.
type migrateLogger struct{ t testing.T }

func (m *migrateLogger) Printf(format string, v ...interface{}) {
	m.t.Logf(format, v...)
}

func (m *migrateLogger) Verbose() bool {
	return true
}

// TestAssertCount is a helper for asserting the expected number of rows exist
// in the DB. It requires that the db argument is passed a *gorm.DB that must
// already have had a Table selection and optionally a where clause added to
// specify what to count. This helper will run `Count()` on the db passed and
// assert it succeeds and finds the desired number of records.
// Examples:
//
//	// Assert foo is empty
//	models.TestAssertCount(t, db.Table("foo"), 0)
//	// Assert 3 providers exist for a given module
//	models.TestAssertCount(t,
//		db.Model(&models.ModuleProvider{}).Where("provider = ?", provider),
//		3)
func TestAssertCount(t testing.T, db *gorm.DB, want int) {
	t.Helper()

	count := 0
	// Assume DB already describes a query that selects the rows required
	db.Count(&count)
	if errs := db.GetErrors(); len(errs) > 0 {
		err := multierror.Append(nil, errs...)
		t.Fatalf("failed counting rows: %s", err)
	}

	if want != count {
		t.Fatalf("got %d rows, want %d", count, want)
	}
}