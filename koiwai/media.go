package koiwai

import "github.com/uptrace/bun"

type Media struct {
	bun.BaseModel `bun:"media"`

	Hash []byte `bun:"hash,pk"`
}
