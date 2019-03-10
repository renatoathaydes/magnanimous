{{ define path "/sections/docs/basic_tutorial.html" }}
{{ define chapter 6 }}
{{ define title "Basic Tutorial" }}
{{ include _docs_header.html }}

# 6. Basic Tutorial

This is the official Magnanimous Tutorial. The intention of this tutorial is to make you capable of
using Magnanimous to create awesome website in the least amount of time possible, so let's get right
to the action!

## Part 1 - Website header, footer and home page

In Part 1, we'll set up the basic parts of the website. Things like the header and footer most pages will share.

We'll also start creating the website's Home Page and get a HTTP server to serve it locally.

### 1.1 Create the website's skeleton

Before we can start building actual content, we need to build some scaffolding. This is the stuff that
will probably be present in all of your pages, including the main CSS styles, header and footer etc.

First of all, we need a root directory for our project:

```
$ mkdir my-website
$ cd my-website
```

Magnanimous expects all of your content to be put into the `source` directory.

Within that, there should be a `processed` directory for files that Magnanimous should process, and a `static`
directory for files which should simply be copied as they are.

Files that go into the `processed` directory can use [Magnanimous instructions](expression_lang.html), which let you
process the contents of a file by including other files into it, adding or removing things to it depending on certain
conditions, define and use variables to avoid repetition, and many other cool things, as we'll see later.

But continuing with the website skeleton: we'll also definitely need some CSS to make the website look great...
let's create a `css` directory under `static` to store the stylesheets:

```
$ mkdir -p source/static/css
$ mkdir -p source/processed
```

