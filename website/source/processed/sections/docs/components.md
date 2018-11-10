{{ define path "/sections/docs/components.html" }}\
{{ define chapter 4 }}\
{{ define title "Using Components" }}\
{{ include _docs_header.html }}\

# Using Components

As you develop your website, you will notice that you will keep using similar patterns again and again
to create things like text boxes, data tables, dividers, and more obvious things like page headers and
footers.

To avoid repeating the same code all over the place, you can use Magnanimous components.

## Creating and using a component

A component should be declared in a non-writable file (i.e. its name starts with underscore, as in `_example`)
somewhere under the `source/processed/` directory. That's because components are not usually full HTML pages,
but fragments which can be customized to display different content, sometimes even in a customizable shape
(so you wouldn't want them served by the server by themselves).

As an example, let's say we want to create a component to represent a simple _warning_ box.

It could look like this:

```html
\{{ doc This component takes a 'message' and displays it in a warning box }}\\
<div class="warning">\{{ eval message }}</div>
```

{{ component /processed/components/file-box.html }}\
    {{ define file "source/processed/components/_warning.html" }}
{{ end }}

That's it!

You could then use this component in other files like this:

```html
<div>Some content</div>

\{{ component /processed/components/_warning.html }}
    Warning Box!! This text is ignored, only definitions in this scope are evaluated.
    \{{ define message "This is a warning message!" }}
\{{ end }}

<div>More content</div>
```

> Notice that all customizable components must go in the `source/processed/` directory (so they can receive
  the customization via variables), but where under that directory does not matter for Magnanimous... they could
  go under any sub-directory you wish. If you create static components that don't need to be customized,
  you could even put them under `source/static`, which guarantees they will be included without changes.

### Component scopes

Notice that our example component expects a variable called `message` to exist in its scope. For that reason, 
in the above example, we provide a `message` binding within the scope of the component itself, so that we don't
mess with the surrounding scope.

If that's not a concern, we could have just re-used an existing binding, as in this example:

```html
\{{ define message "This is a warning message!" }}\\
<div>Some content</div>

\{{ component /processed/components/_warning.html }}\{{ end }}

<div>More content</div>
```

Notice that the usual scope rules which apply to file inclusions also apply to components.

So, the `message` variable for the example component could even have come from another file which includes this file,
or ultimately, from the `_global_context`.

### A more advanced example

To really understand what can be achieved with components, let's look at a more advanced example.

Let's create a data table which will contain a summary of the metadata for each file in a given directory.

This could be used as a table of contents, or blog post summary, for example.

First, we expect the `source/processed/posts/` directory to only contain files that declare the following variables:

* `name`
* `path`
* `date`

To make things more concrete, here's what the files could look like:

```markdown
\{{ define name "My first post" }}\\
\{{ define path "/processed/posts/first_post.html" }}\\
\{{ define date "2018-04-05" }}\\

# My first post

etc...
```

{{ component /processed/components/file-box.html }}\
    {{ define file "source/processed/posts/first_post.md" }}
{{ end }}

```markdown
\{{ define name "My second post" }}\\
\{{ define path baseURL + "/processed/posts/second_post.html" }}\\
\{{ define date "2018-06-07" }}\\

# My second post

etc...
```

{{ component /processed/components/file-box.html }}\
    {{ define file "source/processed/posts/second_post.md" }}
{{ end }}

Now, we can define a simple HTML component that will put the posts' metadata in a HTML table:

```html
\{{ doc
    Arguments:
      * dataDirectory - path to a directory containing files with the following properties:
          * name
          * path
          * date
}}\\
<table class="data-table-component">
    <thead>
    <th>Post Date</th>
    <th>Post Name</th>
    </thead>
    <tbody>
    \{{ for post (sortBy date) eval dataDirectory }}\\
    <tr>
        <td>\{{ eval post.date }}</td>
        <td><a href="\{{ eval post.path }}">\{{ eval post.name }}</a></td>
    </tr>
    \{{ end }}\\
    </tbody>
</table>
```

{{ component /processed/components/file-box.html }}\
    {{ define file "source/processed/components/_data_table.html" }}
{{ end }}

Finally, we can add the component to our index page (and any other pages we want!):

```html
<html>
<body>
<h2>These are my posts</h2>
\{{ component /processed/components/_data_table.html }}\\
    \{{ define dataDirectory "/processed/posts" }}\\
\{{ end }}\\
</body>
</html>
```

{{ component /processed/components/file-box.html }}\
    {{ define file "source/processed/index.html" }}
{{ end }}

And we're done!

Here's what the result should look like:

```html
<html>
<body>
<h2>These are my posts</h2>
<table class="data-table-component">
    <thead>
    <th>Post Date</th>
    <th>Post Name</th>
    </thead>
    <tbody>
    <tr>
        <td>2018-04-05</td>
        <td><a href="/processed/posts/first_post.html">My first post</a></td>
    </tr>
    <tr>
        <td>2018-06-07</td>
        <td><a href="/processed/posts/second_post.html">My second post</a></td>
    </tr>
    </tbody>
</table>
</body>
</html>
```

{{ component /processed/components/file-box.html }}\
    {{ define file "target/index.html" }}
{{ end }}

Now, go on and create your own awesome components!

{{ include _docs_footer.html }}
