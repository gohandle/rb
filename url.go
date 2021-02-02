package rb

import (
	"fmt"

	"go.uber.org/zap"
)

type urlOptions struct {
	pairs []string
	port  int
}

// URLOption configures how urls are generated
type URLOption func(*urlOptions)

// URLVar configures a url variable if a route needs that to be generated
func URLVar(k, v string) URLOption {
	return func(o *urlOptions) {
		o.pairs = append(o.pairs, []string{k, v}...)
	}
}

// GenerateURL will generate a url by the name of a route. If this fails an error is returned.
func (a *App) GenerateURL(s string, opts ...URLOption) (string, error) {
	var o urlOptions
	for _, opt := range opts {
		opt(&o)
	}

	if s[0] == '/' {
		return s, nil
	}

	r := a.mux.Get(s)
	if r == nil {
		return "", fmt.Errorf("no route with name '%s'", s)
	}

	loc, err := r.URL(o.pairs...)
	if err != nil {
		return "", fmt.Errorf("failed to generate url from route: %w, pairs: %v", err, o.pairs)
	}

	return loc.String(), nil
}

// URL generates a URL, it calls GenerateURL but only logs any errors that occure and returns an
// empty string instead.
func (a *App) URL(s string, opts ...URLOption) string {
	s, err := a.GenerateURL(s, opts...)
	if err != nil {
		a.L().Error("failed to generate url",
			zap.String("s", s),
			zap.Error(err))
		return ""
	}
	return s
}
