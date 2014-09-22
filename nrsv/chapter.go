package nrsv

import (
	"code.google.com/p/go.net/html"
	"io"
	// "io/ioutil"
	"net/http"
	// "regexp"
	"fmt"
	"testing"
)

func getRawVerseText(ch Chapter, t *testing.T) (*html.Node, error) {
	res, err := http.Get(ch.url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return getTextNode(res.Body, t)
	// bytes, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	return
	// }
	// text = string(bytes)
	// return
}

func getTextNode(r io.Reader, t *testing.T) (n *html.Node, err error) {
	doc, err := html.Parse(r)
	if err != nil {
		return
	}
	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		if t != nil {
			fmt.Printf("%#v\n", c)
		}
		if c.Type != html.ElementNode || c.Data != "div" {
			continue
		}
		for _, a := range c.Attr {
			if a.Key == "class" && a.Val == "passage-text" {
				n = c
				fmt.Printf("%#v\n", *n)
				return
			}
		}
	}
	return
}

func findTextDiv(n *html.Node) (n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "div" {
		for _, a := range n.Attr {
			if a.Key == "class" && a.Val == "passage-text" {
				return n
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		findTextDiv(c)
	}
}
