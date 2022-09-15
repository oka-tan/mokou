package koiwai

import "github.com/uptrace/bun"

//Oekaki is a tegaki replay hash on the
//database
type Oekaki struct {
	bun.BaseModel `bun:"table:oekaki"`

	Hash []byte `bun:"hash,pk"`
}
