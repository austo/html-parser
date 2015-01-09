package main

import (
	"fmt"
	"github.com/austo/html-parser/nrsv"
	"log"
)

func main() {
	chapters := make(chan nrsv.Chapter)
	handler := func(chap nrsv.Chapter) {
		chapters <- chap
	}
	chaptersDone, err := nrsv.GetChaptersFromWeb(handler)
	if err != nil {
		log.Fatal(err)
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
	getVerses(chapterQueue)
}

// TODO: handle one chapter as soon as it is available on channel
func getVerses(chapters []nrsv.Chapter) {
	records, done := make(chan nrsv.VerseRecord), make(chan bool)
	nDone := 0
	for _, chap := range chapters {
		go nrsv.GetVerseRecordsFromWeb(chap, records, done)
	}
	for nDone < len(chapters) {
		select {
		case vr := <-records:
			fmt.Println(vr)
		case <-done:
			nDone++
		}
	}
}
