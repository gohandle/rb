## Backlog
- [ ] Think about "base_page" functionality, making it ergonomic without inject
- [ ] Think about translate
- [ ] Redirect should support a option that allows for redirecting to a route. Since
      url should allow returning error
- [ ] Have customizable error handling for server errors instead of http.Error in
      app.Action
- [ ] Can we abstract the mux router by including functionality for getting url vars and
      current route name to our context
- [ ] Have a special render type that arrows streaming responses that also supports
      custom headers as a replacement for not having access to the response in the
      the context
- [ ] How about setting up dafault middleware and helpers
- [ ] Support setting csrf values in the the template
- [ ] SHOULD add erroring NoValueDecoder, NoTemplates, NoValidation implementations for 
      respective cores.
- [ ] SHOULD add more translate options for default message etc

## Testing Backlog
- [ ] Test different form sources
