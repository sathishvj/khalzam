package musiclibrary

import (
	"database/sql"
	"fmt"
	"github.com/kisasexypantera94/khalzam/fingerprint"
	"github.com/remeh/sizedwaitgroup"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strings"
)

const (
	dbUser     = "kisasexypantera94"
	dbPassword = ""
	dbName     = "khalzam"
)

// MusicLibrary ...
type MusicLibrary struct {
	db *sql.DB
}

// Open return pointer to existing library
func Open() (*MusicLibrary, error) {
	fmt.Printf("Initializing library...\n\n")

	dbinfo := fmt.Sprintf("user=%s dbname=%s sslmode=disable", dbUser, dbName)
	db, err := sql.Open("postgres", dbinfo)
	return &MusicLibrary{db}, err
}

// Close closes library
func (lib *MusicLibrary) Close() {
	lib.db.Close()
}

// Index inserts song into library
func (lib *MusicLibrary) Index(filename string) error {
	dotIdx := strings.LastIndex(filename, ".")
	slashIdx := strings.LastIndex(filename, "/")
	if dotIdx == -1 {
		return fmt.Errorf("Index: invalid file '%s'", filename)
	}
	songName := filename[slashIdx+1 : dotIdx]
	var songID int
	err := lib.db.QueryRow("INSERT INTO songs(song) VALUES($1) returning sid;", songName).Scan(&songID)
	if err != nil {
		return err
	}

	hashArray := fingerprint.Fingerprint(filename)
	for time, hash := range hashArray {
		lib.db.QueryRow("INSERT INTO hashes(hash, time, sid) VALUES($1, $2, $3) returning hid;", hash, time, songID).Scan()
	}

	return nil
}

// IndexDir indexes whole directory
func (lib *MusicLibrary) IndexDir(path string) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return fmt.Errorf("IndexDir: invalid directory '%s'", path)
	}

	fmt.Printf("Indexing '%s'...\n", path)

	wg := sizedwaitgroup.New(8)
	for _, f := range files {
		filename := path + "/" + f.Name()
		if filepath.Ext(f.Name()) == ".mp3" {
			wg.Add()
			go func() {
				defer wg.Done()
				lib.Index(filename)
				// checkErr(err)
			}()
		}
	}
	wg.Wait()

	return nil
}

type table struct {
	absoluteBest  uint          // highest number of matches among timedeltaBest
	timedeltaBest map[uint]uint // highest number of matches for every timedelta
}

type candidate struct {
	songID   uint
	matchNum uint
}

// Recognize searches library and returns table
func (lib *MusicLibrary) Recognize(filename string) (result string) {
	cnt := make(map[uint]*table)

	hashArray := fingerprint.Fingerprint(filename)
	for t, h := range hashArray {
		rows, err := lib.db.Query("SELECT * FROM hashes WHERE hash=$1;", h)
		checkErr(err)

		for rows.Next() {
			var hid int
			var hash string
			var sid uint
			var time uint
			err = rows.Scan(&hid, &hash, &time, &sid)
			checkErr(err)
			if cnt[sid] == nil {
				cnt[sid] = &table{}
				cnt[sid].timedeltaBest = make(map[uint]uint)
			}

			cnt[sid].timedeltaBest[time-(uint)(t)]++
			if cnt[sid].timedeltaBest[time-(uint)(t)] > cnt[sid].absoluteBest {
				cnt[sid].absoluteBest = cnt[sid].timedeltaBest[time-(uint)(t)]
			}
		}
	}

	matchings := make([]*candidate, 0)
	fmt.Printf("Recognizing '%s'...\n", filename)
	for song, table := range cnt {
		matchings = append(matchings, &candidate{song, table.absoluteBest})
	}

	sort.Slice(matchings, func(i, j int) bool {
		return matchings[i].matchNum > matchings[j].matchNum
	})

	var songName string
	lib.db.QueryRow("SELECT song FROM songs WHERE sid=$1;", matchings[0].songID).Scan(&songName)

	result = fmt.Sprintf("Best match: %s (%d%% matched)\n", songName, (int)(100*(float64)(matchings[0].matchNum)/(float64)(len(hashArray))))
	return
}

// Delete deletes song from library
func (lib *MusicLibrary) Delete(song string) error {
	fmt.Printf("Deleting '%s'\n", song)
	statement, err := lib.db.Prepare("DELETE FROM songs WHERE song=$1;")
	checkErr(err)
	_, err = statement.Exec(song)
	return err
}

func checkErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
