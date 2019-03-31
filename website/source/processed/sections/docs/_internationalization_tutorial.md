{{ define path "/sections/docs/internationalization_tutorial.html" }}
{{ define chapter 7 }}
{{ define title "Advanced Tutorial" }}
{{ include _docs_header.html }}

# 7. Advanced Tutorial

This is the Advanced Magnanimous Tutorial. In the [Basic Tutorial](basic_tutorial.html), you've learned how to create
a website's skeleton (header, footer, styles), the home page, how to add navigation, content summary, customize included
files, and use the global context to define global variables (like the base path of your website). You also learned how
to test the website locally, using your favourite HTTP server, and how to publish it to the world using
[GitHub Pages](https://pages.github.com/)!

That's quite a lot already! But there's still a few more tricks you should learn to become proficient using all
features of Magnanimous.

This tutorial will show you techniques you can use to internationalize your website (useful for other things as well,
like allowing users to use different themes) and use Magnanimous [components](components.html) to create re-usable
pieces that will make your website more consistent, and you much more efficient.

## Part 1 - Internationalization

The best way to internationalize your website is to simply have a different set of pages for each language and
re-use only Magnanimous [components](components.html) and a few common pages.

The localized pages should go under a prefix path that is equal to the language prefix.

So, considering the `/index.html` default page (in our case, in English), the equivalent page in Portuguese would
be `/pt/index.html`.

The file layout should look like this:

```

```

For the common pages that you want to re-use across different languages, you can externalize the
messages shown on the header, so that their actual values may come from the context. We'll then be able to _override_
the messages for each different language by simply importing the relevant `_messages` file as follows:

```
\{{ include eval "/processed/" + language + "/_messages" }}
```

That means we need to replace the following:

* page title

From:

```html
<title>Magnanimous Tutorial Website</title>
```

To:

```html
<title>\{{ eval page_title }}</title>
```

* Home section link

From:

```html
<a href="\{{ eval baseURL }}/index.html" 
    class="button\{{if name != `Home`}} button-outline\{{end}}">Home</a>
```

To:

```html
<a href="\{{ eval baseURL + `/` + language }}/index.html" 
    class="button\{{ if name != home }} button-outline\{{end}}">\{{ eval home }}</a>
```

> _notice that we'll have to define the `home` variable in each `<language>/index.html`_

* Other section links

From:

```html
\{{ for section /processed/sections/ }}\\
<a href="\{{ eval baseURL + section }}" 
    class="button\{{if name != section.name}} button-outline\{{end}}">\{{ eval section.name }}</a>
\{{ end }}
```

To:

```html
\{{ for section eval "/processed/sections/" + language }}\\
<a href="\{{ eval baseURL + section }}" 
    class="button\{{if name != section.name}} button-outline\{{end}}">\{{ eval section.name }}</a>
\{{ end }}
```

* Header main element:

From:

```html
<section class="container">
    <h1>Demo Website</h1>
    <p class="description">A Magnanimous demo website</p>
</section>
```

To:

```html
<section class="container">
    <h1>\{{ eval main_page_name }}</h1>
    <p class="description">\{{ eval main_page_description }}</p>
</section>
```
