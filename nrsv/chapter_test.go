package nrsv

import (
	"testing"
)

var (
	chap = Chapter{
		"Genesis",
		1,
		"https://www.biblegateway.com/passage/?search=Genesis+1&version=NRSV"}
)

func TestGetChapterText(t *testing.T) {
	text, err := getRawVerseText(chap, t)
	if err != nil {
		t.Error(err)
	}
	t.Log(text)
}

// func TestGetTextNode(t *testing.T) {

// }
