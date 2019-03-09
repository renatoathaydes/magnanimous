{{ define path "/sections/about.html" }}
{{ define name "About" }}
{{ define title "About" }}
{{ define index 10 }}
{{ include /processed/_header.html }}
# About Magnanimous

Magnanimous is the fastest and nicest static website generator in the world _(according to me)_.

It was created out of frustration with other existing website generators. They all start simple...
But then, they start adding features until there's so much going on that you need to read a thick book
to fully understand how they work!

I won't let that happen to Magnanimous!

All I really need is a small, fast binary that I can use to stitch together some files and put them in a directory 
my web server can serve without ceremony.

For that to work, a very small set of things need to be available:

* ability to include files in other files.
* ability to customise the contents of included files.
* a converter from MD files to HTML.
* a small expression language to allow necessary automated content (toc, indexes etc.).

So, I set out to build exactly that, and I think the result is pretty great!

Please try it and see for yourself what Magnanimous can do!

{{ include /processed/_footer.html }}