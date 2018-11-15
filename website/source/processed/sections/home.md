{{ define path "/sections/home.html" }}\
{{ define name "Home" }}\
{{ define title "Magnanimous" }}\
{{ define index 0 }}\
{{ include /processed/_header.html }}\

# Magnanimous

> The best and fastest static website generator in the world!

{{ include /processed/components/_spacer.html }}

<hr />

Magnanimous generates static websites from source files at the speed of light.

And it's incredibly simple to use!

### 1. Download the binary

```
$ curl <url-will-go-here>
```

### 2. Write some sources

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

### 3. Profit

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

You might have noticed some funny contents like `\{{ eval title }}` were _translated_ into something else...

These are [Magnanimous instructions]({{ eval INSTRUCTIONS_PATH }}), which let Magnanimous know when you want 
to do things like include a file into another, or declare values to be used somewhere else... the text inside the
`\{{` and `}}` braces are evaluated using the Magnanimous expression language. But don't worry!
You can learn that in about 5 minutes!

Head to the [Documentation]({{ eval baseURL + "/sections/docs.html" }}) to learn more.

{{ include /processed/_footer.html }}
