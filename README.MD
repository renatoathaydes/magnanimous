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

* `{{ include <path> }}` - include the contents of the file in the given path into the current file.
* `{{ define <varible_name> <expression> }}` - define a variable.
* `{{ eval <expression> }}` - evaluates an expression.
* `{{ if <expression> }}` (conditional content) `{{ end }}`.
* `{{ for <item> <iterable> }}` (repeated content for each item) `{{ end }}`.

Markdown files (`*.md`) will be converted to HTML, after Magnanimous instructions are processed, using
[Blackfriday](https://github.com/russross/blackfriday).

Files whose name starts with `_` will not be present in the output website, but may be included by other files.

#### Inclusions

`include` instructions must always have a path to a file as an argument.

An inclusion may look like this: 

```
{{ include /static/header.html }}
```

This whole expression will be replaced with the contents of the file at `source/static/header.html`.
If the file does not exist, a warning will be emitted
during compilation and the expression will appear as it is in the compiled resource.

#### Variables

Variables can be defined with the `define` instruction.
 
For example:

```
{{ define title "My Website" }}
```

The above instructions defines a variable called `title` with the value `My Website`.

To evaluate the value of a variable, the `eval` instruction is used:

```
{{ eval title }}
```

The variable may be either defined in the same file as the evaluation is happening,
or in a file which includes that file.

If the variable cannot be found, a warning will be emitted during compilation, and the expression will appear as it is
in the compiled resource.

> The value of a variable is evaluated as an *expression*. See the next section for more details.

It's possible to write more advanced expressions with variables as in the following example:

```
< div class="header {{ eval pageName == "Home" ? "highlight" : "" }}">Home</div>
```

#### Expressions

Expressions may appear in several different locations:

* as an argument to the `eval` instruction.
* as the value of a new variable defined by the `define` instruction.
* as an argument to the `if` instruction.
* as items of an iterable (array), used with the `for` instruction.

For example, an expression may be evaluated with `eval`, with the result being displayed in the page:

```
{{ eval 2 + 2 * 5 }}
```

This should result in `12` being displayed.

> Notice that expressions use a syntax very similar to C-like languages.

The value of a new variable may also be determined by an expression:

```
{{ define n 10 + 10 * (2 + 4) }}
```

This results in the variable `n` assuming the value `70`.

When an expression is used as an argument to the `if` instruction, it must evaluate to a boolean value (`true` or `false`).

For example, the following is valid (given variable `n` is defined):

```
{{ if n > 50 }}N is large{{ end }}
```

Finally, it's possible to use expressions within arrays declared with the `for` instruction:

```
{{ for x [1 + 1, 2 + 2] }}
X is {{ eval x }}
{{ end }}
```

#### Conditional content

The `if` instruction allows processed files to include content conditionally.

Its argument is a single expression that is expected to evaluate to a boolean value.

> Expressions that evaluate to any non-boolean value emit a warning and evaluate to `false`.

For example:

```
{{ if title == "Home" }}
<h1>Welcome Home!</h1>
{{ end }}
```

In this example, the content `<h1>Welcome Home!</h1>` will only be displayed if the `title` variable has the
value `"Home"`, otherwise, nothing would be displayed.

#### Repeated content (for loops)

The `for` instruction has the following arguments:

* `item`       - a variable name to associate with each item of `iterable`.
* `iterable`   - an iterable object. 

Two kinds of objects are iterable:

* array expressions (e.g. `[1 2 3]`, `["Home", "About"]`).
* file references to directories.

Example iterating over an array:

```
{{ for title ["Home", "About", "FAQ"] }}
<div>{{ eval title }}</div>
{{ end for }}
```

The above, after processing, should look like this:

```html
<div>Home</div>
<div>About</div>
<div>FAQ</div>
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
{{ define columns ["ABC", "GHI", "DEF"] }}
{{ for col $columns sort reverse }}
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
{{ define date "2014-01-01 23:59:59" }}
{{ define title "Example post" }}
```

We could display the newest 5 posts' titles using the following snippet:

```
{{ for post "/processed/posts" sortBy date limit 5 reverse }}
  <div>{{ eval post.title }}</div>
{{ end }}
```