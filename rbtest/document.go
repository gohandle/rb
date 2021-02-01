package rbtest

import (
	"bytes"
	"io"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

var fatalf = func(tb testing.TB, format string, args ...interface{}) {
	tb.Fatalf(format, args...)
}

type Doc struct {
	*goquery.Document
	tb   testing.TB
	read *bytes.Buffer
}

func MustParseHtml(tb testing.TB, r io.Reader) *Doc {
	read := bytes.NewBuffer(nil)
	doc, err := goquery.NewDocumentFromReader(io.TeeReader(r, read))
	if err != nil {
		fatalf(tb, "rbtest: failed to create document from reader: %v", err)
	}

	return &Doc{doc, tb, read}
}

func (d *Doc) ReadHTML() string {
	return d.read.String()
}

func (d *Doc) MustHtml() string {
	ret, err := d.Html()
	if err != nil {
		fatalf(d.tb, "rbtest: failed to convert to HTML: %v", err)
	}

	return ret
}
