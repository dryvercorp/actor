package actor

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	gherkin "github.com/cucumber/gherkin-go"
)

func Test_AWriterCanWriteTags(t *testing.T) {

	inputs := []struct {
		tags   []*gherkin.Tag
		indent int
		output string
	}{
		{
			tags: []*gherkin.Tag{
				{Name: "tag1"},
				{Name: "tag2"},
			},
			output: "@tag1 @tag2",
		},
		{
			tags: []*gherkin.Tag{
				{Name: "tag1"},
				{Name: "tag2"},
			},
			indent: 2,
			output: "        @tag1 @tag2",
		},
		{
			output: "",
		},
	}

	for _, input := range inputs {

		buf := &bytes.Buffer{}
		w := newWriter(buf)
		w.setIndentation(input.indent)

		assert.Nil(t, w.writeTags(input.tags))
		assert.Equal(t, buf.String(), input.output+"\n")
	}
}

func Test_AWriterCanWriteKeywords(t *testing.T) {

	inputs := []struct {
		keyword string
		value   string
		indent  int
		output  string
	}{
		{
			keyword: "Goal",
			value:   "some goal",
			output:  "Goal: some goal",
		},
		{
			keyword: "Goal",
			value:   "some goal",
			indent:  2,
			output:  "        Goal: some goal",
		},
	}

	for _, input := range inputs {

		buf := &bytes.Buffer{}
		w := newWriter(buf)
		w.setIndentation(input.indent)

		assert.Nil(t, w.writeKeyword(input.keyword, input.value))
		assert.Equal(t, buf.String(), input.output+"\n")
	}
}

func Test_AWriterCanWriteBlurbs(t *testing.T) {

	inputs := []struct {
		value  string
		indent int
		output string
	}{
		{
			value:  "some blurb",
			output: "some blurb",
		},
		{
			value:  "some blurb",
			indent: 2,
			output: "        some blurb",
		},
	}

	for _, input := range inputs {

		buf := &bytes.Buffer{}
		w := newWriter(buf)
		w.setIndentation(input.indent)

		assert.Nil(t, w.writeBlurb(input.value))
		assert.Equal(t, input.output+"\n", buf.String())
	}
}

func Test_AWriterCanWriteComments(t *testing.T) {

	inputs := []struct {
		value  string
		indent int
		output string
	}{
		{
			value:  "some comment",
			output: "# some comment",
		},
		{
			value:  "some comment",
			indent: 2,
			output: "        # some comment",
		},
	}

	for _, input := range inputs {

		buf := &bytes.Buffer{}
		w := newWriter(buf)
		w.setIndentation(input.indent)

		assert.Nil(t, w.writeComment(input.value))
		assert.Equal(t, input.output+"\n", buf.String())
	}
}
