# rb
A web framework designed for modern server-side rendered applications

Focus: this framework will probably do too much for you if you're just designing APIs. This 
repos focusses on application that render HTML on the server.

OnShoulder: this framework brings together well-known and active packages from the cummunity 
and doesn't try to re-invent the wheel.

Batteries: but through dependency injection, all these batteries can be swapped. As interfaces 
are provided

Testing: comes with html testing utilities that makes asserting redered content easy and fun

Just an advanced implementation of the "Respond" pattern: https://youtu.be/rWBSMsLG8po?t=1555
And advanced implementation of "decode" pattern: https://youtu.be/rWBSMsLG8po?t=1591

Thesis: what if every handler get's created by fx, with its dependencies?

https://www.veritone.com/blog/how-i-write-go-http-services-after-seven-years/

## TODO Sunday Release
- [ ] helper for url
- error pages for rendering
- structural logging for each part (with request scoped logging)
- rbtest for easy assertions
- field based errors (settle on validation first)
- helper for i18n
- helper for session data
- integrated CSRF 

## Ideas
- Come with I18n
- Come with flash messages
- Come with template render
- Come with form binding
- Come with template helpers
- Come with live reload of loaded templates
- Come with pretty error pages in development
- Come with a neat logging window in development
- Monitor and reload when static files change
- Must be able to retrieve Params from a request

## i18n
- Jet translator interface: https://pkg.go.dev/github.com/CloudyKit/jet#Translator
- Validator translator: https://github.com/go-playground/validator/blob/master/_examples/translations/main.go

## TODO
- MUST be able to use a.URL() directly into a.Render with Redirect()
- COULD allow default template options that are applied on each render
- COULD implement multi-render and multi-bind that does it based on content types and accept headers
- SHOULD add options to FormBind that selects which part of the request the bind will perform
- COULD add an empty bind that just calls parse Form
- COULD make the validator optional to provide, but it should error cleary if bind is called with
        the validation option while non is available
- COULD allow passing validation options, such as "filter", "except" and allow use of the "var" validation
- SHOULD add session options to configure how the session is saved
- SHOULD add an option to make all generated urls on top of a basepath
