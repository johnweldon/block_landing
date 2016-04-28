# block_landing

This is a small Go website that just handles squidGuard redirects to display a friendly message when the proxy blocks a url.

It's intended to run in a docker container, but the build artifact binary can also be run directly.

To generate the docker image you can just use the provided Makefile, and the `build` target.

To deploy the image so that it's persistent you can use the `deploy` target in the Makefile.

## Url format

This webserver expects the url format to be:
`http://example.com/?target=%t&uri=%u&clientip=%a`

 * `%t` is the squidGuard replacement code for "target category"; usually the name of the blocked category.
 * `%u` is the originally requested URL.
 * `%a` is the IP address of the client - it's not displayed now, but it could be if you modified the code.


## Visual appearance

As it is now, the visual appearance is specific to my family.
You'd probably want to tweak the CSS and possibly the HTML template to suit your environment.
Unless you just like the picture of my family.
