package main

import (
	"testing"
)

const (
	html = `<!DOCTYPE html>
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
)

func TestGetChaptersFromPageSource(t *testing.T) {
	chapters := make([]Chapter)
	handler := func(chap Chapter) {
		chapters <- chap
	}

}
