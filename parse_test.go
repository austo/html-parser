package main

import (
	"testing"
)

const (
	htmlStr = `<!DOCTYPE html>
<html>

<body>
    <div>
        <a href="/passage/?search=Genesis+1&amp;version=NRSV" title="Genesis 1">1</a>
        <a href="/passage/?search=Genesis+2&amp;version=NRSV" title="Genesis 2">2</a>
        <a href="/passage/?search=Genesis+3&amp;version=NRSV" title="Genesis 3">3</a>
        <a href="/passage/?search=Genesis+4&amp;version=NRSV" title="Genesis 4">4</a>
        <a href="/passage/?search=Genesis+5&amp;version=NRSV" title="Genesis 5">5</a>
        <a href="/passage/?search=Genesis+6&amp;version=NRSV" title="Genesis 6">6</a>
    </div>
</body>

</html>`
	nExpectedChapters = 6
)

func TestGetChaptersFromPageSource(t *testing.T) {
	var chapters []Chapter
	handler := func(chap Chapter) {
		chapters = append(chapters, chap)
	}
	done, err := getChaptersFromPageSource(htmlStr, handler)
	if err != nil {
		t.Error(err)
	}
	finished := false
	for !finished {
		select {
		case <-done:
			finished = true
		}
	}
	nFoundChapters := len(chapters)
	if nFoundChapters != nExpectedChapters {
		t.Errorf("incorrect number of chapters: expected %d, had %d\n",
			nExpectedChapters, nFoundChapters)
	}
}

func TestGetChaptersFromPageSourceError(t *testing.T) {
	handler := func(chap Chapter) {
		t.Logf("%s", chap)
	}
	_, err := getChaptersFromPageSource(`<!-- [CDATA[foo]]`,
		handler)
	if err != nil {
		t.Errorf("error parsing text string: %v\n", err)
	}
}
