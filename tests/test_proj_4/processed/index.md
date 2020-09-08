{{ define title "My blog"}}{{ includeHTML _header.html }}
# Welcome
<div class ="interessant">This is a website.</div>
My blog posts:
{{ for post (sortBy date reverse limit 4) posts }}
* {{ eval post.date }} - <a href="posts/{{ eval post.file }}">{{ eval post.title }}</a>
{{ end }}
{{ includeHTML _footer.html }}