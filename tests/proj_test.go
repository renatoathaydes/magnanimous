package tests

import (
	"github.com/renatoathaydes/magnanimous/mg"
	"os"
	"strings"
	"testing"
)

func TestProj0(t *testing.T) {
	dir := runMg(t, "test_proj_0")
	defer os.RemoveAll(dir)

	files, err := readAll(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 0 {
		t.Fatalf("Expected no output files, but got: %v", files)
	}
}

func TestProj1(t *testing.T) {
	dir := runMg(t, "test_proj_1")
	defer os.RemoveAll(dir)

	files, err := readAll(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 2 {
		t.Fatalf("Expected 2 output files, but got: %v", files)
	}

	assertFileContents(t, files, dir, "a.txt", "")
	assertFileContents(t, files, dir, "main.txt", "A and B:\n\n10\n20")
}

func TestProj2(t *testing.T) {
	dir := runMg(t, "test_proj_2")
	defer os.RemoveAll(dir)

	files, err := readAll(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 4 {
		t.Fatalf("Expected 4 output files, but got: %v", files)
	}

	assertFileContents(t, files, dir, "sections/section_a.txt", "This is section A.")
	assertFileContents(t, files, dir, "sections/section_b.txt", "This is section B.")
	assertFileContents(t, files, dir, "sections/section_c.txt", "This is section C.")
	assertFileContents(t, files, dir, "main.txt", "Main.\n\nSections:\n\n"+
		"  * Section A\n\n"+
		"  * Section B\n\n"+
		"  * Section C\n\n"+
		"End.")
}

func TestProj3(t *testing.T) {
	dir := runMg(t, "test_proj_3")
	defer os.RemoveAll(dir)

	files, err := readAll(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 4 {
		t.Fatalf("Expected 4 output files, but got: %v", files)
	}

	assertFileContents(t, files, dir, "index.html",
		`<link rel="stylesheet" href="/style.css">
<h1>Welcome</h1>

<p><div class ="interessant">This is a website.</div>
My blog posts:</p>

<ul>
<li><p>01 Jan 2018 - My first blog post</p></li>

<li><p>05 July 2018 - One more blog post</p></li>
</ul>

<h4>the footer</h4>
`)

	assertFileContents(t, files, dir, "posts/p1.html",
		`<link rel="stylesheet" href="/style.css">
<h2>Post 1</h2>

<p>Hello.</p>

<blockquote>
<p>A note.</p>
</blockquote>

<p>Bye.</p>
`)

	assertFileContents(t, files, dir, "posts/p2.html",
		`<link rel="stylesheet" href="/style.css">
<h2>Post 2</h2>

<p>Short one.</p>
`)

	assertFileContents(t, files, dir, "style.css", `h2 {
    color: blue;
}

.interessant {
    font-weight: bolder;
}`)

}

func TestProj4(t *testing.T) {
	dir := runMg(t, "test_proj_4")
	defer os.RemoveAll(dir)

	files, err := readAll(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 7 {
		t.Fatalf("Expected 7 output files, but got: %v", files)
	}

	assertFileContents(t, files, dir, "index.html",
		`<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0, shrink-to-fit=no">
    <title>My blog</title>
    <link rel="stylesheet" href="style.css">
</head>
<body>
<h1>Welcome</h1>

<p><div class ="interessant">This is a website.</div>
My blog posts:</p>

<ul>
<li><p>2019-02-23 - <a href="posts/brocolli.html">Broccoli</a></p></li>

<li><p>2019-01-31 - <a href="posts/capsicum.html">Capsicum</a></p></li>

<li><p>2018-08-23 - <a href="posts/potato.html">Potatoes</a></p></li>

<li><p>2018-07-05 - <a href="posts/second.html">One more blog post</a></p></li>
</ul>
<h4>the footer</h4>
</body>
</html>`)

	assertFileContents(t, files, dir, "posts/first_post.html",
		`<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0, shrink-to-fit=no">
    <title>My first blog post</title>
    <link rel="stylesheet" href="style.css">
</head>
<body>
<h2>Post 1</h2>

<p>Hello.</p>

<blockquote>Note: This is a note.</blockquote>

<p>Bye.</p>
<h4>the footer</h4>
</body>
</html>`)

	assertFileContents(t, files, dir, "posts/second.html",
		`<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0, shrink-to-fit=no">
    <title>One more blog post</title>
    <link rel="stylesheet" href="style.css">
</head>
<body>
<h2>Post 2</h2>

<p>Short one.</p>
<h4>the footer</h4>
</body>
</html>`)

	assertFileContents(t, files, dir, "posts/potato.html",
		`<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0, shrink-to-fit=no">
    <title>Potatoes</title>
    <link rel="stylesheet" href="style.css">
</head>
<body>
<h2>Potatoes</h2>

<p>Potatoes are nice.</p>
<h4>the footer</h4>
</body>
</html>`)

	assertFileContents(t, files, dir, "posts/capsicum.html",
		`<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0, shrink-to-fit=no">
    <title>Capsicum</title>
    <link rel="stylesheet" href="style.css">
</head>
<body>
<h2>Capsicum</h2>

<p>Capsicum is good.</p>
<h4>the footer</h4>
</body>
</html>`)

	assertFileContents(t, files, dir, "posts/brocolli.html",
		`<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0, shrink-to-fit=no">
    <title>Broccoli</title>
    <link rel="stylesheet" href="style.css">
</head>
<body>
<h2>Broccoli and you</h2>

<p>You should eat more broccoli.</p>
<h4>the footer</h4>
</body>
</html>`)

	assertFileContents(t, files, dir, "style.css", `h2 {
    color: blue;
}

body {
    background-color: aliceblue;
}

.interessant {
    font-weight: bolder;
}`)

}

func TestProj5(t *testing.T) {
	dir := runMg(t, "test_proj_5")
	defer os.RemoveAll(dir)

	files, err := readAll(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 4 {
		t.Fatalf("Expected 4 output files, but got: %v", files)
	}

	assertFileContents(t, files, dir, "a.txt", "")
	assertFileContents(t, files, dir, "main.txt", "A and B:\n\n"+
		"/my-website/10\n"+
		"/my-website/20")
	assertFileContents(t, files, dir, "folder/example.txt", "Full path: /my-website/folder/example.txt")
	assertFileContents(t, files, dir, "scopes.txt", "Base URL: /my-website/\n\n"+
		"/other-website/A\n"+
		"/other-website/<nil>\n\n"+
		"File sees /my-website/\n"+
		"Full path: /other-website/folder/example.txt\n"+
		"After unset, base URL: /my-website/")
}

func TestProj6(t *testing.T) {
	dir := runMg(t, "test_proj_6")
	defer os.RemoveAll(dir)

	files, err := readAll(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 2 {
		t.Fatalf("Expected 2 output files, but got: %v", files)
	}

	assertFileContents(t, files, dir, "index.html", `<h2>Simple Component Example</h2>
<div class="wrapper">
    Hello components
</div>`)

	assertFileContents(t, files, dir, "example/properties.html", `<div>Component example</div>
<h1>Component with properties</h1>

<p>Text: This is some text
Number: 23</p>
`)
}

// Initial results:
// 789098 ns/op
// 801477 ns/op
//
// After avoiding expression evaluation for repeated inclusion chain
// 768579 ns/op
// 766237 ns/op
//
func BenchmarkProject4(b *testing.B) {
	for n := 0; n < b.N; n++ {
		benchmarkProject(b, "test_proj_4")
	}
}

func benchmarkProject(b *testing.B, project string) {
	mag := mg.Magnanimous{SourcesDir: project}
	webFiles, err := mag.ReadAll()

	if err != nil {
		b.Fatal(err)
	}

	for _, webFile := range webFiles.WebFiles {
		var w strings.Builder
		err = webFile.Write(&w, webFiles, nil)

		if err != nil {
			b.Fatal(err)
		}

		if len(w.String()) < 10 {
			b.Errorf("Expected a String longer than 10, got %s", w.String())
		}
	}
}
