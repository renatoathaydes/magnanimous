{{ define path "/sections/about.html" }}
{{ define name "About" }}
{{ define title "About" }}
{{ define index 10 }}
{{ include /processed/_header.html }}
# About Magnanimous

Magnanimous is the simplest and fastest static website generator in the world.

It was created out of frustration with other existing website generators. They all start simple...
But then, they start adding features until there's so much going on that you need to read a thick book
to fully understand how they work!

All I wanted was a small, fast binary that could just stitch together some files and put them in a directory,
which I could then deploy to my web server without ceremony.

For that to work, a very small set of things need to be available:

* ability to include files in other files.
* ability to customise the contents of included files.
* a converter from MD files to HTML.
* a small expression language to allow necessary automated content (toc, indexes etc.).

So, I set out to build exactly that, and I think the result is pretty great!

Please try it and see for yourself what Magnanimous can do!

{{ include /processed/_footer.html }}