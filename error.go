package rb

import "net/http"

// ErrorHandleFunc can be casted to to implement the ErrorCore
type ErrorHandleFunc func(http.ResponseWriter, *http.Request, error)

// HandleError handles errors for the core by calling itself
func (f ErrorHandleFunc) HandleError(w http.ResponseWriter, r *http.Request, err error) {
	f(w, r, err)
}

// BasicErrorHandler will simply write the error to the response and set the status code
// to an internal server code.
var BasicErrorHandler = ErrorHandleFunc(func(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
})
