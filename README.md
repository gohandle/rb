## Backlog

- [ ] MUST extend csrf middleware with origin checking: https://github.com/gorilla/csrf/blob/master/csrf.go#L249
- [ ] MUST provide more indepth coverage of csrf package and carefully thinkg about security implications
- [ ] MUST add build-in support for setting a language in the session and detecting in translate core
- [ ] SHOULD add more options for the rbcore/Configured method
- [ ] SHOULD add an option to make all generated urls on top of a basepath
- [ ] SHOULD add more c.Translate() options for default message etc
- [ ] SHOULD Test different form sources
- [ ] COULD configure the jet helpers through a slice instead of arguments 1-by-1 while still supporting fx di
- [ ] COULD add better support for configuring names of cookies, fields, headers, etc
- [ ] COULD add a function to rbtest for creating requests with session date and valid csrf 
- [ ] COULD Think about "base_page" functionality, making it ergonomic without inject
- [ ] COULD Have a special render type that arrows streaming responses that also supports
      custom headers as a replacement for not having access to the response in the
      the context
- [ ] COULD check if our JIT response doesn't incure too much performance overhead
- [ ] COULD Test the various jet helpers through using the app/core
- [ ] COULD add a more ergonomic .Save method directly to the session that allows for defer with error logging
- [ ] COULD move the csrf middleware to a new package
- [ ] COULD add more testing for the rbapp package
- [ ] COULD test the rbtest package a bit better
- [ ] COULD add erroring NoValueDecoder, NoTemplates, NoValidation implementations for 
      respective cores.
- [ ] COULD add a way to make test requests that include certain session values (and flashes)
- [x] MUST provide flash functionality out of the box, see _flashbak dir
- [x] BUG: Calling sessions save multiple times causes multiple cookies with the same name being written
      to the Set-Cookie header
- [x] Think about translate
- [x] Have customizable error handling for server errors instead of http.Error in
      app.Action
- [x] Can we abstract the mux router by including functionality for getting url vars and
      current route name to our context
- [x] MUST How about setting up dafault middleware and helpers, helpers now need access to the ctx.
- [x] MUST Support setting csrf values in the the template
- [x] SHOULD include default middleware in NewApp
- [x] MUST re-add csrf using existing session tools
- [x] MUST add helpers for: url(), t(), field_error(). Depending on the respecitve Core 
      is available in the app.

## OLD README

# rb
Package rb provides a framework for creating server-side rendered web applications

Unlike other frameworks it doesn't prevent you from using regular http.Handler methods
for handling. Such as those created using http.HandlerFunc

## Inspiration
- https://www.veritone.com/blog/how-i-write-go-http-services-after-seven-years/

## Feature TODO
- [x] helper for url
- [x] error pages for rendering
- [x] rbtest for easy assertions
- [x] allow easy flash messages from session: ReadAndDelete
- [x] field based errors (settle on validation first)
- [x] helper for non_field errors
- [x] structural logging for each part (with request scoped logging)
- [x] helper for i18n, especially for translating validation messages
- [x] integrated CSRF 

## Future Ideas
- Come with live reload of loaded templates
- Come with pretty error pages in development
- Come with a neat logging window in development
- Asset helpders with Monitor and reload when they change

## i18n
- Jet translator interface: https://pkg.go.dev/github.com/CloudyKit/jet#Translator
- Validator translator: https://github.com/go-playground/validator/blob/master/_examples/translations/main.go

## TODO
- [x] MUST: rb.Redirect shoul also set the status code to some redirect status if the user doesn't 
            set it since then the redirect won't work at all
- [x] MUST: default error handler fails (probably because rendering is halfway)
- [ ] SHOULD rethink injectors, current implementation was build in a hurry
- [x] MUST be able to use a.URL() directly into a.Render with Redirect()
- [x] SHOULD add options to FormBind that selects which part of the request the bind will perform
- [ ] MUST  allow url generation errors to be scoped to a request, but that means the helpers also
            need acess to the request's context
- [ ] COULD add logging to helpers, but hard to get request scoped logger
- [ ] COULD add an option that disables the default middleware application if the user wants control
            over the order
- [ ] COULD allow default template options that are applied on each render
- [ ] COULD implement multi-render and multi-bind that does it based on content types and accept headers
- [ ] COULD add an empty bind that just calls parse Form
- [ ] COULD make the validator optional to provide, but it should error cleary if bind is called with
            the validation option while non is available
- [ ] COULD allow passing validation options, such as "filter", "except" and allow use of the "var" validation
- [ ] COULD find a way to know for sure if the header was already written and have a flag to indicate to 
            the error handler that  the header was written already
- [ ] COULD add more assertion utilites to our rbtest.Document type, would be cool if we could scope assertion
            such that error messages only show part of the html when it fails
- [ ] COULD add a helper that renders (field) errors that have not been rendered yet, that means keeping
            track of what has been rendered. But can be helpfull for debugging
- [ ] COULD add session options to configure how the session is saved (expires, etc)
- [ ] COULD add an option to make all generated urls on top of a basepath
- [x] COULD have SOME documentation
- [ ] COULD add integration point and allow users to provide their own checks for the field_error helper
            so they can bring their own validator and still filter errors for a certain field
- [ ] COULD add a rendering option that buffers the response so rendering errors can be shows a 
            a completely new page (Possibly with dev options)