To keep with the minimalistic nature of Magnanimous, let's use a minimalistic CSS framework,
[Milligram](https://milligram.io), which will make it easy for us to create a good looking website at the cost of
adding a 2kb gzipped CSS file to our pages, which is about as cheap as it can get!

To get the Milligram CSS file, you can either download it from its [CDN](https://cdnjs.com/libraries/milligram) or
install with NPM:

```
$ npm install milligram
$ cp node_modules/milligram/dist/milligram.min.css source/static/css
```

Great! We've got our first website file (and we didn't have to write anything!)... now we need some HTML fragments
(incomplete pieces of HTML that will be used by other files to create complete HTML files).

Create a HTML header at `source/processed/_header.html` with the following contents:

```html
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Magnanimous Tutorial Website</title>
    <link rel="stylesheet" href="/css/milligram.min.css">
</head>
<body>
<header id="home" class="header">
<section class="container">
    <h1>Demo Website</h1>
    <p class="description">A Magnanimous demo website</p>
</section>
</header>
```

{{ component /processed/components/_file-box.html }}\
    {{ define file "source/processed/_header.html" }}
{{ end }}

Now, create a footer:

```html
</body>
<footer class="footer">
    <section class="container">
        <p>This is the website footer</p>
        <p>Check out the <a href="https://renatoathaydes.github.io/magnanimous">Magnanimous</a> website.</p>
    </section>
</footer>
</html>
```

{{ component /processed/components/_file-box.html }}\
    {{ define file "source/processed/_footer.html" }}
{{ end }}


This is the last piece we needed to create full HTML pages! Notice how the header and footer fit together, so we
should `include` them both in all our HTML and Markdown pages.

### 1.2. Create the Home page (index.html)

We're ready to create the `index.html` page now!

```html
\{{ include _header.html }}
<section class="container">
    <h3>Creating great Websites</h3>
    <p>Isn't this a good looking page?!</p>
</section>
\{{ include _footer.html }}
```

{{ component /processed/components/_file-box.html }}\
    {{ define file "source/processed/index.html" }}
{{ end }}

Done!

### 1.3. Build the website

Just run Magnanimous:

```
$ magnanimous
```

It will print something like this:

```
2019/03/09 20:22:14 No global context file defined.
2019/03/09 20:22:14 Creating file target/index.html from source/processed/index.html
2019/03/09 20:22:14 Creating file target/css/milligram.min.css from source/static/css/milligram.min.css
2019/03/09 20:22:14 Magnanimous generated website in 1.766838ms
```

> Don't worry about the `No global context file defined` message for now... we'll see later how to use a global
  context to support internationalization and deploying the website to different server root paths.

Check out the `target` directory! Your website is there, ready to be served.

Here's what it looks like:

<iframe src="demo/demo-website-part-1.html"
    title="Demo Website V1" width="90%" height="360"></iframe>

### 1.4. Serve the website locally

To serve the website locally during development, we'll need a simple HTTP server.

There's [a lot of choices](https://gist.github.com/willurd/5720255) available!

Here are a few of them, choose your favourite language:

> all of the commands below will serve the `target/` directory on port 8080.
  These web servers are recommended for debugging only. See the [Publishing your website](#deploying) section for
  more robust hosting alternatives.

#### Python HTTP Server

```
$ cd target && python -m SimpleHTTPServer 8080
```

Docs: [https://docs.python.org/2/library/simplehttpserver.html](https://docs.python.org/2/library/simplehttpserver.html)

#### NodeJS HTTP Server

```
$ npm install -g http-server
$ http-server target
```

Docs: [https://www.npmjs.com/package/http-server](https://www.npmjs.com/package/http-server)

#### Java HTTP Server

> Disclaimer: I'm the author of RawHTTP.

```
$ curl https://jcenter.bintray.com/com/athaydes/rawhttp/rawhttp-cli/1.0/rawhttp-cli-1.0-all.jar -o rawhttp.jar
$ java -jar ./rawhttp.jar serve .
```

Docs: [https://renatoathaydes.github.io/rawhttp/rawhttp-modules/cli/](https://renatoathaydes.github.io/rawhttp/rawhttp-modules/cli/)

#### Go HTTP Server

```
$ go get github.com/vwochnik/gost
$ gost target
```

Docs: [https://github.com/vwochnik/gost](https://github.com/vwochnik/gost)

## Part 2 - Solving common problems

In Part 2, we'll explore common features of Magnanimous and learn how we can use them to solve the most common problems
we're likely to face when building static websites.

### 2.1. Adding Navigation

Most websites require some form of navigation tools, like menu bars or tabs, to let users get around the website.

This implies the website has different pages the user can visit! Let's create a few dummy pages so we can demonstrate
how Magnanimous can help.

By the way, this is a good time to switch from HTML to [Markdown](https://www.markdownguide.org/) for content. It's
much easier to write.

> There's a [Markdown Guide](markdown_guide.html) in the Magnanimous Docs, check it out if you want to know more about
  Markdown support in Magnanimous.

Create the following files under the `source/processed/sections` directory:

> The directory maps directly to the page URL, so files inside `source/processed/sections/` will be served under the
  `sections/` URL path (e.g. `http://www.example.com/sections/home.html`).

```markdown
\{{ include /processed/_header.html }}
## About

This website demonstrates some [Magnanimous](https://renatoathaydes.github.io/magnanimous/) features.

\{{ include /processed/_footer.html }}
```

{{ component /processed/components/_file-box.html }}\
    {{ define file "source/processed/sections/about.md" }}
{{ end }}

{{include /processed/components/_spacer.html }}\

```markdown
\{{ include /processed/_header.html }}
## Contact Us

Do not hesitate to contact us!

[Click here](https://renatoathaydes.github.io/magnanimous/) to visit our website.

\{{ include /processed/_footer.html }}
```

{{ component /processed/components/_file-box.html }}\
    {{ define file "source/processed/sections/contact.md" }}
{{ end }}

After running `magnanimous` again, you can visit the new pages directly on their URL:

* [http://localhost:8080/sections/about.html](http://localhost:8080/sections/about.html).
* [http://localhost:8080/sections/contact.html](http://localhost:8080/sections/contact.html).

> Notice that markdown pages are automatically converted to HTML by Magnanimous (browsers can only show HTML file).

But users wouldn't be able to find out about these pages unless you linked to them from the home page.

Without Magnanimous, you'd need to manually maintain links to the sections in the home page. But Magnanimous makes
it much easier: the [for](expression_lang.html#for) instruction lets you generate links to all pages in a particular
directory automatically.

The `for` instruction we need to add to the header looks like this:

```html
\{{ for section /processed/sections/ }}\\
    <a href="\{{ eval section }}" class="button">\{{ eval section }}</a>
\{{ end }}\\
```

> The `\` at the end of the instruction lines escapes the new line in the resulting file, so there won't be lots of
  empty lines in it.

The `for section /processed/sections/` line opens a _for loop_ over the files under the `/processed/sections` directory.
Each file is then assigned to a variable named `section`.

Next, we create HTML links by calling `eval section` both in the `href` and the actual text shown to the users. We'll
improve this later... but for now, here's what the `processed/_header.html` file should look like at this point:

```html
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Magnanimous Tutorial Website</title>
    <link rel="stylesheet" href="/css/milligram.min.css">
</head>
<body>
<header id="home" class="header">
<nav class="navigation">
    <section class="container">
        <a href="/index.html" class="button">Home</a>
        \{{ for section /processed/sections/ }}\\
        <a href="\{{ eval section }}" class="button">\{{ eval section }}</a>
        \{{ end }}\\
    </section>
</nav>
<section class="container">
    <h1>Demo Website</h1>
    <p class="description">A Magnanimous demo website</p>
</section>
</header>
```

{{ component /processed/components/_file-box.html }}\
    {{ define file "source/processed/_header.html" }}
{{ end }}

{{include /processed/components/_spacer.html }}\

With the new page header, the new `About` page now looks like this:

<iframe src="demo/demo-website-part-2.html"
    title="Demo Website V1" width="90%" height="360"></iframe>


### 2.2 Adding Content Summary

You may have noticed that the section links on the top of the website show the path to each section file, rather than
the name of the sections. That's because Magnanimous doesn't know what the user-friendly name of each file should be.

But we can tell it! On each section page, we can define a variable called `name` which can then be used by the header
to display the section name, rather than its path:

```markdown
\{{ define name "About" }}

## About

...
```

Now, the header can refer to the `name` variable on each file (notice we write `eval section.name` instead of just
`eval section` for the link's text):

```html
\{{ for section /processed/sections/ }}\\
    <a href="\{{ eval section }}" class="button">\{{ eval section.name }}</a>
\{{ end }}\\
```

This feature allows us to create summaries of files very easily... for example, you could put your blog posts under the
`processed/blog` directory, and add a few definitions on top of each file, like this:

```markdown
\{{ define name "My first blog post" }}\\
\{{ define date "2019-03-10" }}\\
\{{ define abstract "A short discussion about what I want to write about on my blog" }}\\

## My first blog post

### Abstract

\{{ eval abstract }}\\

...
```

On the blog summary page, you could then list your blog posts:

```html
<h2>Most recent blog posts:</h2>
<ul> 
\{{ for blog_post /processed/blog }}\\
    <li>
        <h3><a href="\{{ eval blog_post }}">\{{ eval blog_post.name }}</a></h3>
        <div>\{{ eval blog_post.date  }}</div>
        <div>\{{ eval blog_post.abstract }}</div>
    </li>    
\{{ end }}\\
</ul>
```

The [`for`](expression_lang.html#for) instructions supports sorting and limits:

```html
\{{ for blog_post (sortBy date reverse limit 10) }}\\
...
\{{ end }}
```

### 2.3 Customizing included content 

The navigation bar on our current website now correctly displays the names of the sections and links to the right
paths. However, it does not display the button for the current page differently from the other buttons, which can be
disorientating for users.

To fix that, we need a way to customize the header, so that it applies some extra style on the button for the current
page.

To do that, we'll use the [`if`](expression_lang.html#if) instruction, which works like this:

```html
<a href="/index.html" 
   class="button\{{if name != `Home`}} button-outline\{{end}}">Home</a>
```

> Notice the use of <code>`</code> (back-ticks) inside the class String to avoid confusion with the surrounding
  double-quotes. Magnanimous allows both back-ticks and double-quotes as String de-markers.

In plain english: If the current page has defined a `name` variable with a value different from `Home`, the 
`button-outline` will be applied to the button, otherwise, it won't.

We haven't defined the name of the Home page yet... to make this work, we need to add this line to the top of the
`index.html` file (it must be defined before the header is included):

```html
\{{ define name "Home" }}
```

Next, we add similar `if`s for the section pages (which already have names!), so that now the HTML header will 
look like this:

```html
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Magnanimous Tutorial Website</title>
    <link rel="stylesheet" href="/css/milligram.min.css">
</head>
<body>
<header id="home" class="header">
<nav class="navigation">
    <section class="container">
        <a href="/index.html" 
            class="button\{{if name != `Home`}} button-outline{{end}}">Home</a>
        \{{ for section /processed/sections/ }}\\
        <a href="\{{ eval section }}" 
            class="button\{{if name != section.name}} button-outline\{{end}}">\{{ eval section.name }}</a>
        \{{ end }}\\
    </section>
</nav>
<section class="container">
    <h1>Demo Website</h1>
    <p class="description">A Magnanimous demo website</p>
</section>
</header>
```

The navigation buttons will now look highlighted only on the page they refer to:

<iframe src="demo/demo-website-part-2-3-home.html"
    title="Demo Website V1" width="90%" height="250"></iframe>

### 2.4 The global context

When you run `magnanimous`, you'll see the following message:

```
2019/03/10 18:49:34 No global context file defined.
```

This means you have not yet defined a global context file. The global context file is just a processed file which
you can use to define global variables that can be used anywhere.

By default, the global context file is located at `source/processed/_global_context`, but you can set it to another
file by using the `-globalctx` option (the path must be relative to `source/processed`:

```
$ magnanimous -globalctx _local_global_context
```

This may seem like a bad idea, but it's extremely useful to solve problems like internationalization and making links
that work on different environments. In the next couple of sections we'll see how we can do that.

#### 2.4.1 Supporting multiple base paths

When your website is not going to be served from the root path `/`, you will run into trouble with absolute links.

That's because absolute links will point to the server's root path, but your website will not be there!

For example, if you deploy your website to GitHub Pages, the URL to your website will look like this:

```
https://username.github.io/projectname/
```

> This website is at [https://renatoathaydes.github.io/magnanimous/](https://renatoathaydes.github.io/magnanimous/).

Notice that the root path is `/projectname`, so absolute links to a page under `source/processed/example/file.html`
need to look like `/projectname/example/file.html`, not just `/example/file.html` as you in your local machine!

To fix this, you can use the global context to define the `baseURL` of your website, which may be different depending
on where you deploy it to.

Let's say we want to be able to both deploy the website locally, as we've been doing, and to 
[GitHub pages](https://pages.github.com/), and that our
project name on GitHub is `magnanimous-tutorial`.

Then, we could define the following 2 global context files:

```
Global context for local deployment: the website will be under the web server's root path
\{{ define baseURL "" }}
```

{{ component /processed/components/_file-box.html }}\
    {{ define file "source/processed/_global_context" }}
{{ end }}

And:

```
GitHub Pages context file
\{{ define baseURL "/magnanimous-tutorial" }}
```

{{ component /processed/components/_file-box.html }}\
    {{ define file "source/processed/_github_global_context" }}
{{ end }}

Now, we need to change every absolute link in the website to take into consideration the `baseURL` variable.

There's a whole chapter on [Paths and Links](paths.md) as it's easy to confuse the two. But suffice to say that links
are things you want the browser to resolve, and paths are things you want Magnanimous to resolve at "compile" time.

Given this knowledge, things like `include` instructions will work fine as they are, but things like links to the
stylesheet will not.

Every link we've added to the demo website so far is in the header page, so it's pretty easy to change them as follows:

* stylesheet import:

```html
<link rel="stylesheet" href="\{{eval baseURL}}/css/milligram.min.css">
```

* section links:

```html
<a href="\{{eval baseURL}}/index.html" 
    class="button\{{if name != `Home`}} button-outline\{{end}}">Home</a>
\{{ for section /processed/sections/ }}\\
<a href="\{{ eval baseURL + section }}" 
    class="button\{{if name != section.name}} button-outline\{{end}}">\{{ eval section.name }}</a>
{{ end }}\\
```

And that's it!

Now, to deploy locally you just do as you've been doing so far:

```
$ magnanimous
```

But when you want to deploy the website to GitHub, let Magnanimous know you want to use the `_github_global_context`
file as your global context:

```
$ magnanimous -globalctx _github_global_context
```

Reading the output, you can make sure Magnanimous picked up the right file:

```
2019/03/10 20:23:44 Using global context file: source/processed/_github_global_context
2019/03/10 20:23:44 Creating file target/sections/about.html from source/processed/sections/about.md
2019/03/10 20:23:44 Creating file target/sections/contact.html from source/processed/sections/contact.md
2019/03/10 20:23:44 Creating file target/css/milligram.min.css from source/static/css/milligram.min.css
2019/03/10 20:23:44 Creating file target/index.html from source/processed/index.html
2019/03/10 20:23:44 Magnanimous generated website in 5.581399ms
```

One last thing: GitHub Pages lets you use the master branch as the deployment branch, but you want only the files in the
`target` directory to be deployed... luckily they also let you configure a folder as the root of your website, but the
directory must be called `docs` for some reason!

So just make sure to move the files there when you want to deploy:

```
$ mv target docs
``` 

And then `git push`! Your website should now be available online!

[Here's the one](https://renatoathaydes.github.io/magnanimous-tutorial/) we've created in this Tutorial!

#### 2.4.2 Internationalization

The way Magnanimous supports internationalization is via the global context.

For pages that you want to re-use across different languages, you can define a set of messages in the different 
global contexts for each language, then refer to those messages rather than hardcode the messages in the pages.

### 2.5 Creating re-usable components


{{ include _docs_footer.html }}
