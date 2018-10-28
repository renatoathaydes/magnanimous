{{ define title "My blog"}}{{ include _header.html }}
# Welcome
<div class ="interessant">This is a website.</div>
My blog posts:
{{ for post (sortBy date reverse limit 4) posts }}
* {{ eval post.date }} - {{ eval post.title }}
{{ end }}
{{ include _footer.html }}