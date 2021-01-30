package rb

type urlOptions struct {
	pairs []string
	port  int
}

type URLOption func(*urlOptions)

func URLVar(k, v string) URLOption {
	return func(o *urlOptions) {
		o.pairs = append(o.pairs, []string{k, v}...)
	}
}

// URL generates a URL
func (a *App) URL(s string, opts ...URLOption) string {
	var o urlOptions
	for _, opt := range opts {
		opt(&o)
	}

	if s[0] == '/' {
		return s
	}

	r := a.mux.Get(s)
	if r == nil {
		// @TODO log an error
		return ""
	}

	loc, err := r.URL(o.pairs...)
	if err != nil {
		// @TODO log an error
		return ""
	}

	return loc.String()
}
