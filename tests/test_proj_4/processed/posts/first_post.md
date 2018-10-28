{{ define date "2018-01-01" }}
{{ define title "My first blog post" }}
{{ include ../_header.html }}

## Post 1

Hello.

{{ define note "This is a note." }}{{ include /processed/components/_note.html }}

Bye.
{{ include ../_footer.html }}