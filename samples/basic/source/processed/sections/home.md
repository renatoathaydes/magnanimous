{{ define title "Home" }}
{{ define index 1 }}
{{ include /processed/_header.html }}

# Home

These are my latest posts:

{{ 
    for post /processed/posts/ sortBy date limit 10 reverse true
}}
### {{ eval post.title }}
{{ end }}

{{ include /footer.html }}