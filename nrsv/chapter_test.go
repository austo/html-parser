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
	chapters = []Chapter{Chapter{
		"Genesis",
		1,
		"https://www.biblegateway.com/passage/?search=Genesis+1&version=NRSV"},
		Chapter{
			"1 Corinthians",
			12,
			"https://www.biblegateway.com/passage/?search=1+Corinthians+12&version=NRSV"}}
)

func TestGetChapterText(t *testing.T) {
	for _, chap := range chapters {
		getChapterVerses(t, chap)
	}
}

func TestGetTextNode(t *testing.T) {
	f := getChapterFile(t)
	defer f.Close()
	node, err := getTextNode(f)
	checkError(t, err)
	t.Logf("%#v\n", node)
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

func getChapterVerses(t *testing.T, chap Chapter) {
	textNode, err := getRawVerseTextNodeFromWeb(chap)
	checkError(t, err)
	verses := getVersesFromPassageTextNode(textNode)
	t.Log(verses)
}
