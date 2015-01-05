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

var (
	textDivClassRe = regexp.MustCompile(`(?i:.*?result\-text\-style\-normal.*?$)`)
	numberRe       = regexp.MustCompile(`^\d+$`)
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
	insideFootnote, finished := false, false
	var v verse
	// TODO: remember verse number and only advance to next verse
	// once next verse number has been found.
	var f func(n *html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			s := strings.TrimSpace(n.Data)
			if len(s) > 0 {
				if s == "Footnotes:" {
					finished = true
					return
				}
				if s[0] == '[' {
					insideFootnote = true
				} else if s[0] == ']' {
					insideFootnote = false
				} else if !insideFootnote {
					if numberRe.MatchString(s) {
						n, _ := strconv.ParseInt(s, 10, 16)
						if n > 1 {
							v.text = strings.TrimSpace(v.text)
							verses = append(verses, v)
						}
						v = verse{uint16(n), ""}
					} else {
						v.text += s + " "
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if finished {
				return
			}
			f(c)
		}
	}
	f(node)
	verses = append(verses, v)
	return
}
