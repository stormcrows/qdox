package cmd

import (
	"bytes"
	"testing"
	"text/template"
)

func TestSearch(t *testing.T) {
	Tpl = template.Must(template.ParseGlob("../templates/*.gohtml"))
	app := NewApp()
	buf := new(bytes.Buffer)
	app.Writer = buf
	app.Run([]string{"qdox", "search", "../books/", "wild weekend"})

	expected := "92% \"../books/Grand Teton National Park.txt\"\n40% \"../books/Around the End - Ralph Henry Barbour.txt\"\n"
	out := buf.String()
	if expected != out {
		t.Errorf("expected: %q, got: %q", expected, out)
	}
}

func TestSearchWithDifferentN(t *testing.T) {
	Tpl = template.Must(template.ParseGlob("../templates/*.gohtml"))
	app := NewApp()
	buf := new(bytes.Buffer)
	app.Writer = buf
	app.Run([]string{"qdox", "search", "../books/", "wild weekend", "-n", "1"})

	expected := "92% \"../books/Grand Teton National Park.txt\"\n"
	out := buf.String()
	if expected != out {
		t.Errorf("expected: %q, got: %q", expected, out)
	}
}

func TestSearchWithDifferentThreshold(t *testing.T) {
	Tpl = template.Must(template.ParseGlob("../templates/*.gohtml"))
	app := NewApp()
	buf := new(bytes.Buffer)
	app.Writer = buf
	app.Run([]string{"qdox", "search", "../books/", "wild weekend", "-t", "0.5"})

	expected := "92% \"../books/Grand Teton National Park.txt\"\n"
	out := buf.String()
	if expected != out {
		t.Errorf("expected: %q, got: %q", expected, out)
	}
}
