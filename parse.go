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
	done, err := nrsv.GetChaptersFromWeb(handler)
	if err != nil {
		log.Fatal(err)
	}
	finished := false
	for !finished {
		select {
		case ch := <-chapters:
			fmt.Printf("%s\n", ch)
		case <-done:
			finished = true
		}
	}
}
