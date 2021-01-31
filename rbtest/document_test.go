package rbtest

import (
	"errors"
	"strings"
	"testing"
)

type errReader struct{}

func (er errReader) Read([]byte) (int, error) { return 0, errors.New("expected") }

func expectFatal(tb testing.TB, expFormat string) func() {
	old := fatalf
	var actFormat string
	fatalf = func(tb testing.TB, format string, args ...interface{}) {
		actFormat = format
	}

	return func() {
		if actFormat != expFormat {
			tb.Fatalf("expected fatalf: '%v' got: '%v'", expFormat, actFormat)
		}
		fatalf = old
	}
}

func TestMustParseDocument(t *testing.T) {
	t.Run("fail", func(t *testing.T) {
		defer expectFatal(t, "rbtest: failed to create document from reader: %v")()
		MustParseHtml(t, errReader{})
	})

	doc := MustParseHtml(t, strings.NewReader(`<p></p>`))
	if doc == nil {
		t.Fatalf("got: %v", doc)
	}

	if act := doc.MustHtml(); act != `<html><head></head><body><p></p></body></html>` {
		t.Fatalf("got: %v", act)
	}
}
