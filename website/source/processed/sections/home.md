{{ define path "/sections/home.html" }}\
{{ define name "Home" }}\
{{ define title "Magnanimous" }}\
{{ define index 0 }}\
{{ include /processed/_header.html }}\

<div id="top-mag-title"></div>

# Magnanimous

<div id="bottom-mag-title"></div>

> The friendliest website generator in the world!

{{ include /processed/components/_spacer.html }}

<hr />

Magnanimous generates static websites from source files really fast.

And it's incredibly simple to use!

### Features:

{{component /processed/components/_tiles.html}}
{{slot a}}
Simple templating mechanism based on a tiny [expression language]({{eval INSTRUCTIONS_PATH}}),

```javascript
2 + 2 is \{{ eval 2 + 2 }}
```
{{end}}
{{slot b}}
Compose snippets to create content.

```html
\{{ include _header.html }}
<h2>Content!</h2>
\{{ include footer.html }}
```
{{end}}
{{slot c}}
Define variables for easy re-use.

```javascript
\{{ define title "Hello world" }}
```
{{end}}
{{slot d}}
[Components]({{eval baseURL + "/sections/docs/components.html"}}) and slots inspired by web components.

```javascript
\{{ component other-file.html }}
  \{{ slot main }}Main content\{{ end }}
\{{ end }}
```
{{end}}
{{slot e}}
[Conditional content]({{eval INSTRUCTIONS_PATH + "#if"}}).

```javascript
\{{ if title == "Home" }}We are home\{{ end }}
```
{{end}}
{{slot f}}
[For loops]({{eval INSTRUCTIONS_PATH + "#for"}}) over files and variables.

```javascript
\{{ for file /path/to/files }}
* Post name: \{{eval file.postName}}
\{{end}}
```
{{end}}
{{slot g}}
[markdown content]({{eval baseURL + "/sections/docs/markdown_guide.html"}}) automatically converted to HTML.

```markdown
## Markdown is easy!

> Note
```
{{end}}
{{slot h}}
[Source code highlighting]({{eval baseURL + "/sections/docs/markdown_guide.html#source-code"}}) in markdown.

````markdown
```javascript
const hello = () => "Hello world";
```
````
{{end}}
{{define tiles [a, b, c, d, e, f, g, h]}}
{{end}}

### 1. Download the binary

Go to the [releases page](https://github.com/renatoathaydes/magnanimous/releases) and download the binary for your OS.

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

Add some [markdown](https://en.wikipedia.org/wiki/Markdown) for comfortably writing content. 

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

You might have noticed how some funny contents like `\{{ eval title }}`, in the `title` tag, and 
`\{{ define title "My Website" }}`, in the markdown file...

These are [Magnanimous instructions]({{ eval INSTRUCTIONS_PATH }}), which let Magnanimous know when you want it to **process** files,
which means doing things like including a file into another, or declaring values to be used elsewhere... the text inside the
`\{{` and `}}` braces are evaluated using the Magnanimous expression language. But don't worry!
You can learn that in about 15 minutes!

Also notice that all markdown content is converted automatically to HTML! Magnanimous will even highlight source code blocks in pretty much
[any language](https://github.com/alecthomas/chroma#supported-languages), just provide the block language as you do on GitHub markdown:

````markdown
```javascript
const magnanimous = () => "Awesome";
```
````

Head to the [Documentation]({{ eval baseURL + "/sections/docs.html" }}) to learn more.

{{ include /processed/_footer.html }}
