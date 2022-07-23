package importers

import "github.com/uptrace/bun"

//Service imports data from schema to schema.
type Service struct {
	AsagiDb   *bun.DB
	EienteiDb *bun.DB
	BatchSize int
}
