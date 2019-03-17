# Magnanimous

![Magnanimous mascot](website/source/static/images/mag-mid-logo.jpg)

[https://renatoathaydes.github.io/magnanimous](https://renatoathaydes.github.io/magnanimous)

The simplest and nicest static website generator in the world!!

Magnanimous has the following guiding principles:

* no configuration.
* no magic.
* explicit is better than implicit.
* generate static websites well, nothing else.

## Features

* simple templating mechanism based on a tiny [expression language](https://renatoathaydes.github.io/magnanimous/sections/docs/expression_lang.html)
 (simple as `2 + 2 is {{ eval 2 + 2 }}`).
* include file into another (`{{ include other-file.html }}`).
* define variables (`{{ define title "Hello world" }}`).
* [components](https://renatoathaydes.github.io/magnanimous/sections/docs/components.html) and slots inspired by web components.
* conditional content (`{{ if title == "Home" }}We are home{{ end }}`).
* for loops over files and variables (`{{ for file /path/to/files }}* Post name: {{eval file.postName}} {{end}}`).
* markdown content automatically converted to HTML.
* [source code highlighting](https://renatoathaydes.github.io/magnanimous/sections/docs/markdown_guide.html#source-code) in markdown.

## Documentation

Visit [the website](https://renatoathaydes.github.io/magnanimous) for the full documentation.

## Basic Design

Magnanimous reads source files from the `source/` directory, and generates a website in `target/`.

The following `source` sub-directories are the only ones treated _specially_ by Magnanimous:

* `source/static`    - static files that should be served as-is.
* `source/processed` - files that will be processed with [Magnanimous instructions](https://renatoathaydes.github.io/magnanimous/sections/docs/expression_lang.html).

Any other files and sub-directories are ignored by default, but may be referred to from processed files.

[Magnanimous instructions](https://renatoathaydes.github.io/magnanimous/sections/docs/expression_lang.html) allow the inclusion of files within
templates, the use of variables and [slots](https://renatoathaydes.github.io/magnanimous/sections/docs/expression_lang.html#slot)
to customize inclusions and [components](https://renatoathaydes.github.io/magnanimous/sections/docs/components.html), the addition of
conditional and repeated content, and much more!

Check the [Documentation](https://renatoathaydes.github.io/magnanimous/sections/docs) for more information.
