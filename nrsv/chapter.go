package nrsv

import (
	"golang.org/x/net/html"
	"io"
	// "io/ioutil"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type verse struct {
	number uint16
	text   string
}

var (
	textDivClassRe = regexp.MustCompile(`(?i:.*?result\-text\-style\-normal.*?$)`)
	orRe           = regexp.MustCompile(`(?i:^or$)`)
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

func logPassageText(node *html.Node) {
	insideFootnote := false
	// TODO: remember verse number and only advance to next verse
	// once next verse number has been found.
	// TODO: break reliably once footnotes have been encountered
	var f func(n *html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			s := strings.TrimSpace(n.Data)
			if len(s) > 0 {
				if s == "Footnotes:" || orRe.MatchString(s) {
					return
				}
				if s[0] == '[' {
					insideFootnote = true
				} else if s[0] == ']' {
					insideFootnote = false
				} else if !insideFootnote {
					fmt.Printf("%s\n", s)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(node)
}
