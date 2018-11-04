{{ define path "/sections/docs.html" }}
{{ define name "Documentation" }}
{{ define title "Magnanimous Docs" }}
{{ define index 5 }}
{{ include /processed/_header.html }}

# Magnanimous Documentation

Here, you'll find all information you need to use Magnanimous effectively.

{{ define getstartedlink baseURL + "/sections/docs/get_started.html" }}

If you're new to Magnanimous, the best place to go first is the
[Getting Started]({{ eval getstartedlink }}) page.

To learn more once you've been there, follow one of the Tutorials listed below.

{{ define referencelink baseURL + "/sections/docs/reference.html" }}

Once you've mastered Magnanimous (which should be really quick and easy), you can go directly to the
[Reference]({{ eval referencelink }}) section when you just need to refresh your memory on
how to achieve what you want, fast.

### Table Of Contents

{{ for doc (sortBy chapter) docs }}
    {{ if doc.path != nil }}
        {{ define docLink baseURL + doc.path }}
#### {{ eval doc.chapter }}. [{{ eval doc.title }}]({{ eval docLink }})
    {{ end }}
{{ end }}

{{ include /processed/_footer.html }}