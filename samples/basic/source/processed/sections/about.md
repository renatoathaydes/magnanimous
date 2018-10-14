{{ define name "Home" }}
{{ define index 2 }}
# Home

These are my latest posts:

{{ 
    for post posts/ sortBy date limit 10 reverse true
}}
### {{ eval post.title }}
{{ end for }}

{{ include /footer.html }}