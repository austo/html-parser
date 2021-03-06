package nrsv

import (
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type Chapter struct {
	book  string
	index uint8
	url   string
}

func (c Chapter) String() string {
	return fmt.Sprintf("%s, Chapter %d: %s", c.book, c.index, c.url)
}

type ChapterHandler func(Chapter)
type bookList map[string][]Chapter

const (
	baseUrl     = `https://www.biblegateway.com`
	bookListUrl = `/versions/New-Revised-Standard-Version-NRSV-Bible/#booklist`
)

var (
	bookNameRe = regexp.MustCompile(`(?i:^.*?search=((?:\d+\+)*[\w\+]+)\+\d+&.*$)`)
)

func GetChaptersFromWeb(
	handler ChapterHandler) (done <-chan bool, err error) {

	res, err := http.Get(baseUrl + bookListUrl)
	if err != nil {
		return
	}
	doc, err := html.Parse(res.Body)
	res.Body.Close()
	if err != nil {
		return
	}
	done = getChapters(doc, handler)
	return
}

func getChapters(n *html.Node, handler ChapterHandler) <-chan bool {
	done := make(chan bool)
	go func() {
		findChapters(n, handler)
		done <- true
	}()
	return done
}

func findChapters(n *html.Node, handler ChapterHandler) {
	if n.Type == html.ElementNode && n.Data == "a" {
		handleAnchorNode(n, handler)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		findChapters(c, handler)
	}
}

func handleAnchorNode(n *html.Node, handler ChapterHandler) {
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

func buildChapterFromHref(a html.Attribute, chapIndex uint8) (Chapter, error) {
	m := bookNameRe.FindAllStringSubmatch(a.Val, -1)
	if m == nil {
		return Chapter{}, fmt.Errorf("href is not book name")
	}
	bookName := strings.Replace(m[0][1], "+", " ", -1)
	return Chapter{bookName, chapIndex, fmt.Sprintf("%s%s", baseUrl, a.Val)}, nil
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

func getTocPageText(url string) (string, error) {
	res, err := http.Get(url)
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

func getChaptersFromPageSource(s string,
	handler ChapterHandler) (<-chan bool, error) {

	doc, err := html.Parse(strings.NewReader(s))
	if err != nil {
		return nil, err
	}
	done := getChapters(doc, handler)
	return done, err
}

func makeHandler() (<-chan Chapter, ChapterHandler) {
	chapters := make(chan Chapter)
	var f = func(chap Chapter) {
		chapters <- chap
	}
	return chapters, f
}
