package actor

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewParserReturnsANonNilObject(t *testing.T) {
	assert.NotNil(t, NewParser(bytes.NewBuffer(nil)))
}

func Test_NewFileParserCanLoadAFile(t *testing.T) {
	parser, err := NewFileParser("examples/valid.actor")
	assert.NotNil(t, parser)
	assert.Nil(t, err)
}

func Test_ItCanParseAnActorFile(t *testing.T) {

	file := `
# This is a comment
@tag1 @tag2
Actor: Valid actor
    Description and blurb... multi line 1 tab indent
    Some other line of blurb

    Goals:
        Goal number 1
        Goal number 2
    
    @tag3 @tag4
    Goal: Goal number 3
`

	parser := NewParser(bytes.NewBufferString(file))
	_, err := parser.Parse()
	assert.Nil(t, err)
}
