Hello

{{ includeRaw /processed/other.md }}

{{ define _forceMarkdown 1 }}\
```html
{{ include /processed/_some.html }}
```

END
