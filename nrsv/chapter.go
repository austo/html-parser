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

var (
	textDivClassRe = regexp.MustCompile(`(?i:.*?result\-text\-style\-normal.*?$)`)
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

func logPassageText(n *html.Node) {
	if n.Type == html.TextNode {
		s := strings.TrimSpace(n.Data)
		if len(s) > 0 {
			fmt.Printf("%s\n", s)
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		logPassageText(c)
	}
}
