package actor

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_StringsTokeniseAsExpected(t *testing.T) {

	var inputs = []struct {
		line   string
		tokens []token
		err    error
	}{
		///////////////////////////////
		// Tags
		///////////////////////////////

		{
			// Normal tags
			line: "@tag1 @tag2",
			tokens: []token{
				{kind: token_tag, content: "tag1"},
				{kind: token_tag, content: "tag2"},
			},
		},
		{
			// One tag and a comment
			line: "@tag1 # @tag2",
			tokens: []token{
				{kind: token_tag, content: "tag1"},
			},
		},
		{
			// Invalid tags
			line: "@",
			err:  fmt.Errorf("Tag '@' (#1 on the line) is not valid"),
		},
		{
			// Invalid tags
			line: "@1",
			err:  fmt.Errorf("Tag '@1' (#1 on the line) is not valid"),
		},
		{
			// Invalid tags
			line: "@_",
			err:  fmt.Errorf("Tag '@_' (#1 on the line) is not valid"),
		},
		{
			// Invalid tags
			line: "@tag*1",
			err:  fmt.Errorf("Tag '@tag*1' (#1 on the line) is not valid"),
		},

		///////////////////////////////
		// Keywords
		///////////////////////////////

		{
			line: "Actor:",
			tokens: []token{
				{kind: token_actorDefinition, content: ""},
			},
		},
		{
			line: "Actor: #this is a comment",
			tokens: []token{
				{kind: token_actorDefinition, content: ""},
			},
		},
		{
			line: "Actor: Some actor",
			tokens: []token{
				{kind: token_actorDefinition, content: "Some actor"},
			},
		},

		{
			line: "NotAToken: Someactor",
			err:  fmt.Errorf("Unrecognised keyword 'NotAToken'"),
		},

		///////////////////////////////
		// Text
		///////////////////////////////

		{
			line: "This is text",
			tokens: []token{
				{kind: token_text, content: "This is text"},
			},
		},

		{
			line: "This is text #this is comment",
			tokens: []token{
				{kind: token_text, content: "This is text"},
			},
		},

		///////////////////////////////
		// Empty
		///////////////////////////////

		{
			line: "",
		},

		{
			line: "#This is a line comment",
		},

		{
			line: "    #This is a line comment with left padding",
		},
	}

	for _, input := range inputs {
		tkn := newTokeniser()
		tokens, err := tkn.tokenise(&line{content: lineContent(input.line)})
		assert.Equal(t, input.err, err)
		assert.Equal(t, input.tokens, tokens)
	}
}
