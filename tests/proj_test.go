package tests

import (
	"os"
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
		`<link rel="stylesheet" href="/style.css"><h1>Welcome</h1>

<p><div class ="interessant">This is a website.</div>
My blog posts:</p>
<ul>
<li><p>01 Jan 2018 - My first blog post</p></li>

<li><p>05 July 2018 - One more blog post</p></li>
</ul>
<h4>the footer</h4>
`)

	assertFileContents(t, files, dir, "posts/p1.html",
		`<link rel="stylesheet" href="/style.css"><h2>Post 1</h2>

<p>Hello.</p>

<blockquote>
<p>A note.</p>
</blockquote>

<p>Bye.</p>
`)

	assertFileContents(t, files, dir, "posts/p2.html",
		`<link rel="stylesheet" href="/style.css"><h2>Post 2</h2>

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
		`<html>
<head>
    <title>My blog</title>
    <link rel="stylesheet" href="/style.css">
</head>
<body><h1>Welcome</h1>

<p><div class ="interessant">This is a website.</div>
My blog posts:</p>
<ul>
<li><p>2019-02-23 - Broccoli</p></li>

<li><p>2019-01-31 - Capsicum</p></li>

<li><p>2018-08-23 - Potatoes</p></li>

<li><p>2018-07-05 - One more blog post</p></li>
</ul>
<h4>the footer</h4>
</body>
</html>`)

	assertFileContents(t, files, dir, "posts/first_post.html",
		`<html>
<head>
    <title>My first blog post</title>
    <link rel="stylesheet" href="/style.css">
</head>
<body><h2>Post 1</h2>

<p>Hello.</p>
<blockquote>Note: This is a note.</blockquote>
<p>Bye.</p>
<h4>the footer</h4>
</body>
</html>`)

	assertFileContents(t, files, dir, "posts/second.html",
		`<html>
<head>
    <title>One more blog post</title>
    <link rel="stylesheet" href="/style.css">
</head>
<body><h2>Post 2</h2>

<p>Short one.</p>
<h4>the footer</h4>
</body>
</html>`)

	assertFileContents(t, files, dir, "posts/potato.html",
		`<html>
<head>
    <title>Potatoes</title>
    <link rel="stylesheet" href="/style.css">
</head>
<body><h2>Potatoes</h2>

<p>Potatoes are nice.</p>
<h4>the footer</h4>
</body>
</html>`)

	assertFileContents(t, files, dir, "posts/capsicum.html",
		`<html>
<head>
    <title>Capsicum</title>
    <link rel="stylesheet" href="/style.css">
</head>
<body><h2>Capsicum</h2>

<p>Capsicum is good.</p>
<h4>the footer</h4>
</body>
</html>`)

	assertFileContents(t, files, dir, "posts/brocolli.html",
		`<html>
<head>
    <title>Broccoli</title>
    <link rel="stylesheet" href="/style.css">
</head>
<body><h2>Broccoli and you</h2>

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
