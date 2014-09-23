package nrsv

import (
	"code.google.com/p/go.net/html"
	"io"
	// "io/ioutil"
	"net/http"
	"regexp"
	// "fmt"
)

var (
	textDivClassRe = regexp.MustCompile(`(?i:.*?result\-text\-style\-normal.*?$)`)
)

func getPassageTextFromWeb(ch Chapter) (*html.Node, error) {
	res, err := http.Get(ch.url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return getPassageTextNode(res.Body)
}

func getPassageTextNode(r io.Reader) (n *html.Node, err error) {
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
