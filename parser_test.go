package actor

import (
	"bytes"
	"fmt"
	"testing"

	gherkin "github.com/cucumber/gherkin-go"

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

	var inputs = []struct {
		file  string
		err   error
		actor *Actor
	}{
		{
			file: `
# This is a comment
@tag1 @tag2
Actor: Valid actor
    Description and blurb... multi line 1 tab indent
    Some other line of blurb

    @tag3 @tag4
    Goals:
        Goal number 1
        Goal number 2
    
    @tag5 @tag6
    Goal: Goal number 3
`,
			actor: &Actor{
				Name: "Valid actor",
				Tags: []*gherkin.Tag{
					{Name: "tag1"},
					{Name: "tag2"},
				},
				Blurb: []string{
					"Description and blurb... multi line 1 tab indent",
					"Some other line of blurb",
				},
				Goals: []*Goal{
					{
						Name: "Goal number 1",
						Tags: []*gherkin.Tag{
							{Name: "tag3"},
							{Name: "tag4"},
						},
					},
					{
						Name: "Goal number 2",
						Tags: []*gherkin.Tag{
							{Name: "tag3"},
							{Name: "tag4"},
						},
					},
					{
						Name: "Goal number 3",
						Tags: []*gherkin.Tag{
							{Name: "tag5"},
							{Name: "tag6"},
						},
					},
				},
			},
		},
		{
			file: `@tag @ tag`,
			err:  fmt.Errorf("[Line 0001:00] Tag '@' (#2 on the line) is not valid"),
		},

		{
			file: `Actor:`,
			err:  fmt.Errorf("[Line 0001:00] Actor keyword must be followed by an actor name"),
		},

		{
			file: `
Actor: Some actor
Actor: Some other actor`,
			err: fmt.Errorf("[Line 0003:00] Only one actor definition is permitted per file (other actor 'Some actor' : [Line 0002:00])"),
		},
	}

	for _, input := range inputs {

		parser := NewParser(bytes.NewBufferString(input.file))
		actor, err := parser.Parse()

		assert.Equal(t, input.err, err)

		if input.actor != nil {
			compareActors(t, input.actor, actor)
		}
	}
}
