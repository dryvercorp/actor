package actor

import (
	"fmt"
	"regexp"
	"strings"
)

//go:generate stringer -type=tokenKind
type tokenKind int

var commentMatcher = regexp.MustCompile(`#.+$`)
var tagMatcher = regexp.MustCompile(`^@([a-zA-Z][a-zA-Z0-9_-]*)$`)
var keywordMatcher = regexp.MustCompile(`^(.+):\s?(.+)?$`)

const (
	token_comment tokenKind = iota
	token_tag
	token_actorDefinition
	token_text
	token_goals
	token_goal
)

var tokenKindsByString = map[string]tokenKind{
	"actor": token_actorDefinition,
	"goals": token_goals,
	"goal":  token_goal,
}

type token struct {
	kind    tokenKind
	content string
}

type tokeniser struct{}

func newTokeniser() *tokeniser {
	return &tokeniser{}
}

func (t *tokeniser) tokenise(l *line) (tokens []token, err error) {

	// Remove comments
	content := strings.Trim(string(commentMatcher.ReplaceAll([]byte(l.content), []byte(""))), " \t")

	// Check for empty
	if content == "" {
		return
	}

	// Tag lines start with a @
	if content[0] == '@' {
		return t.tokeniseTags(lineContent(content))
	}

	// Check for Something:
	if keywordMatcher.Match([]byte(content)) {
		return t.tokeniseKeyword(lineContent(content))
	}

	// Assume the result is text
	return []token{
		token{kind: token_text, content: content},
	}, nil
}

func (t *tokeniser) tokeniseTags(content lineContent) (tokens []token, err error) {
	tags := strings.Fields(string(content))

	for i, v := range tags {
		if matched := tagMatcher.Find([]byte(v)); matched != nil {
			tokens = append(tokens, token{kind: token_tag, content: string(matched[1:])})
		} else {
			return nil, fmt.Errorf("Tag '%s' (#%d on the line) is not valid", v, i+1)
		}
	}

	return
}

func (t *tokeniser) tokeniseKeyword(content lineContent) (tokens []token, err error) {

	terms := keywordMatcher.FindStringSubmatch(string(content))
	typ, ok := tokenKindsByString[strings.ToLower(terms[1])]

	if !ok {
		return nil, fmt.Errorf("Unrecognised keyword '%s'", terms[1])
	}

	tokens = append(tokens, token{kind: typ, content: terms[2]})

	return
}
