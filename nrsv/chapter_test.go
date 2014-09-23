package nrsv

import (
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
	text, err := getRawVerseText(chap)
	checkError(t, err)
	t.Log(text)
}

func TestGetTextNode(t *testing.T) {
	f := getChapterFile(t)
	defer f.Close()
	node, err := getTextNode(f)
	checkError(t, err)
	t.Log(node)
}

func getChapterFile(t *testing.T) io.ReadCloser {
	f, err := os.Open(chapFilename)
	if err != nil {
		t.Error(err)
	}
	return f
}

func checkError(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}
