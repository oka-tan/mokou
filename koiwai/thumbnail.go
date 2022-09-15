package koiwai

import "github.com/uptrace/bun"

//Thumbnail is a thumbnail hash on the database
type Thumbnail struct {
	bun.BaseModel `bun:"table:thumbnail"`

	Hash []byte `bun:"hash,pk"`
}
