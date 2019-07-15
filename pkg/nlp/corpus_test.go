package nlp

import (
	"regexp"
	"testing"
)

func TestLoad(t *testing.T) {
	c := NewCorpus()
	r := regexp.MustCompile("\\.txt")
	err := c.Load("../../docs", r)
	if err != nil {
		t.Errorf("error reading folder %s", err.Error())
	}

	l := len(c.documents)
	expected := 4

	if l != expected {
		t.Errorf("expected to read %d documents, got: %d", expected, l)
	}

	for _, doc := range c.documents {
		if doc.content == "" || doc.path == "" {
			t.Errorf("expecting content and path to be set for: %q", doc.path)
		}
	}
}

func TestRelease(t *testing.T) {
	c := NewCorpus()
	r := regexp.MustCompile("\\.txt")
	err := c.Load("../../docs", r)
	if err != nil {
		t.Errorf("error reading folder %s", err.Error())
	}

	c.Release()

	for _, doc := range c.documents {
		if doc.content != "" {
			t.Errorf("expecting content to be empty for: %q", doc.path)
		}
		if doc.path == "" {
			t.Errorf("expecting path to be not empty but it is")
		}
	}
}

func TestCountDocuments(t *testing.T) {
	r := regexp.MustCompile("\\.txt")
	n, err := CountDocuments("../../docs", r)
	if err != nil {
		t.Errorf("error counting documents: %s", err.Error())
	}

	if n != 4 {
		t.Errorf("expected 4 documents, but got: %d", n)
	}
}
