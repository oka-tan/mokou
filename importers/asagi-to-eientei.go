package importers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"mokou/asagi"
	"mokou/config"
	"mokou/eientei"
	"mokou/utils"
	"strconv"
	"strings"
	"time"
)

//AsagiToEientei imports the board given described in boardConfig
//from Asagi to Eientei
func (s *Service) AsagiToEientei(boardConfig *config.BoardConfig) error {
	log.Printf("Importing data from Asagi to Eientei for board %s\n", boardConfig.Name)

	postRestorer := NewPostRestorer(boardConfig)

	eienteiTx, err := s.EienteiDb.BeginTx(context.Background(), &sql.TxOptions{})
	defer eienteiTx.Rollback()

	if err != nil {
		return err
	}

	asagiTx, err := s.AsagiDb.BeginTx(context.Background(), &sql.TxOptions{})
	defer asagiTx.Rollback()

	if err != nil {
		return err
	}

	asagiPosts := make([]asagi.Post, 0, s.BatchSize)
	eienteiPosts := make([]eientei.Post, 0, s.BatchSize)
	var keyset uint

	for {
		asagiPosts = asagiPosts[0:0]
		eienteiPosts = eienteiPosts[0:0]

		err := asagiTx.NewSelect().
			Model(&asagiPosts).
			ModelTableExpr(fmt.Sprintf("`%s` AS post", boardConfig.Name)).
			Where("`subnum` = 0").
			Where("`doc_id` > ?", keyset).
			Order("doc_id ASC").
			Limit(s.BatchSize).
			Scan(context.Background())

		if err != nil {
			return err
		}

		if len(asagiPosts) == 0 {
			break
		}

		keyset = asagiPosts[len(asagiPosts)-1].DocID

		for _, asagiPost := range asagiPosts {
			exif := make(map[string]string)

			if asagiPost.Exif != nil && *asagiPost.Exif != "" {
				err := json.Unmarshal([]byte(*asagiPost.Exif), &exif)

				if err != nil {
					log.Printf("Error parsing post %d, exif field: %s, error: %s", asagiPost.DocID, *asagiPost.Exif, err)
				}
			}

			resto := asagiPost.ThreadNum
			if resto == 0 {
				resto = asagiPost.Num
			}

			timestamp := time.UnixMilli(int64(asagiPost.Timestamp)).Add(4 * time.Hour)

			name := asagiPost.Name
			if name == nil || *name == "Anonymous" {
				name = nil
			} else {
				filteredName := utils.FilterHtml(*name)
				name = &filteredName
			}

			trip := asagiPost.Trip
			if trip == nil || *trip == "" {
				trip = nil
			} else {
				filteredTrip := utils.FilterHtml(*trip)
				trip = &filteredTrip
			}

			var capcode *string
			switch asagiPost.Capcode {
			case "N":
				{
					capcode = nil
				}
			case "M":
				{
					modString := "mod"
					capcode = &modString
				}
			case "A":
				{
					adminString := "admin"
					capcode = &adminString
				}
			case "D":
				{
					developerString := "developer"
					capcode = &developerString
				}
			case "F":
				{
					founderString := "founder"
					capcode = &founderString
				}
			}

			var country *string
			if asagiPost.PosterCountry != nil && len(*asagiPost.PosterCountry) == 2 {
				country = asagiPost.PosterCountry
			}

			sub := asagiPost.Title
			if sub != nil && *sub == "" {
				sub = nil
			}

			var com *string
			if asagiPost.Comment != nil && *asagiPost.Comment != "" {
				com = postRestorer.RestoreComment(&asagiPost, exif)
			}

			var tim *int64
			var ext *string
			if asagiPost.PreviewOrig != nil {
				lastIndex := strings.LastIndex(*asagiPost.PreviewOrig, ".")

				if lastIndex != -1 {
					timV, err := strconv.ParseInt((*asagiPost.PreviewOrig)[:lastIndex], 10, 64)

					if err != nil {
						tim = &timV
						extV := (*asagiPost.PreviewOrig)[lastIndex:]
						ext = &extV
					}
				}
			}

			var filename *string
			if asagiPost.MediaFilename != nil {
				lastIndex := strings.LastIndex(*asagiPost.MediaFilename, ".")

				if lastIndex != -1 {
					filenameV := (*asagiPost.MediaFilename)[:lastIndex]
					filename = &filenameV
				}

			}

			var fsize *int64
			if asagiPost.MediaSize > 0 {
				fsizeV := int64(asagiPost.MediaSize)
				fsize = &fsizeV
			}

			var w *int16
			if asagiPost.MediaW > 0 {
				wV := int16(asagiPost.MediaW)
				w = &wV
			}

			var h *int16
			if asagiPost.MediaH > 0 {
				hV := int16(asagiPost.MediaH)
				h = &hV
			}

			var tnW *int16
			if asagiPost.PreviewW > 0 {
				tnWV := int16(asagiPost.PreviewW)
				tnW = &tnWV
			}

			var tnH *int16
			if asagiPost.PreviewH > 0 {
				tnHV := int16(asagiPost.PreviewH)
				tnH = &tnHV
			}

			var since4Pass *int16
			since4PassString, hasSince4Pass := exif["since4pass"]
			if hasSince4Pass {
				since4PassV64, err := strconv.ParseInt(since4PassString, 10, 16)
				if err != nil {
					log.Printf("Error processing post %d, since4pass in exif isn't number: %s", asagiPost.DocID, since4PassString)
				} else {
					since4PassV := int16(since4PassV64)
					since4Pass = &since4PassV
				}
			}

			eienteiPosts = append(eienteiPosts, eientei.Post{
				Board:      boardConfig.Name,
				No:         int64(asagiPost.Num),
				Resto:      int64(resto),
				Time:       timestamp,
				Name:       name,
				Trip:       trip,
				Capcode:    capcode,
				Country:    country,
				Since4Pass: since4Pass,
				Com:        com,
				Tim:        tim,
				MD5:        asagiPost.MediaHash,
				Filename:   filename,
				Ext:        ext,
				Fsize:      fsize,
				W:          w,
				H:          h,
				TnW:        tnW,
				TnH:        tnH,
				Deleted:    asagiPost.Deleted,
				Spoiler:    asagiPost.Spoiler,
				Op:         asagiPost.Op,
				Sticky:     asagiPost.Sticky,
			})
		}

		_, err = eienteiTx.NewInsert().
			Model(&eienteiPosts).
			On("CONFLICT DO NOTHING").
			Exec(context.Background())

		if err != nil {
			return err
		}
	}

	if err = eienteiTx.Commit(); err != nil {
		return err
	}

	return nil
}
