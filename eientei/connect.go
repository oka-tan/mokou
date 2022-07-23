package eientei

import (
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

//Connect creates a connection to postgres given
//a connection string.
func Connect(s string) *bun.DB {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(s)))

	db := bun.NewDB(sqldb, pgdialect.New())

	return db
}
