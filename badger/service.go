package badger

import (
	"github.com/uptrace/bun"
)

//Service imports data from schema to schema.
type Service struct {
	JsonFolder string
	Pg         *bun.DB
}
