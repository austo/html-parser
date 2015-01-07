package nrsv

import (
	"golang.org/x/net/html"
	"io"
	// "io/ioutil"
	// "fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type verse struct {
	number uint16
	text   string
}

type Verse struct {
	bookIndex    uint8
	bookName     string
	chapterIndex uint8
	verseIndex   uint16
	text         string
}

var (
	textDivClassRe   = regexp.MustCompile(`(?i:.*?result\-text\-style\-normal.*?$)`)
	numberRe         = regexp.MustCompile(`^\d+$`)
	verseTextClassRe = regexp.MustCompile(`^text\s\w+\-\d+\-(\d+).*$`)
)

func getRawVerseTextNodeFromWeb(ch Chapter) (*html.Node, error) {
	res, err := http.Get(ch.url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return getTextNode(res.Body)
}

func getTextNode(r io.Reader) (n *html.Node, err error) {
	doc, err := html.Parse(r)
	if err != nil {
		return
	}
	n = findPassageTextDiv(doc)
	return
}

func findPassageTextDiv(n *html.Node) (node *html.Node) {
	if n.Type == html.ElementNode && n.Data == "div" {
		for _, a := range n.Attr {
			if a.Key == "class" && textDivClassRe.MatchString(a.Val) {
				node = n
				return
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		node = findPassageTextDiv(c)
		if node != nil {
			return
		}
	}
	return
}

func getVersesFromPassageTextNode(node *html.Node) (verses []verse) {
	var currentVerseNum uint16 = 1
	v := verse{currentVerseNum, ""}
	var f func(n *html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			if isVerse, verseNum := isVerseNode(n); isVerse {
				s := strings.TrimSpace(n.Data)
				if verseNum > currentVerseNum { // clean up old verse and make new verse
					v.text = strings.TrimSpace(v.text)
					verses = append(verses, v)
					currentVerseNum = verseNum
					v = verse{verseNum, s}
				} else {
					v.appendText(s)
				}
			} else if isSmallCaps, text := isSmallCapsNode(n); isSmallCaps {
				v.appendText(text)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(node)
	verses = append(verses, v)
	return
}

func isVerseNode(n *html.Node) (found bool, verseNum uint16) {
	parent := n.Parent
	if parent.Data != "span" || parent.Parent.Data != "p" {
		return
	}
	for _, a := range parent.Attr {
		if a.Key == "class" {
			m := verseTextClassRe.FindStringSubmatch(a.Val)
			if len(m) > 0 {
				vNum, _ := strconv.ParseInt(m[1], 10, 16)
				verseNum = uint16(vNum)
				found = true
				return
			}
		}
	}
	return
}

func isSmallCapsNode(n *html.Node) (found bool, text string) {
	parent := n.Parent
	if parent.Data != "span" {
		return
	}
	for _, a := range parent.Attr {
		if a.Key == "class" && a.Val == "small-caps" {
			found = true
			text = strings.TrimSpace(n.Data)
			return
		}
	}
	return
}

func (v *verse) appendText(s string) {
	if v.text == "" { // first verse node
		v.text += s
	} else { // verse node after footnote
		v.text += " " + s
	}
}

// func (v verse) getRecord(c Chapter) (rv Verse) {
// }
