## Backlog
- [ ] Think about "base_page" functionality, making it ergonomic without inject
- [x] Think about translate
- [x] Have customizable error handling for server errors instead of http.Error in
      app.Action
- [x] Can we abstract the mux router by including functionality for getting url vars and
      current route name to our context
- [ ] MUST Have a special render type that arrows streaming responses that also supports
      custom headers as a replacement for not having access to the response in the
      the context
- [ ] MUST How about setting up dafault middleware and helpers
- [ ] MUST Support setting csrf values in the the template
- [ ] SHOULD add erroring NoValueDecoder, NoTemplates, NoValidation implementations for 
      respective cores.
- [ ] SHOULD add more translate options for default message etc
- [ ] SHOULD add build-in support for setting a language in the session
- [ ] COULD provide flash functionality, build int

## Testing Backlog
- [ ] Test different form sources
