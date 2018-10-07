{{ include /processed/_header }}

# Home

These are my latest posts:

{{ 
    for post /processed/posts/ sortBy date limit 10 reverse
}}
### {{ $post.title }}
{{ end }}

{{ include /footer }}