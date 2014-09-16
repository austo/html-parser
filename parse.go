package main

import (
	"code.google.com/p/go.net/html"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const (
	baseUrl     = `https://www.biblegateway.com`
	bookListUrl = `/versions/New-Revised-Standard-Version-NRSV-Bible/#booklist`
)

func main() {
	s, err := getMainPage()
	if err != nil {
		log.Fatal(err)
	}
	doc, err := html.Parse(strings.NewReader(s))
	if err != nil {
		log.Fatal(err)
	}
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			isChap, chap := isChapterLink(n)
			if isChap {
				for _, a := range n.Attr {
					if a.Key == "href" {
						fmt.Printf("%d: %s%s\n", chap, baseUrl, a.Val)
						break
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
}

func getMainPage() (string, error) {
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
