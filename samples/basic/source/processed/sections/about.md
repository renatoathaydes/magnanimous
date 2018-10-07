{{ define name "Home" }}
# Home

These are my latest posts:

{{ 
    for post posts/ sortBy date limit 10 reverse
}}
### {{ $post.title }}
{{ end }}

{{ include /footer }}