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

## Feature TODO
- [x] helper for url
- [x] error pages for rendering
- [x] rbtest for easy assertions
- [x] allow easy flash messages from session: ReadAndDelete
- [x] field based errors (settle on validation first)
- [x] helper for non_field errors
- [ ] structural logging for each part (with request scoped logging)
- [ ] helper for i18n, especially for translating validation messages
- [ ] integrated CSRF 

## Validation Options
- https://github.com/go-playground/validator (oct 2020, ~7k, big boy)
- https://github.com/thedevsaddam/govalidator (apr 2020, ~1k, laravel inspired)
- https://github.com/asaskevich/govalidator (okt 2020, ~4.6k, inspired by validator.js)
- https://github.com/go-ozzo/ozzo-validation (oct 2020, ~1.9k, no tags)
- github.com/gookit/validate (jan 2021, ~400, good track record of releases)

## Future Ideas
- Come with live reload of loaded templates
- Come with pretty error pages in development
- Come with a neat logging window in development
- Asset helpders with Monitor and reload when they change

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
- COULD find a way to know for sure if the header was already written and have a flag to indicate to 
        the error handler that  the header was written already
- COULD add more assertion utilites to our rbtest.Document type, would be cool if we could scope assertion
        such that error messages only show part of the html when it fails
- COULD add a helper that renders (field) errors that have not been rendered yet, that means keeping
        track of what has been rendered. But can be helpfull for debugging
- SHOULD add session options to configure how the session is saved
- SHOULD add an option to make all generated urls on top of a basepath
- SHOULD have SOME documentation
- COULD add integration point and allow users to provide their own checks for the field_error helper
        so they can bring their own validator.