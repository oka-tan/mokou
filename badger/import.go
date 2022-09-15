//Package badger exposes functionality for importing database
//dumps from genetically_enhanced_badger
package badger

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"mokou/config"
	"mokou/koiwai"
	"os"
	"sync"
	"time"

	"github.com/samber/lo"
)

//Import imports the board described in boardConfig
func (s Service) Import(boardConfig config.BadgerBoardConfig) {
	log.Printf("Importing data from Badger to Koiwai for board %s\n", boardConfig.Name)

	tx, err := s.Pg.BeginTx(context.Background(), &sql.TxOptions{})
	defer tx.Rollback()

	if err != nil {
		panic(err)
	}

	_, err = tx.Query(fmt.Sprintf("CREATE TABLE IF NOT EXISTS post_%s PARTITION OF post FOR VALUES IN ('%s')", boardConfig.Name, boardConfig.Name))

	if err != nil {
		panic(err)
	}

	boardJsonFolder := fmt.Sprintf("%s/%s", s.JsonFolder, boardConfig.Name)
	jsonFiles, err := os.ReadDir(boardJsonFolder)

	if err != nil {
		panic(err)
	}

	now := time.Now()

	jsonChannel := make(chan string, len(jsonFiles))
	var wg sync.WaitGroup

	for _, jsonFile := range jsonFiles {
		jsonChannel <- fmt.Sprintf("%s/%s", boardJsonFolder, jsonFile.Name())
	}

	wg.Add(len(jsonFiles))

	for i := 0; i < 50; i++ {
		go func() {
			for jsonFile := range jsonChannel {
				bs, err := os.ReadFile(jsonFile)

				if err != nil {
					panic(err)
				}

				bs = bytes.TrimPrefix(bs, []byte("\xef\xbb\xbf"))

				var t thread

				if err := json.Unmarshal(bs, &t); err != nil {
					panic(fmt.Errorf("Error parsing file %s: %s", jsonFile, err))
				}

				koiwaiPosts := lo.Map(t.Posts, func(p post, _ int) koiwai.Post {
					return toPostModel(boardConfig.Name, now, p)
				})

				_, err = tx.NewInsert().
					Model(&koiwaiPosts).
					Returning("NULL").
					Exec(context.Background())

				if err != nil {
					panic(err)
				}

				wg.Done()
			}
		}()
	}

	wg.Wait()

	if err := tx.Commit(); err != nil {
		panic(err)
	}
}
