{{ define path "/sections/home.html" }}
{{ define name "Home" }}
{{ define title "Magnanimous" }}
{{ define index 0 }}
{{ include /processed/_header.html }}
# Magnanimous

> The best and fastest static website generator in the world!

Magnanimous generates static websites from source files at the speed of light.

And it's incredibly simple to use!

### 1. Download the binary

```
$ curl https://getit
```

### 2. Write some sources

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

**source/_footer.html**

```html
</body>
</html>
```

**source/processed/index.md**

```md
\{{ define title "My Website" }}
\{{ include /processed/_header.html }}
# This is my website

How awesome is it?!

\{{ include /_footer.html }}
```

### 3. Profit

```
$ magnanimous
```

Your website is ready on the `target` directory!

**target/index.html**
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

You might have noticed some funny contents like `\{{ }}` were _translated_ into something else...

These are Magnanimous instructions, which let Magnanimous know when you want to do things like
include a file into another, or declare values to be used somewhere else... the text inside the
`\{{` and `}}` braces are evaluated using the Magnanimous expression language. But don't worry!
You can learn this little expression language in about 5 minutes!

{{ include /processed/_footer.html }}
