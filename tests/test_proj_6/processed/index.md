## Simple Component Example

{{ component /processed/components/_div_wrapper.html }}\
{{ define content "Hello components" }}
{{ end }}

## HTML Component with Markdown contents

{{ component /processed/components/_div_wrapper.html }}\

Bar

### inner markdown

Foo
{{ end }}

Example HTML:

```html
{{ include _html_in.md }}
```
