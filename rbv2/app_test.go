package rb_test

import (
	"net/http"
	"testing"

	rb "github.com/gohandle/rb/rbv2"
)

func HandleFoo(a rb.App) http.Handler {
	type Page struct {
		Err  error
		Form struct {
			Foo string `form:"foo"`
		}
	}

	return a.Action(func(c rb.Ctx) error {
		p, submit := Page{}, false

		if submit, p.Err = c.Bind(rb.Form(&p.Form)); submit && p.Err == nil {
			c.Session().Set("_flash", c.Translate("message.flash.saved"))
			return c.Render(rb.Redirect(c.URL("home")))
		}

		return c.Render(rb.Template("foo.html", p))
	})
}

func TestFooExample(t *testing.T) {

}
