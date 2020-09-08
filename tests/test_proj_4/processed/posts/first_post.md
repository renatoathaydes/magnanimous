{{ define date "2018-01-01" }}
{{ define file "first_post.html" }}
{{ define title "My first blog post" }}
{{ includeHTML ../_header.html }}

## Post 1

Hello.

{{ define note "This is a note." }}{{ include /processed/components/_note.html }}

Bye.
{{ includeHTML ../_footer.html }}