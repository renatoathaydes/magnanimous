{{ define path "/sections/docs/paths.html" }}
{{ define chapter 3 }}
{{ define title "Paths and Links" }}
{{ include _docs_header.html }}

# 3. Paths and Links

When writing your website, you'll probably want to make use of several [Magnanimous instructions]({{ eval INSTRUCTIONS_PATH }})
which need to point to another source file. They do that by telling Magnanimous the **path** to the other file.

You'll also likely add **links** to other files in your pages, so that users can easily navigate between the pages 
of your website.

The two concepts are similar, but not quite the same.

A **path** is a compilation-time concept and doesn't exist at runtime, while **links** only exist at runtime.

In the next sections, you'll learn what that means in detail, so you never have to spend time wondering why a link
does not point where you thought it ought to, or a path does not resolve correctly, giving you compile time warnings.

## Paths

[Components](components.html), the [for](expression_lang.html#for) and the [include](expression_lang.html#include)
instructions all need to declare a reference to some file, or directory in the case of [for](#for).

They need to know where to find the information necessary to do their job: the contents of a file, in the case of
the `include` and `component` instructions, and the directory containing the files to be iterated over, in the case of 
the `for` instruction.

Paths can be of two types:

* Absolute paths
* Relative paths

**Absolute paths** start with a `/`, and are always interpreted as being under the
[`source/`](get_started.html#magnanimous-directories) directory.

**Relative paths**, unlike absolute paths, do not start with a `/` and are resolved by the location of the file they
are declared in.

As an example, if within the file `/processed/components/_comp1.html`, you have a instruction that uses the path
`example/_comp2.html`, then `example/_comp2.html` points to the absolute path `/processed/components/example/_comp2.html`.

Paths can refer to a parent directory with `../`, just like in most file systems. However, trying to access anything
above the `source/` directory is forbidden and will result in the `../` part being ignored.

Notice that paths are a compile-time concept, which means that they are only used by Magnanimous instructions, during
compilation, to resolve the contents of the final files that will be part of the website. For this reason, paths can
point to _hidden_ files (whose file names start with `_`, which are not copied to the final website).

Once the website is created, instructions and paths simply do not exist anymore! All you have if a bunch of static
files ready to be served by a _dumb_ web server.

For this reason, paths should always point to other source files, not to generated files. 

Links, on the other hand, are a different story, as we'll see.

## Links

Links are not really a Magnanimous concept, but a _Markdown_ and _HTML_ concept. But because, like paths, they point
to another file, they can be confused with paths, so some explanation of what's different between them is warranted.

The big difference is simple: links are used at runtime by users of your website, while paths are used by Magnanimous
at compile time, to resolve [Magnanimous instructions]({{ eval INSTRUCTIONS_PATH }}).

This has an important implication: while paths point to source files, links point to generated files (or other websites).

So, if you want to add a link to some content you wrote in a markdown file, say `/processed/blog/post.md`, the link to it
should point to `/blog/post.html`, the generated file, not the source file.

> Notice that all processed `.md` files are converted to `.html` files by Magnanimous.
  See the [Markdown Guide](markdown_guide.html) for details.

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

Now, from anywhere, a link can be safely created as follows:

```html
<a href="\{{ eval basePath + `blog/post.html` }}">A link</a>
```

Or, in markdown:

```markdown
[A link](\{{ eval basePath + "blog/post.html" }})
```

{{ include _docs_footer.html }}
