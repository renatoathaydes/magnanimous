{{ include _header.html }}
# Welcome
<div class ="interessant">This is a website.</div>
My blog posts:
{{ for post posts }}
* {{ eval post.date }} - {{ eval post.title }}
{{ end }}
{{ include _footer.md }}