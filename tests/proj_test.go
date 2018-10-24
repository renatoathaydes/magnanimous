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
