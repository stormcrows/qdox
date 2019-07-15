package nlp

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

// Document holds content of the file and its path
type Document struct {
	content string
	path    string
}

// Corpus is a list of documents
type Corpus struct {
	documents []Document
}

// NewCorpus returns an empty corpus
func NewCorpus() Corpus {
	return Corpus{make([]Document, 0)}
}

// Release contents from memory after training
func (c *Corpus) Release() {
	for i := 0; i < len(c.documents); i++ {
		c.documents[i].content = ""
	}
}

// Load walks given path recursively and adds documents to the corpus, that match the pattern
func (c *Corpus) Load(path string, pattern *regexp.Regexp) error {
	n, err := CountDocuments(path, pattern)
	if err != nil {
		return err
	}

	c.documents = make([]Document, n)
	i := 0

	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if ok := pattern.MatchString(path); !ok {
			return nil
		}

		content, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		c.documents[i] = Document{string(content), path}
		i++

		return nil
	})
}

// CountDocuments returns number of documents under given path for given pattern
func CountDocuments(path string, pattern *regexp.Regexp) (int, error) {
	n := 0
	return n, filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if ok := pattern.MatchString(path); ok {
			n++
		}

		return nil
	})
}

// Contents returns a list of loaded document's contents as strings
func (c *Corpus) Contents() []string {
	contents := make([]string, len(c.documents))
	for i := 0; i < len(c.documents); i++ {
		contents[i] = c.documents[i].content
	}
	return contents
}

// Paths is returning a list of paths to loaded documents
func (c *Corpus) Paths() []string {
	paths := make([]string, len(c.documents))
	for i := 0; i < len(c.documents); i++ {
		paths[i] = c.documents[i].path
	}
	return paths
}

// GetPath returns path for given document's index
func (c *Corpus) GetPath(i int) string {
	return c.documents[i].path
}
