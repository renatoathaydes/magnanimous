# Magnanimous

The simplest and fastest static website generator in the world!!

Magnanimous has the following guiding principles:

* no configuration.
* no magic.
* explicit is better than implicit.
* generate static websites well, nothing else.

This project is currently in early development and is not ready to be used. But stay tuned!

## Basic Design

Magnanimous reads source files from the `source/` directory, and generates a website in `target/`.

The following `source` sub-directories are the only ones treated specially by Magnanimous:

* `source/static`    - static files that should be served as-is.
* `source/processed` - files that will be processed, as explained below.

Any other files and sub-directories are ignored by default, but may be referred to from processed files.

### Processed files

Files meant to be processed by Magnanimous are placed inside the `source/processed` directory.

Processed files may use the following processor instructions:

* `{{ include <file> }}` - include the contents of another file into the current file.
* `{{ define <varible_name> <variable_value> }}` - define a variable.
* `{{ if <expression> }}` (conditional content) `{{ end }}`.
* `{{ for <item> <iterable> }}` (repeated content for each item) `{{ end }}`.

Markdown files (`*.md`) will be converted to HTML after Magnanimous instructions are processed, using
[Blackfriday](https://github.com/russross/blackfriday).

Files whose name starts with `_` will not be present in the output website, but may be included by other files.

#### Inclusions

`include` instructions must always have a file reference as an argument.

An inclusion may look like this: 

```
{{ include /static/header.html }}
```

This whole expression will be replaced with the contents of the file at `source/static/header.html`.
If the file does not exist, a warning will be emitted
during compilation, and the expression will appear as it is in the compiled resource.

The file extension may be omitted if there is no ambiguity.

#### Variables

Variables can be resolved in any processed file by using the syntax: `{{ $file.variable_name }}`.

If no `file` is given, the variable will be resolved as follows, in this order:

* a variable definition in the file which included this file, if any.
* a variable definition in the current file.

This allows inclusions to "override" local variables.

If the variable cannot be found, a warning will be emitted during compilation, and the expression will appear as it is
in the compiled resource.

The following is an example of a simple variable definition:

```
{{ define title "My Website" }}
```

> The value of a variable is evaluated as a [govaluate](https://github.com/Knetic/govaluate) expression.

Using the variable:
                                                               
```
<title>{{ $title }}</title>
```

Whic should compile to:

```
<title>My Website</title>
```

it's also possible to write more advanced expressions with variables as in the following example:

```
< div class="header {{ $pageName == "Home" ? "highlight" : "" }}">Home</div>
```

#### Conditional content

The `if` instruction allows processed files to include content conditionally.

Its argument is a single [govaluate](https://github.com/Knetic/govaluate) expression.

For example:

```
{{ if $pageName == "Home" }}
<h1>Welcome Home!</h1>
{{ end }}
```

In this example, the content `<h1>Welcome Home!</h1>` will only be displayed if the `pageName` variable has the
value `"Home"`, otherwise, nothing would be displayed.

#### Repeated content (for loops)

The `for` instruction has the following arguments:

* `item`       - a variable name to associate with each item of `iterable`.
* `iterable`   - an iterable object. 

Two kinds of objects are iterable:

* arrays, as evaluated by [govaluate](https://github.com/Knetic/govaluate).
* file references to directories.

Example iterating over an array:

```
{{ for (1, 2, 3, 4, 5) n }}
<div>{{ $n }}</div>
{{ end }}
```

The above, after processing, should look like this:

```html
<div>1</div>
<div>2</div>
<div>3</div>
<div>4</div>
<div>5</div>
```

See the section on iterables for details on how to work with iterables.

File references are discussed in the next section.

### File references

References to other files, both from processed files and from HTML files that link to css or images, for example,
may be relative or absolute.

Absolute references start with a `/` and are always resolved with relation to `source/`.

Examples:

* `/images/ico.png` ---> `source/images/ico.png` (no processing).
* `/processed/template.html` ---> `source/processed/template.html` (file contents will be processed).

Relative references are resolved relative to the file declaring the reference.

For example, within the `source/processed/template.html` file:

* `about/index.html` ---> `source/processed/about/index.html` (file contents will be processed).
* `../images/ico.png` ---> `source/images/ico.png` (no processing).

Notice that referring to parent directories with `../` only works up to the `source/` directory:

* `../../images/ico.png` --> `source/images/ico.png` (no processing).

#### Iterating over file references

File references pointing to a directory can be used with the `for` instruction as iterables over each file
inside the directory.

For example, to create a HTML menu containing links to each file under the `source/processed/posts/` directory:

```
{{ for post /processed/posts }}
    <div>
        <a href="/processed/posts/{{ $post.name }}">
        {{ post.date }} - {{ $post.description }}
        </a>
    </div>
{{ end }}
```

### Working with iterables

Iterables are object that can be iterated over with the `for` instruction.

As we've seen before, only arrays and file references to directories can be iterable.

The following operations apply to any iterable:

* `limit <n>`    - limit iterations to at most `n`, where `n` is a non-negative integer.
* `sort`         - sort items based on their value.
* `reverse`      - reverse the items.

File reference iterables also have the following operations:

* `sortBy <var>` - sort items based on the value of a variable within each file.

Examples:

```
<tr>
{{ 
    define columns ("ABC", "GHI", "DEF")
    for col $columns sort reverse
}}
<th>{{ $col }}</th>
{{ end }}
</tr>
```

Results in:

```html
<tr>
<th>GHI</th>
<th>DEF</th>
<th>ABC</th>
</tr>
```

Given a directory `source/processed/posts/` containing 10 files called `p1.md`, `p2.md`, and so on, each of which
defining the `date` and `title` variables as in this example:

```md
{{ define date "2014-01-01 23:59:59" define title "Example post" }}
```

We could display the newest 5 posts' titles using the following snippet:

```
{{ for post "/processed/posts" sortBy date limit 5 reverse }}
  <div>{{ $post.title }}</div>
{{ end }}
```