package rb

import (
	"net/http"

	"go.uber.org/zap"
)

type redirectRender struct{ loc string }

// Redirect creates a render that will write a redirect response. It will redirect to the
// provided url.
func Redirect(loc string) Render { return redirectRender{loc} }

func (rr redirectRender) RenderHeader(
	rc RenderCore,
	w http.ResponseWriter,
	r *http.Request,
	status int,
) (int, error) {
	if status < 1 { // no explicit status code set
		status = http.StatusFound
		L(r).Debug("no explit status for redirect render, use default",
			zap.Int("status_code", status))
	}

	w.Header().Set("Location", rr.loc)
	L(r).Debug("redirect wrote location header", zap.String("location", rr.loc))
	return status, nil
}

func (rr redirectRender) Value() interface{} { return nil }
func (rr redirectRender) Render(RenderCore, http.ResponseWriter, *http.Request) error {
	return nil
}
