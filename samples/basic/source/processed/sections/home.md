{{ define title "Home" }}
{{ define index 1 }}
{{ define path "/sections/home.html" }}
{{ define name "Home" }}
{{ include /processed/_header.html }}

# Home

These are my latest posts:

{{ 
    for post /processed/posts/
}}
### <a href="{{ eval post.path }}">{{ eval post.title }}</a>
{{ end }}

{{ include /footer.html }}