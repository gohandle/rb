package rb

// URLOptions hold options for url generation
type URLOptions struct {
	Pairs    []string
	BasePath string
}

// URLOption configures how urls are generated
type URLOption func(*URLOptions)

// URLVar configures a url variable if a route needs that to be generated
func URLVar(k, v string) URLOption {
	return func(o *URLOptions) {
		o.Pairs = append(o.Pairs, []string{k, v}...)
	}
}

// BasePath configures a base path that should prefix the url generated url
func BasePath(bp string) URLOption {
	return func(o *URLOptions) {
		o.BasePath = bp
	}
}
