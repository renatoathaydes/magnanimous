{{ define path "/sections/docs/components.html" }}\
{{ define chapter 4 }}\
{{ define title "Using Components" }}\
{{ include _docs_header.html }}\

# 4. Using Components

As you develop your website, you will notice that you will keep using similar patterns again and again
to create things like text boxes, data tables, dividers, and more obvious things like page headers and
footers.

To avoid repeating the same code all over the place, you can use Magnanimous components.

{{ component /processed/components/_linked_header.html }}\
{{ define id "creating-and-using-component" }}\
Creating and using a component
{{ end }}

A component should normally be declared in a non-writable file (i.e. its name starts with underscore, as in `_example`)
somewhere under the `source/processed/` directory.

That's because components are not usually full HTML pages, but fragments which can be customized to display 
customizable content (so you wouldn't want them served by the server by themselves, without "filling" it in).

As an example, let's say we want to create a component to represent a simple _warning_ box.

It could look like this:

```html
\{{ doc This component takes its contents and displays it in a warning box }}\\
<div class="warning">\{{ eval __contents__ }}</div>
```

{{ component /processed/components/_file-box.html }}\
    {{ define file "source/processed/components/_warning.html" }}
{{ end }}

> Notice that `__contents__` is a special variable that holds the contents of the component when it's used
  (i.e. the content wrapped between \{{ component path }} and \{{ end }} as we'll see below).

That's it!

You could then use this component in other processed files like this:

```html
<div>Some content</div>

\{{ component /processed/components/_warning.html }}
    This is a warning message!
\{{ end }}

<div>More content</div>
```

Which should render like this:

```html
<div>Some content</div>

<div class="warning">
    This is a warning message!
</div>

<div>More content</div>
```

{{ component /processed/components/_linked_header.html }}\
{{ define id "customizing-components-with-variables" }}\
Customizing components with variables
{{ end }}

Components allow the user to declare variables that customize its contents. The variables may be declared before the
component's declaration, but they can also be placed inside the component's body, making the variables "local" to the
component (i.e. not visible in the surrounding scope).

Let's look at an example to clarify what that means. This component expects 2 variables (`my_variable` and `other_var`)
to be set for customizing it:

```html
<h2>\{{ eval my_variable }}</h2>
<div>
    <span>\{{ eval other_var }}</span>
</div>
```

{{ component /processed/components/_file-box.html }}\
    {{ define file "source/processed/components/_var_example.html" }}
{{ end }}

When using it, we need to provide values for both variables, which we can do both from outside and from inside the
component:

```html
\{{ define my_variable "This is available in the component's scope" }}\\

\{{ component /processed/components/_var_example.html }}
    \{{ define other_var "This also!" }}\\
\{{ end }}
```

{{ component /processed/components/_file-box.html }}\
    {{ define file "source/processed/example.html" }}
{{ end }}

Result:

```html

<h2>This is available in the component's scope</h2>
<div>
    <span>This also!</span>
</div>
```

{{ component /processed/components/_file-box.html }}\
    {{ define file "target/example.html" }}
{{ end }}


{{ component /processed/components/_linked_header.html }}\
{{ define id "customizing-components-with-slots" }}\
Customizing components with slots
{{ end }}

`slot`s make components extremely powerful! They allow the creation of very modular components because the parts
of the component can be provided by slots defined elsewhere.

> Magnanimous slots were inspired by the HTML5 
  [Web Components](https://developer.mozilla.org/en-US/docs/Web/Web_Components/Using_templates_and_slots) standard.

For example, we can create a component that places its contents in 3 areas: top, middle and bottom
(where both top and bottom are optional):

```html
<div class="top">\{{ eval top || "" }}</div>
<div class="middle">\{{ eval middle }}</div>
<div class="bottom">\{{ eval bottom || "" }}</div>
```

{{ component /processed/components/_file-box.html }}\
    {{ define file "source/processed/components/_slots_example.html" }}
{{ end }}

To use this component is pretty easy:

```html
<h1>Hello world</h1>
\{{ component /processed/components/_slots_example.html }}
    \{{ slot top }}
        <h1>Top</h1>
        <p>This goes at the top</p>
    \{{ end }}
    \{{ slot middle }}
        <h3>Middle</h3>
        <p>This goes in the middle</p>
    \{{ end }}
    \{{ slot bottom }}
        <h3>Bottom</h3>
        <p>This goes at the bottom</p>
    \{{ end }}
\{{ end }}
```

{{ component /processed/components/_file-box.html }}\
    {{ define file "source/processed/example.html" }}
{{ end }}

Result:

```html
<div class="top">
    <h1>Top</h1>
    <p>This goes at the top</p>
</div>
<div class="middle">
    <h3>Middle</h3>
    <p>This goes in the middle</p>
</div>
<div class="bottom">
    <h3>Bottom</h3>
    <p>This goes at the bottom</p>
</div>
```

{{ component /processed/components/_file-box.html }}\
    {{ define file "target/example.html" }}
{{ end }}


{{ component /processed/components/_linked_header.html }}\
{{ define id "advanced-example" }}\
A more advanced example
{{ end }}

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

{{ component /processed/components/_file-box.html }}\
    {{ define file "source/processed/posts/first_post.md" }}
{{ end }}

```markdown
\{{ define name "My second post" }}\\
\{{ define path baseURL + "/processed/posts/second_post.html" }}\\
\{{ define date "2018-06-07" }}\\

# My second post

etc...
```

{{ component /processed/components/_file-box.html }}\
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

{{ component /processed/components/_file-box.html }}\
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

{{ component /processed/components/_file-box.html }}\
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

{{ component /processed/components/_file-box.html }}\
    {{ define file "target/index.html" }}
{{ end }}

Now, go on and create your own awesome components!

{{ include _docs_footer.html }}
