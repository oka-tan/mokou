//Package koiwai exposes models for the target db
package koiwai

import "github.com/uptrace/bun"

//Media is a media hash on the database
type Media struct {
	bun.BaseModel `bun:"media"`

	Hash []byte `bun:"hash,pk"`
}
