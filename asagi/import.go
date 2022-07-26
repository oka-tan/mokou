package asagi

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"mokou/config"
	"mokou/koiwai"
	"mokou/utils"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/minio/minio-go/v7"
)

//Import imports the board given described in boardConfig
//from Asagi to Koiwai
func (s *Service) Import(boardConfig *config.AsagiBoardConfig) error {
	log.Printf("Importing data from Asagi to Koiwai for board %s\n", boardConfig.Name)

	//Compiles a bunch of regex
	postRestorer := newPostRestorer(boardConfig)

	tx, err := s.Pg.BeginTx(context.Background(), &sql.TxOptions{})
	defer tx.Rollback()

	if err != nil {
		return err
	}

	_, err = tx.Query(fmt.Sprintf("CREATE TABLE IF NOT EXISTS post_%s PARTITION OF post FOR VALUES IN ('%s')", boardConfig.Name, boardConfig.Name))

	//This error can occur because the partition exists with a different name
	if err != nil {
		return err
	}

	asagiTx, err := s.AsagiDb.BeginTx(context.Background(), &sql.TxOptions{})
	defer asagiTx.Rollback()

	if err != nil {
		return err
	}

	asagiPosts := make([]post, 0, s.BatchSize)
	koiwaiPosts := make([]koiwai.Post, 0, s.BatchSize)

	//We use a single now value for all of the time_createds and last_modifieds
	//to cheat postgres into compressing shit better
	now := time.Now()
	var keyset uint

	for {
		asagiPosts = asagiPosts[0:0]

		err := asagiTx.NewSelect().
			Model(&asagiPosts).
			ModelTableExpr(fmt.Sprintf("`%s` AS post", boardConfig.Name)).
			Where("`subnum` = 0"). //filters out ghost posts
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

		koiwaiPosts = koiwaiPosts[0:0]
		keyset = asagiPosts[len(asagiPosts)-1].DocID

		var mutex sync.Mutex
		var wg sync.WaitGroup

		wg.Add(len(asagiPosts))
		for _, asagiPost := range asagiPosts {
			go func(asagiPost post) {
				defer wg.Done()

				exif := make(map[string]interface{})

				if asagiPost.Exif != nil && *asagiPost.Exif != "" {
					err := json.Unmarshal([]byte(*asagiPost.Exif), &exif)

					if err != nil {
						log.Printf("Error parsing post %d, exif field: %s, error: %s", asagiPost.DocID, *asagiPost.Exif, err)
					}
				}

				postNumber := int64(asagiPost.Num)

				threadNumber := int64(asagiPost.ThreadNum)
				if threadNumber == 0 {
					threadNumber = postNumber
				}

				op := postNumber == threadNumber

				deleted := asagiPost.Deleted

				timePosted := newYorkToUTC(int64(asagiPost.Timestamp))

				name := asagiPost.Name
				if name == nil || *name == "" || *name == "Anonymous" {
					name = nil
				} else {
					filteredName := utils.FilterHtml(*name)
					name = &filteredName
				}

				tripcode := asagiPost.Trip
				if tripcode == nil || *tripcode == "" {
					tripcode = nil
				} else {
					filteredTripcode := utils.FilterHtml(*tripcode)
					tripcode = &filteredTripcode
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

				posterID := asagiPost.PosterHash
				if posterID != nil && *posterID == "" {
					posterID = nil
				}

				var country *string
				var flag *string

				neofuukaTrollCountryCode, neofuukaTrollCountryCodeAvailable := exif["troll_country_code"]
				asagiTrollCountryCode, asagiTrollCountryCodeAvailable := exif["trollCountry"]
				if neofuukaTrollCountryCodeAvailable {
					switch neofuukaTrollCountryCode.(type) {
					case string:
						{
							trollCountryCodeString := neofuukaTrollCountryCode.(string)

							if len(trollCountryCodeString) == 2 {
								flag = &trollCountryCodeString
							}
						}
					}
				} else if asagiTrollCountryCodeAvailable {
					switch asagiTrollCountryCode.(type) {
					case string:
						{
							trollCountryCodeString := asagiTrollCountryCode.(string)

							if len(trollCountryCodeString) == 2 {
								flag = &trollCountryCodeString
							}
						}
					}
				} else if asagiPost.PosterCountry != nil && len(*asagiPost.PosterCountry) == 2 {
					country = asagiPost.PosterCountry
				}

				email := asagiPost.Email
				if email != nil && *email == "" {
					email = nil
				}

				subject := asagiPost.Title
				if subject != nil && *subject == "" {
					subject = nil
				}

				var comment *string
				if asagiPost.Comment != nil && *asagiPost.Comment != "" {
					comment = postRestorer.restoreComment(&asagiPost, exif)
				}

				hasMedia := false

				var mediaTimestamp *int64
				tim, timAvailable := exif["tim"]
				if timAvailable {
					switch tim.(type) {
					case float64:
						{
							mediaTimestampV := int64(tim.(float64))
							mediaTimestamp = &mediaTimestampV
						}
					}
				} else if asagiPost.MediaOrig != nil && *asagiPost.MediaOrig != "" {
					hasMedia = true
					lastIndex := strings.LastIndex(*asagiPost.MediaOrig, ".")

					if lastIndex != -1 {
						mediaTimestampV, err := strconv.ParseInt((*asagiPost.MediaOrig)[:lastIndex], 10, 64)

						if err == nil {
							mediaTimestamp = &mediaTimestampV
						}
					}
				}

				media4chanHash := base64StringToBytes(asagiPost.MediaHash)

				var mediaFileName *string
				var mediaExtension *string
				if asagiPost.MediaFilename != nil && *asagiPost.MediaFilename != "" {
					hasMedia = true
					lastIndex := strings.LastIndex(*asagiPost.MediaFilename, ".")

					if lastIndex != -1 {
						if lastIndex == 0 {
							mediaFileName = nil
						} else {
							mediaFileNameV := (*asagiPost.MediaFilename)[:lastIndex]
							mediaFileName = &mediaFileNameV
						}
						mediaExtensionV := (*asagiPost.MediaFilename)[lastIndex+1:]

						if mediaExtensionV == "web" {
							mediaExtensionV = "webm"
						}

						mediaExtension = &mediaExtensionV
					} else {
						log.Printf("The goofy media filename is %s\n", *asagiPost.MediaFilename)
					}
				}

				var mediaSize *int
				if asagiPost.MediaSize > 0 {
					hasMedia = true
					mediaSizeV := int(asagiPost.MediaSize)
					mediaSize = &mediaSizeV
				}

				var mediaWidth *int16
				var mediaHeight *int16
				if asagiPost.MediaW > 0 && asagiPost.MediaH > 0 {
					hasMedia = true
					mediaWidthV := int16(asagiPost.MediaW)
					mediaHeightV := int16(asagiPost.MediaH)
					mediaWidth = &mediaWidthV
					mediaHeight = &mediaHeightV
				}

				var thumbnailWidth *int16
				var thumbnailHeight *int16
				if asagiPost.PreviewW > 0 && asagiPost.PreviewH > 0 {
					hasMedia = true
					thumbnailWidthV := int16(asagiPost.PreviewW)
					thumbnailHeightV := int16(asagiPost.PreviewH)

					thumbnailWidth = &thumbnailWidthV
					thumbnailHeight = &thumbnailHeightV
				}

				var spoiler *bool
				if hasMedia {
					spoiler = &asagiPost.Spoiler
				}

				var closed *bool
				if op {
					closed = &asagiPost.Locked
				}

				var sticky *bool
				if op {
					sticky = &asagiPost.Sticky
				}

				var posters *int16
				if op {
					postersCount, postersCountAvailable := exif["uniqueIps"]
					if postersCountAvailable {
						switch postersCount.(type) {
						case string:
							{
								postersString := postersCount.(string)
								postersV64, err := strconv.ParseInt(postersString, 10, 16)
								if err == nil {
									postersV := int16(postersV64)
									posters = &postersV
								}
							}
						}
					}
				}

				var replies *int16
				if op {
					var zero int16
					replies = &zero
				}

				var since4Pass *int16
				since4PassRaw, hasSince4Pass := exif["since4pass"]
				if hasSince4Pass {
					switch since4PassRaw.(type) {
					case string:
						{
							since4PassString := since4PassRaw.(string)
							since4PassV64, err := strconv.ParseInt(since4PassString, 10, 16)
							if err != nil {
								log.Printf("Error processing post %d, since4pass in exif isn't number: %s", asagiPost.DocID, since4PassString)
							} else {
								since4PassV := int16(since4PassV64)
								since4Pass = &since4PassV
							}
						}
					}
				}

				var thumbnailInternalHash *[]byte
				if boardConfig.ImportThumbnails && asagiPost.PreviewOrig != nil && len(*asagiPost.PreviewOrig) > 6 {
					thumbnail := findThumbnail(s.AsagiImagesFolder, boardConfig.Name, *asagiPost.PreviewOrig)
					thumbnailInternalHash = utils.HashFile(thumbnail)

					if thumbnailInternalHash != nil {
						alreadyExists, err := s.Pg.NewSelect().
							Model(&koiwai.Thumbnail{}).
							Where("hash = ?", *thumbnailInternalHash).
							Exists(context.Background())

						if err != nil {
							log.Fatalf("Error checking if thumbnail exists: %s", err)
						}

						if !alreadyExists {
							_, err := s.S3MediaClient.FPutObject(
								context.Background(),
								*s.S3ThumbnailsBucket,
								base64.URLEncoding.EncodeToString(*thumbnailInternalHash),
								thumbnail,
								minio.PutObjectOptions{
									CacheControl: "private, immutable, max-age=604800",
									ContentType:  "image/jpg",
								},
							)

							if err != nil {
								log.Fatalf("Error putting file %s on s3: %s", thumbnail, err)
							}

							if _, err := s.Pg.NewInsert().Model(&koiwai.Thumbnail{Hash: *thumbnailInternalHash}).On("CONFLICT (hash) DO NOTHING").Exec(context.Background()); err != nil {
								log.Fatalf("Error putting file hash on db: %s", err)
							}
						}
					}
				}

				var mediaInternalHash *[]byte
				if boardConfig.ImportImages && asagiPost.MediaOrig != nil && len(*asagiPost.MediaOrig) > 6 {
					image := findImage(s.AsagiImagesFolder, boardConfig.Name, *asagiPost.MediaOrig)
					mediaInternalHash = utils.HashFile(image)

					if mediaInternalHash != nil {
						alreadyExists, err := s.Pg.NewSelect().
							Model(&koiwai.Media{}).
							Where("hash = ?", *mediaInternalHash).
							Exists(context.Background())

						if err != nil {
							log.Fatalf("Error checking if media exists: %s", err)
						}

						if !alreadyExists {
							_, err := s.S3MediaClient.FPutObject(
								context.Background(),
								*s.S3MediaBucket,
								base64.URLEncoding.EncodeToString(*mediaInternalHash),
								image,
								minio.PutObjectOptions{
									CacheControl: "private, immutable, max-age=604800",
									ContentType:  mime.TypeByExtension("." + *mediaExtension),
								},
							)

							if err != nil {
								log.Fatalf("Error putting file %s on s3: %s", image, err)
							}

							if _, err := s.Pg.NewInsert().Model(&koiwai.Media{Hash: *mediaInternalHash}).On("CONFLICT (hash) DO NOTHING").Exec(context.Background()); err != nil {
								log.Fatalf("Error putting file hash on db: %s", err)
							}
						}
					}
				}

				mutex.Lock()
				koiwaiPosts = append(koiwaiPosts, koiwai.Post{
					Board:                 boardConfig.Name,
					PostNumber:            postNumber,
					ThreadNumber:          threadNumber,
					Op:                    op,
					Deleted:               deleted,
					Hidden:                false,
					TimePosted:            timePosted,
					LastModified:          now,
					CreatedAt:             now,
					Name:                  name,
					Tripcode:              tripcode,
					Capcode:               capcode,
					PosterID:              posterID,
					Country:               country,
					Flag:                  flag,
					Email:                 email,
					Subject:               subject,
					Comment:               comment,
					HasMedia:              hasMedia,
					MediaDeleted:          nil,
					TimeMediaDeleted:      nil,
					MediaTimestamp:        mediaTimestamp,
					Media4chanHash:        media4chanHash,
					MediaInternalHash:     mediaInternalHash,
					ThumbnailInternalHash: thumbnailInternalHash,
					MediaExtension:        mediaExtension,
					MediaFileName:         mediaFileName,
					MediaSize:             mediaSize,
					MediaHeight:           mediaHeight,
					MediaWidth:            mediaWidth,
					ThumbnailHeight:       thumbnailHeight,
					ThumbnailWidth:        thumbnailWidth,
					Spoiler:               spoiler,
					CustomSpoiler:         nil,
					Sticky:                sticky,
					Closed:                closed,
					Posters:               posters,
					//This field will be corrected with an update statement at the end.
					Replies:    replies,
					Since4Pass: since4Pass,
				})
				mutex.Unlock()
			}(asagiPost)
		}

		wg.Wait()

		_, err = tx.NewInsert().
			Model(&koiwaiPosts).
			On("CONFLICT DO NOTHING").
			Returning("NULL").
			Exec(context.Background())

		if err != nil {
			return err
		}
	}

	replyCountCte := tx.NewSelect().
		Model(&koiwai.Post{}).
		Column("thread_number").
		ColumnExpr("COUNT(*) - 1 AS replies").
		Where("board = ?", boardConfig.Name).
		Group("thread_number")

	_, err = tx.NewUpdate().
		With("_reply_counts", replyCountCte).
		Model(&koiwai.Post{}).
		TableExpr("_reply_counts").
		Set("replies = _reply_counts.replies").
		Where("board = ?", boardConfig.Name).
		Where("op").
		Where("post_number = _reply_counts.thread_number").
		Exec(context.Background())

	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
