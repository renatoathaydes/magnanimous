{{ define path "/sections/docs.html" }}
{{ define name "Documentation" }}
{{ define title "Magnanimous Docs" }}
{{ define index 5 }}
{{ include /processed/_header.html }}

# Magnanimous Documentation

Here, you'll find all information you need to use Magnanimous effectively.

{{ define getstartedlink baseURL + "/sections/docs/get_started.html" }}\

If you're new to Magnanimous, the best place to go first is the
[Getting Started]({{ eval getstartedlink }}) page.

{{ define tutoriallink baseURL + "/sections/docs/basic_tutorial.html" }}\

To learn more once you've been there, follow the 
[Basic Tutorial]({{ eval tutoriallink }}) and you should be ready to go!

{{ define langlink baseURL + "/sections/docs/expression_lang.html" }}\

When you need help remembering Magnanimous syntax and instructions, head to the 
[Expression Language]({{ eval langlink }}) page.

### Table Of Contents

{{ for doc (sortBy chapter) docs }}
    {{ if doc.path != nil }}
        {{ define docLink baseURL + doc.path }}
#### {{ eval doc.chapter }}. [{{ eval doc.title }}]({{ eval docLink }})
    {{ end }}
{{ end }}

{{ include /processed/_footer.html }}