package asagi

import (
	"github.com/uptrace/bun"
)

//Post is a post in the Asagi schema.
//This includes OPs.
type Post struct {
	bun.BaseModel `bun:"table:post,alias:post"`

	//Auto-incremented pk.
	DocID            uint    `bun:"doc_id,pk"`
	Num              uint    `bun:"num,notnull"`
	Subnum           uint    `bun:"subnum,notnull"`
	ThreadNum        uint    `bun:"thread_num,notnull"`
	Op               bool    `bun:"op,notnull"`
	Timestamp        uint    `bun:"timestamp,notnull"`
	TimestampExpired uint    `bun:"timestamp_expired,notnull"`
	PreviewOrig      *string `bun:"preview_orig"`
	PreviewW         uint16  `bun:"preview_w"`
	PreviewH         uint16  `bun:"preview_h"`
	MediaFilename    *string `bun:"media_filename"`
	MediaW           uint16  `bun:"media_w"`
	MediaH           uint16  `bun:"media_h"`
	MediaSize        uint16  `bun:"media_size"`
	MediaHash        *string `bun:"media_hash"`
	MediaOrig        *string `bun:"media_orig"`
	Spoiler          bool    `bun:"spoiler"`
	Deleted          bool    `bun:"deleted"`
	Capcode          string  `bun:"capcode"`
	Email            *string `bun:"email"`
	Name             *string `bun:"name"`
	Trip             *string `bun:"trip"`
	Title            *string `bun:"title"`
	Comment          *string `bun:"comment"`
	Delpass          *string `bun:"delpass"`
	Sticky           bool    `bun:"sticky"`
	Locked           bool    `bun:"locked"`
	PosterHash       *string `bun:"poster_hash"`
	PosterCountry    *string `bun:"poster_country"`
	Exif             *string `bun:"exif"`
}
