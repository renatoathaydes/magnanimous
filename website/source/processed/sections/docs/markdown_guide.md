{{ define path "/sections/docs/markdown_guide.html" }}
{{ define chapter 5 }}
{{ define title "Markdown Guide" }}
{{ include _docs_header.html }}

# 5. Markdown Guide

Magnanimous converts all markdown content (files in the `source/processed` directory and with the `.md` extension)
into HTML, changing the resulting file's extension from `.md` to `.html`.

In this guide, you'll learn details about how this conversion is performed and how you can use markdown
to make writing your website easier.

## Technologies used

The conversion markdown-HTML is done via a [Go](https://golang.org/) library called
[Blackfriday](https://github.com/russross/blackfriday).

Source code within markdown (content found between 3 back-ticks) is color-highlighted via 
[Chroma](https://github.com/alecthomas/chroma), which supports a large number of
[languages](https://github.com/alecthomas/chroma#supported-languages).

## Why write markdown

Markdown is much easier to write by hand than HTML, especially if the content being created is mostly composed of
text, images and a simple layout.

For example, take the following HTML content:

```html
<h1>Title</h1>

<p>This is a <a href="https://example.org">link</a>.</p>

<img src="image.png" alt="some image"/>
```

The equivalent content could be written using markdown as follows:

```markdown
# Title

This is a [link](https://example.org).

![some image](image.png)
```

And in case you want to control exactly the layout of the content, you can just embed HTML into markdown!

```markdown
## Example with HTML in Markdown

<div class="float-top-right">HTML can be embedded inside markdown!</div>
```

To learn more about markdown, check out the [GitHub Guide](https://guides.github.com/features/mastering-markdown/)
which briefly describes all you'll need to know!

## Writing full HTML pages via markdown

Magnanimous, in order to stay simple and easy to learn, does not do anything very "magical"!

When you put a markdown file into `source/processed/`, all Magnanimous will do is:

* process all Magnanimous instructions in the file.
* convert the result from markdown to HTML.
* save the result in a `.html` file, not `.md`.

It will not, for example, wrap the contents of the page into a fully formed HTML document automatically.

So, if you write this simple markdown:

**source/processed/index.md**

```markdown
## hello
```

You'll get something like this:

**target/index.html**

```html
<h2>hello</h2>
```

That is not going to show properly in a browser as it's not a full HTML document.

To create a full document, you need to _wrap_ the markdown inside a HTML header and a HTML footer.

Here's an example header:

**source/processed/_header.html**

```html
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>\{{ eval title }}</title>
</head>
<body>
```

And the matching footer:

**source/_footer.html**

```html
</body>
</html>
```

Now, we can `include` these two in the beginning and end of all our markdown files:

**source/processed/index.md**

```markdown
\{{ define title "My Website" }}
\{{ include /processed/_header.html }}
## hello
\{{ include /_footer.html }}
```

Notice that we can define some properties before we include the header in order to customize the header.
Check the [Components](components.html) chapter to learn more about this pattern.

Finally, running `magnanimous` should result in a valid `index.html` file in the `target/` directory:

**target/index.html**

```html
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>My Website</title>
</head>
<body>

<h2>hello</h2>

</body>
</html>
```

## Including source code

Markdown files may contain sample source code which gets automatically highlighted.

To use that, wrap the code sample within three back-ticks, with the name of the language right after the first ticks,
as in this example:

````markdown
```java
class Main {
  public static void main(String[] args) {
    System.out.println("Hello world");
  }
}
```
````

This website is itself created with Magnanimous, so the above formatting is actually a sample of what code samples
should look like in your own website!

> Notice that [this file](https://github.com/renatoathaydes/magnanimous/blob/master/website/source/processed/sections/docs/markdown_guide.md)
  is written in markdown!

The code highlighting is done by [Chroma](https://github.com/alecthomas/chroma).

You can find which [languages](https://github.com/alecthomas/chroma#supported-languages) are supported in their
documentation.                   


{{ include _docs_footer.html }}
