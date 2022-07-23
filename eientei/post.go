package eientei

import (
	"time"

	"github.com/uptrace/bun"
)

//Post is a post in the Eientei db.
type Post struct {
	bun.BaseModel `bun:"table:post,alias:post"`

	Board         string    `bun:"board"`
	No            int64     `bun:"no,pk"`
	Resto         int64     `bun:"resto,notnull"`
	Time          time.Time `bun:"time,notnull"`
	Name          *string   `bun:"name"`
	Trip          *string   `bun:"trip"`
	Capcode       *string   `bun:"capcode"`
	Country       *string   `bun:"country"`
	Since4Pass    *int16    `bun:"since4pass"`
	Sub           *string   `bun:"sub"`
	Com           *string   `bun:"com"`
	Tim           *int64    `bun:"tim"`
	MD5           *string   `bun:"md5"`
	Filename      *string   `bun:"filename"`
	Ext           *string   `bun:"ext"`
	Fsize         *int64    `bun:"fsize"`
	W             *int16    `bun:"w"`
	H             *int16    `bun:"h"`
	TnW           *int16    `bun:"tn_w"`
	TnH           *int16    `bun:"tn_h"`
	Deleted       bool      `bun:"deleted"`
	FileDeleted   bool      `bun:"file_deleted"`
	Spoiler       bool      `bun:"spoiler"`
	CustomSpoiler *int8     `bun:"custom_spoiler"`
	Op            bool      `bun:"op"`
	Sticky        bool      `bun:"sticky"`
	CreatedAt     time.Time `bun:"created_at,notnull"`
	LastModified  time.Time `bun:"last_modified,notnull"`
}
