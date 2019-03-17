{{ define path "/sections/home.html" }}\
{{ define name "Home" }}\
{{ define title "Magnanimous" }}\
{{ define index 0 }}\
{{ include /processed/_header.html }}\

<div id="top-mag-title"></div>

# Magnanimous

<div id="bottom-mag-title"></div>

> The fastest and easiest static website generator in the world!

{{ include /processed/components/_spacer.html }}

<hr />

Magnanimous generates static websites from source files really fast.

And it's incredibly simple to use!

Features:

* simple templating mechanism based on a tiny [expression language]({{eval INSTRUCTIONS_PATH}})
 (simple as `2 + 2 is \{{ eval 2 + 2 }}`).
* include file into another (`\{{ include other-file.html }}`).
* define variables (`\{{ define title "Hello world" }}`).
* [components]({{eval baseURL + "/sections/docs/components.html"}}) and slots inspired by web components.
* [conditional content]({{eval INSTRUCTIONS_PATH + "#if"}}) (`\{{ if title == "Home" }}We are home\{{ end }}`).
* [for loops]({{eval INSTRUCTIONS_PATH + "#for"}}) over files and variables (`\{{ for file /path/to/files }}* Post name: \{{eval file.postName}} \{{end}}`).
* [markdown content]({{eval baseURL + "/sections/docs/markdown_guide.html"}}) automatically converted to HTML.
* [source code highlighting]({{eval baseURL + "/sections/docs/markdown_guide.html#source-code"}}) in markdown.

### 1. Download the binary

```
$ curl <url-will-go-here>
```

### 2. Write some sources

How about some HTML fragments?! 

```html
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>\{{ eval title }}</title>
</head>
<body>
```

{{ component /processed/components/_file-box.html }}\
    {{ define file "source/processed/_header.html" }}
{{ end }}

```html
</body>
</html>
```

{{ component /processed/components/_file-box.html }}\
    {{ define file "source/_footer.html" }}
{{ end }}

Add some [markdown](https://en.wikipedia.org/wiki/Markdown) content for comfort. 

```markdown
\{{ define title "My Website" }}
\{{ include /processed/_header.html }}
# This is my website

How awesome is it?!

\{{ include /_footer.html }}
```

{{ component /processed/components/_file-box.html }}\
    {{ define file "source/processed/index.md" }}
{{ end }}

### 3. Run magnanimous

```
$ magnanimous
```

Your website is ready on the `target` directory!

```html
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>My Website</title>
</head>
<body>

<h1>This is my website</h1>

<p>How awesome is it?!</p>

</body>
</html>
```

{{ component /processed/components/_file-box.html }}\
    {{ define file "target/index.html" }}
{{ end }}

You might have noticed how some funny contents like `\{{ eval title }}` were _translated_ into something else...

These are [Magnanimous instructions]({{ eval INSTRUCTIONS_PATH }}), which let Magnanimous know when you want 
to do things like include a file into another, or declare values to be used somewhere else... the text inside the
`\{{` and `}}` braces are evaluated using the Magnanimous expression language. But don't worry!
You can learn that in about 5 minutes!

Also notice that all markdown content is converted automatically to HTML!

Head to the [Documentation]({{ eval baseURL + "/sections/docs.html" }}) to learn more.

{{ include /processed/_footer.html }}
