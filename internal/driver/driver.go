package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// DB holds the database connection pool
type DB struct {
	SQL *sql.DB
}

var dbConnection = &DB{}

const maxOpenDbConnection = 10        // Maximum number of connections that can be open simultaneously: 10
const maxIdleDbConnection = 5         // Maximum number of idle (unused) connections allowed: 5
const maxDbLifetime = 5 * time.Minute // Maximum lifetime of a connection: 5 minutes

// ConnectSQL creates database pool for Postgres
func ConnectSQL(dsn string) (*DB, error) {
	db, err := NewDatabase(dsn)
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(maxOpenDbConnection)
	db.SetMaxIdleConns(maxIdleDbConnection)
	db.SetConnMaxLifetime(maxDbLifetime)
	dbConnection.SQL = db

	return dbConnection, nil
}

// NewDatabase creates a new database for the application
func NewDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
