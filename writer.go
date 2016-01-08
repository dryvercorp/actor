package actor

import (
	"fmt"
	"io"
	"strings"

	gherkin "github.com/cucumber/gherkin-go"
)

type writer struct {
	writer      io.Writer
	indentation int
}

func newWriter(w io.Writer) *writer {
	return &writer{writer: w}
}

func (w *writer) indent() {
	w.indentation++
}

func (w *writer) unindent() {

	w.indentation--

	if w.indentation < 0 {
		w.indentation = 0
	}
}

func (w *writer) setIndentation(i int) {
	w.indentation = i
}

func (w *writer) indentString() string {
	return strings.Repeat(" ", 4*w.indentation)
}

func (w *writer) newLine() error {
	_, err := w.writer.Write([]byte("\n"))
	return err
}

func (w *writer) writeLine(b []byte) error {

	if _, err := w.writer.Write(b); err != nil {
		return err
	}

	if err := w.newLine(); err != nil {
		return err
	}

	return nil
}

func (w *writer) writeTags(tags []*gherkin.Tag) error {

	tagString := w.indentString()

	for _, tag := range tags {
		tagString += fmt.Sprintf("@%s ", tag.Name)
	}

	return w.writeLine([]byte(strings.TrimRight(tagString, " ")))
}

func (w *writer) writeKeyword(keyword, value string) error {

	keywordString := fmt.Sprintf("%s%s: %s", w.indentString(), keyword, value)

	return w.writeLine([]byte(strings.TrimRight(keywordString, " ")))
}

func (w *writer) writeBlurb(value string) error {

	blurbString := fmt.Sprintf("%s%s", w.indentString(), value)

	return w.writeLine([]byte(blurbString))
}

func (w *writer) writeComment(value string) error {

	commentString := fmt.Sprintf("%s# %s", w.indentString(), value)

	return w.writeLine([]byte(commentString))
}
