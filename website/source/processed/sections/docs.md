{{ define path "/sections/docs.html" }}
{{ define name "Documentation" }}
{{ define title "Magnanimous Docs" }}
{{ define index 5 }}
{{ include /processed/_header.html }}

# Magnanimous Documentation

Here, you'll find all information you need to use Magnanimous effectively.

If you're new to Magnanimous, the best place to go first is the
[Getting Started]({{ eval baseURL + "/sections/docs/get_started.html" }}) page.

To learn more once you've been there, follow one of the Tutorials listed below.

Once you've mastered Magnanimous (which should be really quick and easy), you can just go directly to the
[Reference]({{ eval baseURL + "/sections/docs/reference.html" }}) when you just need to refresh your memory on
how to achieve what you want, fast.

### Table Of Contents

{{ for doc (sortBy index) docs }}
1. [{{ eval doc.title }}]({{ eval baseURL + doc.path }}) page.
{{ end }}