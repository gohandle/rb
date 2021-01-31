package rbtest

import (
	"io"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

var fatalf = func(tb testing.TB, format string, args ...interface{}) {
	tb.Fatalf(format, args...)
}

type Doc struct {
	*goquery.Document
	tb testing.TB
}

func MustParseHtml(tb testing.TB, r io.Reader) *Doc {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		fatalf(tb, "rbtest: failed to create document from reader: %v", err)
	}

	return &Doc{doc, tb}
}

func (d *Doc) MustHtml() string {
	ret, err := d.Html()
	if err != nil {
		fatalf(d.tb, "rbtest: failed to convert to HTML: %v", err)
	}

	return ret
}
