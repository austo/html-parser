package nrsv

import (
	"code.google.com/p/go.net/html"
	"io"
	// "io/ioutil"
	"net/http"
	// "regexp"
	// "fmt"
)

func getRawVerseText(ch Chapter) (*html.Node, error) {
	res, err := http.Get(ch.url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return getTextNode(res.Body)
	// bytes, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	return
	// }
	// text = string(bytes)
	// return
}

func getTextNode(r io.Reader) (n *html.Node, err error) {
	doc, err := html.Parse(r)
	if err != nil {
		return
	}
	n = findTextDiv(doc)
	return
}

func findTextDiv(n *html.Node) (node *html.Node) {
	if n.Type == html.ElementNode && n.Data == "div" {
		for _, a := range n.Attr {
			if a.Key == "class" && a.Val == "passage-text" {
				node = n
				return
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		node = findTextDiv(c)
		if node != nil {
			return
		}
	}
	return
}
