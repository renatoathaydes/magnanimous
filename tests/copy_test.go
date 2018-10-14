package tests

import (
	"github.com/renatoathaydes/magnanimous/mg"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestCopy(t *testing.T) {
	f, err := ioutil.TempFile("", "copy_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	text := "hello Magnanimous!"
	_, err = f.Write([]byte(text))
	if err != nil {
		t.Fatal(err)
	}

	result := mg.Copy(f.Name(), "b", true)

	if result.BasePath != "b" {
		t.Errorf("Expected basePath 'b' but was '%s'", result.BasePath)
	}

	if result.Context != nil {
		t.Errorf("Expected nil context but was %v", result.Context)
	}

	contents := result.Processed.Contents

	if len(contents) != 1 {
		t.Errorf("Contents does not have length 1: %v", contents)
	}

	w := strings.Builder{}
	m := mg.WebFilesMap{}
	me := contents[0].Write(&w, m, nil)

	if me != nil {
		t.Fatal(me)
	}

	if w.String() != text {
		t.Errorf("Expected copied contents to be '%s', but was '%s'", text, w.String())
	}
}
