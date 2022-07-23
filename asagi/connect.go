package asagi

import (
	"database/sql"

	//MySql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
)

//Connect creates a connection to the MySQL instance
func Connect(mysqlConnectionString string) *bun.DB {
	sqldb, err := sql.Open("mysql", mysqlConnectionString)

	if err != nil {
		panic(err)
	}

	db := bun.NewDB(sqldb, mysqldialect.New())

	return db
}
