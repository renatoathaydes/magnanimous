{{ define path "/sections/docs/expression_lang.html" }}\
{{ define chapter 2 }}\
{{ define title "Expression Language" }}\
{{ include _docs_header.html }}\

# 2. Expression Language
 
Magnanimous defines a very simple _expression language_ that is used to **process** documents in the `source/processed`
directory.

This expression language is composed solely of instructions. Each instruction has the following form:

```
\{{ instruction-name args }}
```

_Where:_

* `instruction-name` - the name of the instruction to be used.
* `args` - optional arguments for the instruction. Its syntax is instruction-dependent.

Only a single instruction may be present within double-braces (i.e. between `\{{` and `}}`).

Here's a list of all Magnanimous instructions:

* [`define`](#define)         - defines a variable. Its value is given by an [expression](#expressions).
* [`eval`](#eval)             - evaluates an [expression](#expressions) and inserts the result into the current position.
* [`include`](#include)       - includes contents of another file into the current position.
* [`includeB64`](#includeB64) - includes base64-encoded contents of file into the current position.
* [`component`](#component)   - includes a [Component](components.html) into the current position.
* [`slot`](#slot)             - defines a variable whose content is the body of the instruction.
* [`if`](#if)                 - conditionally includes some content into the current position.
* [`for`](#for)               - repeats some content for each item in an [iterable](#iterables).
* [`doc`](#doc)               - allows documentation to be added to sources (not included in the resource).
* [`end`](#end)               - ends a scoped instruction (`component`, `slot`, `if` and `for`).

{{ component /processed/components/_linked_header.html }}\
{{ define id "instructions" }}\
{{ define text "Instructions" }}\
{{ end }}

In this section, we'll see details about every available instruction in Magnanimous, and how to use them.

{{ component /processed/components/_linked_header.html }}\
{{ define id "define" }}{{ define tag "h3" }}\
{{ end }}

#### Syntax:

```
\{{ define <variable-name> <expression> }}
```

_where:_

* `variable-name` is an [identifier](#identifiers).
* `expression` is an [expression](#expressions).

The `define` instruction is used to _define_ a variable that can later be used in one or more expressions.

For example, the following instruction defines a variable called `text` with the value `Hello Magnanimous`:

```
\{{ define text "Hello Magnanimous" }}
```

After this declaration, any expression where the `text` variable is used will have its value, `Hello Magnanimous`,
used upon evaluation.

For example, we could add the following definition after the one we've just shown above:

```
\{{ define other_text text + ", this is great!" }}
```

Because `text` was previously defined as `Hello Magnanimous`, this definition results in `other_text` getting the value
`Hello Magnanimous, this is great!`!

See [Expressions](#expressions) for more details about the kind of expressions you can use.

{{ component /processed/components/_linked_header.html }}\
{{ define id "eval" }}{{ define tag "h3" }}\
{{ end }}

#### Syntax:

```
\{{ eval <expression> }}
```

_where:_

* `expression` is an [expression](#expressions).

The `eval` instruction is similar to `define`, except that the result of the expression is inserted into the 
processed document instead of being set to a variable.

As an example, suppose you have the following _md_ file:

```markdown
## The dog is \{{ eval "big" }}
```

After processing, this file would be turned into HTML and its instructions processed, resulting in the following
content:

```html
<h2>The dog is big</h2>
```

Of course, `eval` is mostly useful when combined with variables, [slots](#slot) and expressions.

For example, you could define a number of variables beforehand, to be used later in several other places
(so you could change them in only one place if you ever changed your mind about their values):

```markdown
\{{ define title "OurWebsite" }}\\
\{{ define visitorsPerMonth 10000 }}\\

The best website in the world, \{{ eval title }}, has an average number of
\{{ eval visitorsPerMonth / 30 }} visitors per day!

Many people love \{{ eval title }} because it is so great!  
```

> You may have noticed that instructions often end with a `\` character. That's to avoid a new-line
  character being inserted where the instructions were in the source file, as the `\` can escape new lines.

{{ component /processed/components/_linked_header.html }}\
{{ define id "include" }}{{ define tag "h3" }}\
{{ end }}

#### Syntax:

```
\{{ include <path> }}
```

_where:_

* `path` is a [path](paths.html) to another file.

The `include` statement is used to include the contents of a file into another file.

Example:

```html
<div id="other-file-contents">
\{{ include path/to/other/file.md }}
</div>
```

Notice that paths can be [expressions](#expressions) if starting with `"` or <code>\`</code>,
or explicitly with `eval <expr>`:

```html
<!-- Include expression path -->
\{{ include "other/" + fileName }}

<!-- Include evaluated path -->
\{{ include eval fileName + "2" }}
```

Paths may be absolute, relative or _up-paths_ (i.e. paths that refer to a file at the same directory, or a directory
up the source tree from the current file being written).

Absolute path example: `/processed/partials/_header.html`

Relative path example: `partials/_header.html`

Up-path example: `.../_header.html`

See [Paths and Links](paths.html) for more details.

{{ component /processed/components/_linked_header.html }}\
{{ define id "includeB64" }}{{ define tag "h3" }}\
{{ end }}

#### Syntax:

```
\{{ includeB64 <path> }}
```

_where:_

* `path` is a [path](paths.html) to another file.

The `includeB64` statement is used to include the [base64](https://wikipedia.org/wiki/Base64)-encoded contents of a file into another file.

Example:

```html
<!-- include a gif file encoded as a data URL -->
<img src="data:image/gif;base64,\{{ includeB64 path/to/file.gif }}">
```

See the [include](#include) documentation for more information about paths that may be used.

{{ component /processed/components/_linked_header.html }}\
{{ define id "component" }}{{ define tag "h3" }}\
{{ end }}

#### Syntax:

```
\{{ component <path> }}
<content>
\{{ end }}
```

_where:_

* `path` is a [path](paths.html) to another file.
* `content` the body of the component.

The `component` instruction is quite similar to `include`. It also includes the contents of another file, 
the component file (which is usually designed specifically for this purpose), into the file using it.

However, components are more powerful as they may also declare some content which can be placed anywhere inside the 
component (typically using `slot`s).

Example:

A simple component that displays its contents inside a `div` element, with an optional custom class:

```html
<div class="\{{ eval cssClass || `component-example` }}">
\{{ eval __contents__ }}
</div>
```

Including the above component in another file:

```html
\{{ component path/to/component.html }}
\{{ define cssClass "large-text" }}
Include this in my component.
\{{ end }}
```

See [Components](components.html) for more details about using components.

{{ component /processed/components/_linked_header.html }}\
{{ define id "slot" }}{{ define tag "h3" }}\
{{ end }}

#### Syntax:

```
\{{ slot <slot-name> }}
<content>
\{{ end }}
```

_where:_

* `slot-name` the name of this slot.
* `content` the content the slot evaluates to.

A `slot` instruction, like [`define`](#define), defines a variable and does not include any content at the location where it is
declared. But unlike `define`, `slot` uses the body of its declaration as it value.

That means that evaluating the variable defined by `slot` with [`eval`](#eval) results in the body of the slot 
being inserted into the processed document.

Slots are commonly used together with [Components](components.html) (but may also be used on their own), so they are
explained in more detail in the [Components](components.html) page.

{{ component /processed/components/_linked_header.html }}\
{{ define id "if" }}{{ define tag "h3" }}\
{{ end }}

#### Syntax:

```
\{{ if <expression> }}
<content>
\{{ end }}
```

_where:_

* `expression` is an [expression](#expressions), expected to evaluate to a boolean (`true` or `false`).
* `content` is some content that will be included in the document if `<expression>` is true.

The `if` instruction can be used to include some content in a document only if some condition is true.

For example, you may want to include a certain CSS class on an element only if it's the currently active element:

```html
<div class="\{{ if currentPage == page }}active\{{ end }}"></div>
```

{{ component /processed/components/_linked_header.html }}\
{{ define id "for" }}{{ define tag "h3" }}\
{{ end }}

#### Syntax:

```
\{{ for <variable-name> [ (<for-instruction>...) ] <iterable> }}
<content>
\{{ end }}
```

_where:_

* `variable-name` is an [identifier](#identifiers) to be bound for each item of `<iterable>`.
* `for-instruction` instructions for iteration (see below).
* `iterable` is an [iterable](#iterables).
* `content` is some content to be repeatedly included in the document, once for each item.

{{ component /processed/components/_linked_header.html }}\
{{ define id "for-instructions" }}{{ define tag "h4" }}\
{{ define text "for sub-instructions" }}\
{{ end }}

* `sort`            - sort the elements alphabetically.
* `sortBy <field>`  - sort the file items by the value of some property.
* `reverse`         - reverse the order of the items.
* `limit <max>`     - limit the number of items to include.

The `for` instruction allows some content to be repeated for each item of an [iterable](#iterables).

For example, you could use an _array_ to iterate over some values, including the same content for each item:

```html
\{{ for item ["Home", "About", "Documentation"] }}
<div>\{{ eval item }}</div>
\{{ end }}
```

More commonly, the `for` instruction is used to show properties of files in a certain directory.

Example:

```html
\{{ for item (sortBy date reverse limit 10) /path/to/directory }}
<div>Date: \{{ eval date }}</div>
<div>Post name: \{{ eval name }}</div>
\{{ end }}
```

Notice that the path is not normally given as an [expression](#expressions), but as simple text
(notice that the path is not wrapped into double-quotes).
If you need to pass in an expression, or just a variable instead of a hardcoded path, you must call `eval` first:

```html
\{{ define postDirectories "/path/to/directory" }}

<!-- Somewhere else in the file -->
\{{ for item eval postDirectories }}
<div>Date: \{{ eval date }}</div>
<div>Post name: \{{ eval name }}</div>
\{{ end }}
```

See [Iterables](#iterables) for details about what iterable types can be used with the `for` instruction.

{{ component /processed/components/_linked_header.html }}\
{{ define id "doc" }}{{ define tag "h3" }}\
{{ end }}

#### Syntax:

```
\{{ doc <text> }}
```

_where:_

* `text` can be any text (excluding `}}`, which ends the instruction).

The `doc` instruction is used to document your templates and does not produce any visible output in generated
documents.

You can use it to make clear what some complex parts of your templates work, or document the variables expected to
be set for a [Component](components.html), for example.

{{ component /processed/components/_linked_header.html }}\
{{ define id "end" }}{{ define tag "h3" }}\
{{ end }}

#### Syntax:

```
\{{ end }}
```

`end` is not an independent instruction; it is used to end the scope of a previously started scoped instruction.

It does not accept any argument.

All scoped instructions must always be finalized with an `\{{ end }}`.

Scoped instructions are:

* `component`
* `slot`
* `if`
* `for`

Each `end` instruction always matches the nearest unclosed scoped instruction.

{{include /processed/components/_spacer.html }}\
<hr>

{{ component /processed/components/_linked_header.html }}\
{{ define id "expressions" }}\
{{ define text "Expressions" }}\
{{ end }}

Magnanimous expressions use a syntax that's similar to C-like languages, including Java, JavaScript and Go.

For example, the following are all valid expressions:

```javascript
2
2 + 2.42
(5 * 2) - 10
variable == true
!negated == false
something + other != "hello" + yet_another
thing == null
`multiline
string`
[1, 2, 3, 4]
path["path/to/some/file.html"]
date["2019-03-20"]
date["now"]
date[some_path]
date["2019-03-20T08:24:00"]["Mon Jan 2 15:04:05 2006"]
```

The above expressions include all types of variables available:

* Strings: double-quoted as in `"hello""`.
* Multiline Srings: delimited with back-ticks: `` `example` ``.
* Numbers: like `2` or `2.42`.
* Booleans: `true` or `false`.
* Null: the `null` value (i.e. a variable that has not been defined).
* Array: arrays of values of any type (can be used with the [for](#for) instruction.
* Dates: see the dates section below for details.

They also show the use of _variables_, such as `variable` and `negated` above, which must be declared via the
[define](#define) instruction.

The names of the variables must be valid _identifiers_, as we'll see below.

{{ component /processed/components/_linked_header.html }}\
{{ define id "identifiers" }}{{ define tag "h3" }}\
{{ define text "Identifiers" }}\
{{ end }}

_Identifiers_ are used to name variables declared with the [define](#define) instruction.

They must start with a letter or `_`, but may contain any of the following characters:

* `a-z`
* `A-Z`
* `0-9`
* `_`

Examples of valid identifiers:

```
hello
_first_character
abs123
CONSTANT_1
```

{{ component /processed/components/_linked_header.html }}\
{{ define id "operators" }}{{ define tag "h3" }}\
{{ define text "Operators" }}\
{{ end }}

As shown above, expressions may comprise identifiers, numbers, Strings and also operators.

The following operators are supported:

#### Arithmetic operators

* `+` - addition
* `-` - subtraction 
* `*` - multiplication
* `/` - division
* `%` - remainder

#### Comparison operators

* `>`  - greater than
* `>=` - greater than or equal
* `<`  - less than
* `<=` - less than or equal
* `==` - equal
* `!=` - not-equal

#### Boolean operators

* `&&` - AND
* `||` - OR
* `!`  - NOT

The `||` (OR) operator can be used to declare default values for variables that may be missing a value.

For example, in the following [eval](#eval) declaration, the document should get the value of the `example_var`
if it's been defined, _or_ the default value, `Not here`, otherwise. 

```
\{{ eval example_var || "Not here" }}
```

#### Dates

Dates are represented in the following full format (some parts are optional, as we'll see):

```javascript
date["2018-03-20T22:55:00"]["Mon Jan 2 15:04:05 2006"]
```

The String between the first brackets is the actual date value, `20th of March 2018, 10:55:00AM` in this example,
and the second is the optional layout, which uses the fixed date-time `02nd of January 2006, 03:04:05PM` as a
layout-by-example (this is copied from the [Go Time Formatter](https://golang.org/pkg/time/#Time.Format)).

If not given, the following layout is used:

```
02 Jan 2006, 03:04 PM
```

The first String can also be `"now"` to evaluate the current time (when Magnanimous is running):

```javascript
date["now"]
```

When a path value is given to `date`, it resolves to the time the relevant file was last updated.

```javascript
\{{ define p1 path["path/to/file.txt"] }}
File path/to/file.txt was last updated on \{{ eval date[p1] }}
```

To make things clearer, here are a few examples:

| Date expression | Evaluated result |
|-----------------|------------------|
|`date["2018-03-20T22:55:00"]["Mon Jan 2 15:04:05 2006"]` | `Tue Mar 20 22:55:00 2018` |
|`date["2018-03-20T22:55"]["2 Jan 2006 03:04:05PM"]` | `20 Mar 2018 10:55:00PM` |
|`date["2018-03-20"]["2 Jan 2006"]` | `20 Mar 2018` |
|`date["2018-03-20"]` | `20 Mar 2018, 12:00 AM` |
|`date["2018-03-20T22:55"]` | `20 Mar 2018, 10:55 PM` |
|`date["now"]["2016"]` | `2019` (current year) |

{{ component /processed/components/_linked_header.html }}\
{{ define id "iterables" }}\
{{ define text "Iterables" }}\
{{ end }}

Iterables are collections of values which can be used with the [for](#for) instruction.

They can be of two types:

* Arrays.
* Paths to directories.

#### Arrays

Arrays have the following form:

```
[ expr1, expr2 ... ]
```

Examples:

```javascript
\{{ for number [1, 2, 3, 4, 5] }}\{{ end }}
\{{ for word [ "ABC" ] }}\{{ end }}
\{{ for section [ "Home", "About", "Docs" ] }}\{{ end }}
```

To iterate over an array defined elsewhere, the variable must be eval'd:

```javascript
\{{ define array [1,2,3] }}
\{{ for item eval array }}
    \{{ eval item }}
\{{ end }}
```

#### Paths to directories

A `for` expression may iterate over each file of a directory by declaring a _path_ to a directory as its iterable.

For example, suppose there is a directory with the following contents:

```
my-sections/
    /something.md
    /other.md
```

If each file [defines](#define) a variable named `section_name`, for example, then a summary file could contain
the following [for](#for) instruction to display the `section_name` of each file:

```
\{{ for section my-sections/ }}\\
  The name of the section is \{{ eval section_name }}!
\{{ end }}\\
```

See [Paths and Links](paths.html) for more details about paths.

{{ include _docs_footer.html }}
