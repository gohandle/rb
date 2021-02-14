package rb

import "net/http"

// AppOptions holds configurable settings for an rb.App
type AppOptions struct {
	noDefaultMiddleware bool
	noDefaultHelpers    bool
	basePath            string
	errorHandler        func(a *App, w http.ResponseWriter, r *http.Request, err error) error
}

// AppOption configures the App
type AppOption func(*AppOptions)

// NoDefaultHelpers disables the adding of default template helpers into
// the template helper
func NoDefaultHelpers() AppOption {
	return func(opts *AppOptions) {
		opts.noDefaultHelpers = true
	}
}

// NoDefaultMiddleware disables the adding of default middleware into the
// app's router. This normally done right when the App is constructed
func NoDefaultMiddleware() AppOption {
	return func(opts *AppOptions) {
		opts.noDefaultMiddleware = true
	}
}

// BasePath configures the app to generate all absolute URLs prefixed with this
// base path. Usefull if the application is served behind a proxy that
// routes all request on a sub-path
func BasePath(p string) AppOption {
	return func(opts *AppOptions) {
		opts.basePath = p
	}
}

// ErrorHandler can be configured to get called when an error occured during rendering
func ErrorHandler(eh func(a *App, w http.ResponseWriter, r *http.Request, err error) error) AppOption {
	return func(opts *AppOptions) {
		opts.errorHandler = eh
	}
}