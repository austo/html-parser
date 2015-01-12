package main

import (
	"fmt"
	"github.com/austo/html-parser/config"
	"github.com/austo/html-parser/db"
	"github.com/austo/html-parser/nrsv"
	// "log"
	"os"
	"sync"
)

func main() {
	d, err := getDb()
	if err != nil {
		errExit("error obtaining database connection\n", err)
	}
	chapters := make(chan nrsv.Chapter)
	handler := func(chap nrsv.Chapter) {
		chapters <- chap
	}
	chaptersDone, err := nrsv.GetChaptersFromWeb(handler)
	if err != nil {
		errExit("error getting chapters\n", err)
	}
	finished := false
	var chapterQueue []nrsv.Chapter
	for !finished {
		select {
		case ch := <-chapters:
			chapterQueue = append(chapterQueue, ch)
		case <-chaptersDone:
			fmt.Println("CHAPTERS DONE LOADING")
			finished = true
		}
	}
	getVerses(d, chapterQueue)
}

// TODO: handle one chapter as soon as it is available on channel
func getVerses(d *db.Db, chapters []nrsv.Chapter) {
	records, done := make(chan nrsv.VerseRecord), make(chan bool)
	nDone := 0
	for _, chap := range chapters {
		go nrsv.GetVerseRecordsFromWeb(chap, records, done)
	}
	// var wg sync.WaitGroup
	var verseRecords []interface{}
	for nDone < len(chapters) {
		select {
		case vr := <-records:
			verseRecords = append(verseRecords, vr)
			// wg.Add(1)
			// go insertVerseRecord(d, vr, &wg)
		case <-done:
			nDone++
		}
	}
	// wg.Wait()
	err := d.InsertVerseRecords(verseRecords)
	if err != nil {
		errExit("error inserting verse records", err)
	}
	fmt.Printf("processed %d chapters, %d verses\n", nDone, len(verseRecords))
}

func insertVerseRecord(d *db.Db, vr nrsv.VerseRecord, wg *sync.WaitGroup) {
	defer wg.Done()
	err := d.InsertVerseRecord(vr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error inserting verse: %v, %v\n", err, vr)
	}
}

func getDb() (d *db.Db, err error) {
	cfg, err := config.ReadEnvironmentFromFile(config.Filename, "local")
	if err != nil {
		return
	}
	d, err = db.Connect(cfg)
	return
}

func errExit(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args)
	os.Exit(1)
}
