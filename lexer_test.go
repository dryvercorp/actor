package actor

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CreateLexerReturnsNonNil(t *testing.T) {
	assert.NotNil(t, newLexer(&bytes.Buffer{}), "No new lexer supplied")
}

func Test_CreateLineReturnsNonNil(t *testing.T) {
	assert.NotNil(t, newLine(1, 1, ""), "No new token supplied")
}

func Test_LexerLoadsInputIntoTokens(t *testing.T) {

	var inputs = []struct {
		file          string
		expected_tree lexerTree
	}{
		{
			file: `
@tag1 @tag2
Actor: Valid actor
    Description and blurb... multi line 1 tab indent
    Some other line of blurb

    Goals:
        Goal number 1
        Goal number 2
    
    @tag3 @tag4
    Goal: Goal number 3`,

			expected_tree: lexerTree{
				&line{line: 2, column: 0, content: "@tag1 @tag2"},
				&line{line: 3, column: 0, content: "Actor: Valid actor", children: []*line{
					&line{line: 4, column: 4, content: "Description and blurb... multi line 1 tab indent"},
					&line{line: 5, column: 4, content: "Some other line of blurb"},
					&line{line: 7, column: 4, content: "Goals:", children: []*line{
						&line{line: 8, column: 8, content: "Goal number 1"},
						&line{line: 9, column: 8, content: "Goal number 2"},
					}},
					&line{line: 11, column: 4, content: "@tag3 @tag4"},
					&line{line: 12, column: 4, content: "Goal: Goal number 3"},
				}},
			},
		},
		{
			file: `

Actor: Valid actor


    Description and blurb... multi line 1 tab indent
    Some other line of blurb
    Another other line of blurb
`,

			expected_tree: lexerTree{
				&line{line: 3, column: 0, content: "Actor: Valid actor", children: []*line{
					&line{line: 6, column: 4, content: "Description and blurb... multi line 1 tab indent"},
					&line{line: 7, column: 4, content: "Some other line of blurb"},
					&line{line: 8, column: 4, content: "Another other line of blurb"},
				}},
			},
		},

		{
			file: `
Actor: Valid actor
 Description and blurb... multi line 1 tab indent
 Some other line of blurb
Another other line of blurb
`,

			expected_tree: lexerTree{
				&line{line: 2, column: 0, content: "Actor: Valid actor", children: []*line{
					&line{line: 3, column: 1, content: "Description and blurb... multi line 1 tab indent"},
					&line{line: 4, column: 1, content: "Some other line of blurb"},
				}},
				&line{line: 5, column: 0, content: "Another other line of blurb"},
			},
		},
	}

	for i, v := range inputs {
		lex := newLexer(bytes.NewBufferString(v.file))
		lines, err := lex.lex()
		assert.Nil(t, err, fmt.Sprintf("%d: Expected no error", i))
		compareLexerTrees(t, lines, v.expected_tree, 0)
	}
}

func Test_CanLoadValidActorDefinitionFromFile(t *testing.T) {

	var inputs = []struct {
		file          string
		expected_tree lexerTree
	}{
		{
			file: "examples/valid.actor",
			expected_tree: lexerTree{
				&line{line: 1, column: 0, content: "@tag1 @tag2"},
				&line{line: 2, column: 0, content: "Actor: Valid actor", children: []*line{
					&line{line: 3, column: 4, content: "Description and blurb... multi line 1 tab indent"},
					&line{line: 4, column: 4, content: "Some other line of blurb"},
					&line{line: 6, column: 4, content: "Goals:", children: []*line{
						&line{line: 7, column: 8, content: "Goal number 1"},
						&line{line: 8, column: 8, content: "Goal number 2"},
					}},
					&line{line: 10, column: 4, content: "@tag3 @tag4"},
					&line{line: 11, column: 4, content: "Goal: Goal number 3"},
				}},
			},
		},
	}

	for _, v := range inputs {
		file, err := os.Open("examples/valid.actor")
		assert.Nil(t, err)
		defer file.Close()

		lex := newLexer(file)
		lines, err := lex.lex()
		assert.Nil(t, err)

		compareLexerTrees(t, lines, v.expected_tree, 0)
	}
}

func compareLexerTrees(t *testing.T, a, b lexerTree, indent int) {
	assert.Equal(t, len(a), len(b), fmt.Sprintf("Index %d: Trees should be of equal length"))

	for i := 0; i < len(a); i++ {
		assert.Equal(t, a[i].line, b[i].line)
		assert.Equal(t, a[i].column, b[i].column)
		assert.Equal(t, a[i].content, b[i].content)

		compareLexerTrees(t, a[i].children, b[i].children, indent+1)
	}
}

func printTree(indent string, input lexerTree) {

	for _, v := range input {
		fmt.Printf("%s%+v\n", indent, v)
		printTree(indent+"-- ", v.children)
	}
}
