## Backlog
- [ ] BUG: Calling sessions save multiple times causes multiple cookies with the same name being written
      to the Set-Cookie header
- [ ] Think about "base_page" functionality, making it ergonomic without inject
- [x] Think about translate
- [x] Have customizable error handling for server errors instead of http.Error in
      app.Action
- [x] Can we abstract the mux router by including functionality for getting url vars and
      current route name to our context
- [ ] MUST Have a special render type that arrows streaming responses that also supports
      custom headers as a replacement for not having access to the response in the
      the context
- [ ] MUST How about setting up dafault middleware and helpers, helpers now need access to the ctx.
- [ ] MUST Support setting csrf values in the the template
- [ ] SHOULD add erroring NoValueDecoder, NoTemplates, NoValidation implementations for 
      respective cores.
- [ ] SHOULD add more translate options for default message etc
- [ ] SHOULD add build-in support for setting a language in the session
- [ ] COULD provide flash functionality, build int
- [ ] SHOULD include default middleware in NewApp
- [ ] MUST re-add csrf using existing session tools
- [ ] MUST add helpers for: url(), t(), field_error(). Depending on the respecitve Core 
      is available in the app.
- [ ] COULD add an option to make all generated urls on top of a basepath
- [ ] MUST extend csrf middleware with origin checking: https://github.com/gorilla/csrf/blob/master/csrf.go#L249
- [ ] MUST provide more indepth coverage of csrf package
- [ ] MUST Test different form sources
- [ ] COULD Test the various jet helpers through using the app/core

## Flash 
- [ ] Add middleware that pops from the session core and make it available to the context