{{ define path "/sections/docs/paths.html" }}
{{ define chapter 3 }}
{{ define title "Paths and Links" }}
{{ include _docs_header.html }}

# 3. Paths and Links

When writing your website, you'll probably want to make use of several [Magnanimous instructions]({{ eval INSTRUCTIONS_PATH }})
which need to point to another source file. They do that by telling Magnanimous the **path** to the other file.

You'll also likely add **links** to other files in your pages, so that users can easily navigate between the pages 
of your website.

The two concepts are somewhat similar, but definitely not the same.

A **path** is a compilation-time concept and doesn't exist at runtime, while **links** only exist at runtime.

In the next sections, you'll learn what that means in detail, so you never have to spend time wondering why a link
does not point where you thought it ought to, or a path does not resolve correctly, giving you compile time warnings.

## Paths

[Components](components.html), the [for](expression_lang.html#for) and the [include](expression_lang.html#include)
instructions all need to declare a reference to some file, or directory in the case of [for](expression_lang.html#for).

They need to know where to find the information necessary to do their job: the contents of a file, in the case of
the `include` and `component` instructions, and the directory containing the files to be iterated over, in the case of 
the `for` instruction.

Paths can be of three types:

* Absolute paths
* Relative paths
* Up paths

### Absolute paths

**Absolute paths** start with a `/`, and are always interpreted as being under the
[`source/`](get_started.html#magnanimous-directories) directory.

Example:

```
\{{ include /processed/_header.html }}
```

### Relative paths

**Relative paths**, unlike absolute paths, do not start with a `/` and are resolved by the location of the file they
are declared in.

Example:

```
\{{ include _header.html }}
```

If the above inclusion is done from the file `/processed/components/_comp1.html`, then the path, `_header.html`, will point to the file at `/processed/components/_header.html`.

If the inclusion is from `/processed/_comp2.html`, then it points to `/processed/_header.html`.

Relative paths can also refer to a parent directory if they start with `../`, just like in most file systems.

However, trying to access anything above the `source/` directory is forbidden and will result in the `../` part
being ignored.

### Up paths

**Up paths** refer to a file located at the same location of the **file currently being written** by Magnanimous (not necessarily the location of the file the `include` statement is declared in, as with absolute and relative paths),
or any parent directory **up to** the `source` directory. 

Example:

```
\{{ include .../_messages }}
```

_Up paths_ are special in that the file they point to depends on the context. They are commonly used to allow overriding which file should be used, which is useful for things like internationalization.

For example, given this file structure:

```
source
└── processed
    ├── _header.html
    ├── index.html
    ├── _messages
    └── pt
    |   ├── _messages
    |   └── index.html
    └── sv
        └── index.html
```

Assuming `source/processed/_header.html` has the following contents:

```html
\{{ include .../_messages }}
<html>
<head>
<title>\{{ eval title }}</title>
</head>
<body>
```

> Notice that the header file uses an _up path_, `.../_messages`.

If the `source/processed/index.html` file `include`s this header, it will get messages from `source/processed/_messages` because that's the nearest `_messages` file looking up the directory tree.

If the `source/processed/pt/index.html` file `include`s the header, its messages will come from `source/processed/pt/_messages` instead!

If the `source/processed/sv/index.html` file `include`s the header, because there's no `_messages` file on the same directory, Magnanimous will look up the directory tree, and find the `source/processed/_messages` file, so it will
include that.

Simple, but effective.

### Expression paths

Paths can also be given by [expressions](expression_lang.html#expressions) if starting with `"` or <code>\`</code>,
or explicitly with `eval <expr>`.

For example, the paths in these [include](expression_lang.html#include) instructions all evaluate to 
`/hello/world/file.html`:

```
\{{ include /hello/world/file.html }}
\{{ include "/hello" + "/world" + "/file.html" }}

\{{ define dir "/hello/world/" }}
\{{ include eval dir + "file.html }}
```

### Path variables

Even though most of the time, paths are declared as simple expressions or unquoted strings, it is possible to
explicitly declare a path as follows:

```javascript
\{{ define my_path path["/processed/example/file.txt"] }}
```

This is useful to read the top-level variables defined by another file.
For example, suppose the file at `/processed/example/file.txt` contained the following definition:

```javascript
\{{ define astronaut "Neil Armstrong" }}
```

Then, given `my_path` defined above, we can evaluate this definition in another file:

```javascript
The first astronaut to step on the Moon was \{{ eval my_path.astronaut }}.
```

It is even possible to refer to the current file:

```javascript
This file is at \{{ eval path["."] }}.
```

Path variables can also be used with [date expressions](expression_lang.html#expressions):

```javascript
\{{ define this path["."] }}
This file `\{{eval this}}` was last updated on \{{ eval date[this] }}.
```

Here's how the above code sample renders on this document itself!

{{ define this path["."] }}\
> This file ``{{eval this}}`` was last updated on {{ eval date[this] }}.

### When are paths resolved?

Paths are resolved when Magnanimous is writing the generated website files.

Notice that paths are a compile-time concept, which means that they are only used by Magnanimous instructions, during
compilation, to resolve the contents of the files that will be part of the website. For this reason, paths can
point to _no-writable_ files (whose file names start with `_`, which are not copied to the final website).

Once the website is created, instructions and paths simply do not exist anymore! All you have is a bunch of static
files ready to be served by a _dumb_ web server.

For this reason, paths should always point to other source files, not to generated files.

Links, on the other hand, are a different story, as we'll see.

## Links

Links are not really a Magnanimous concept, but a _Markdown_ and _HTML_ concept. But because, like paths, they point
to another file (or website), they can be confused with paths, so some explanation of what's different between them
is warranted.

Example:

#### Markdown

```markdown
My favourite site is [Wikipedia](https://www.wikipedia.org/).
```

#### HTML

```html
<p>My favourite site is <a href="https://www.wikipedia.org/">Wikipedia</a>.</p>
```

The big difference between links and paths is simple: links are only ever used at runtime,
while paths are used by Magnanimous at compile time to resolve [Magnanimous instructions]({{ eval INSTRUCTIONS_PATH }}).

This has an important implication: while paths point to source files, links point to generated files (or other websites).

So, if you want to add a link to some content you wrote in a markdown file, say `/processed/blog/post.md`, the link to it
should point to `/processed/blog/post.html`, the generated file, not the source file.

> Notice that all processed `.md` files are converted to `.html` files by Magnanimous.
  See the [Markdown Guide](markdown_guide.html) for details.

### When are links resolved?

Links are resolved by the browser when users click on links on your HTML pages.

They are part of the language you're using to write content, be it HTML or Markdown, and are not used by Magnanimous.

### Consider the base path in your links

If your website is not served from the root path within your domain, you might also need to consider the base path!

For example, if your website will be served under the `https://user1.github.io/my-project` URL (as a project called
`my-project` by `user1`, served via [GitHub pages](https://pages.github.com/) would), then your links need to take the
`/my-project` path (the base path) into account, so an absolute link to the source file `/processed/blog/post.md`
should actually point to `/my-project/blog/post.html`.

A good way to approach this is to define a `basePath` variable in the `/processed/_global_context` file, as in this
example:

```
\{{ define basePath "/my-project/" }}
```

{{ component /processed/components/_file-box.html }}\
{{ define file "source/processed/_global_context" }}\
{{ end }}

Now, from any other processed file, a link can be safely created as follows:

```html
<a href="\{{ eval basePath + `blog/post.html` }}">A link</a>
```

Or, in markdown:

```markdown
[A link](\{{ eval basePath + "blog/post.html" }})
```

{{ include _docs_footer.html }}
