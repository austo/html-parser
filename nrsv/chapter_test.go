package nrsv

import (
	"code.google.com/p/go.net/html"
	"io"
	"os"
	"testing"
)

const (
	chapFilename = "genesis_1.html"
)

var (
	chap = Chapter{
		"Genesis",
		1,
		"https://www.biblegateway.com/passage/?search=Genesis+1&version=NRSV"}
)

func TestGetChapterText(t *testing.T) {
	text, err := getPassageTextFromWeb(chap)
	checkError(t, err)
	t.Log(text)
}

func TestGetTextNode(t *testing.T) {
	node := getTestPassageTextNode(t)
	t.Log(node)
}

func getTestPassageTextNode(t *testing.T) *html.Node {
	f := getChapterFile(t)
	defer f.Close()
	n, err := getPassageTextNode(f)
	checkError(t, err)
	return n
}

func getChapterFile(t *testing.T) io.ReadCloser {
	f, err := os.Open(chapFilename)
	checkError(t, err)
	return f
}

func checkError(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}
