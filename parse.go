package main

import (
	"code.google.com/p/go.net/html"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// TODO: before writing more code,
// refactor into chapter package and add test coverage.
// chapter package sould deal with http in terms of interfaces.
type chapter struct {
	book  string
	index uint8
	url   string
}

func (c chapter) String() string {
	return fmt.Sprintf("%s, Chapter %d: %s", c.book, c.index, c.url)
}

type chapterHandler func(chapter)
type bookList map[string][]chapter

const (
	baseUrl     = `https://www.biblegateway.com`
	bookListUrl = `/versions/New-Revised-Standard-Version-NRSV-Bible/#booklist`
)

var (
	bookNameRe = regexp.MustCompile(`(?i:^.*?search=((?:\d+\+)*\w+)\+\d+&.*$)`)
)

func main() {
	s, err := getTocPageText()
	if err != nil {
		log.Fatal(err)
	}
	getChaptersFromPageSource(s)
}

func getChaptersFromPageSource(s string) {
	doc, err := html.Parse(strings.NewReader(s))
	if err != nil {
		log.Fatal(err)
	}
	chapters, handler := makeHandler()
	done := getChapters(doc, handler)
	finished := false
	for !finished {
		select {
		case chap := <-chapters:
			fmt.Printf("%s\n", chap)
		case <-done:
			finished = true
		}
	}
}

func makeHandler() (<-chan chapter, chapterHandler) {
	chapters := make(chan chapter)
	var f = func(chap chapter) {
		chapters <- chap
	}
	return chapters, f
}

func getChapters(n *html.Node, handler chapterHandler) <-chan bool {
	done := make(chan bool)
	go func() {
		findChapters(n, handler)
		done <- true
	}()
	return done
}

func findChapters(n *html.Node, handler chapterHandler) {
	if n.Type == html.ElementNode && n.Data == "a" {
		handleAnchorNode(n, handler)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		findChapters(c, handler)
	}
}

func handleAnchorNode(n *html.Node, handler chapterHandler) {
	isChap, chap := isChapterLink(n)
	if !isChap {
		return
	}
	for _, a := range n.Attr {
		if a.Key != "href" {
			continue
		}
		c, err := buildChapterFromHref(a, chap)
		if err == nil {
			handler(c)
		}
		break
	}
}

func buildChapterFromHref(a html.Attribute, chapIndex uint8) (chapter, error) {
	m := bookNameRe.FindAllStringSubmatch(a.Val, -1)
	if m == nil {
		return chapter{}, fmt.Errorf("href is not book name")
	}
	bookName := strings.Replace(m[0][1], "+", " ", -1)
	return chapter{bookName, chapIndex, fmt.Sprintf("%s%s", baseUrl, a.Val)}, nil
}

func getTocPageText() (string, error) {
	res, err := http.Get(baseUrl + bookListUrl)
	if err != nil {
		return "", err
	}
	text, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return "", err
	}
	return string(text), nil
}

func isChapterLink(n *html.Node) (bool, uint8) {
	if n.FirstChild == nil {
		return false, 0
	}
	data := n.FirstChild.Data
	chap, err := strconv.ParseUint(data, 10, 8)
	if err != nil {
		return false, 0
	}
	return true, uint8(chap)
}
