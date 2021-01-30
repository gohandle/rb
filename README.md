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

## TODO
- MUST be able to use a.URL() directly into a.Render with Redirect()
- COULD allow default template options that are applied on each render