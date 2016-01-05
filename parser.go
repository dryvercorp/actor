package actor

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

type Parser interface {
	Parse() (*Actor, error)
}

type parser struct {
	reader io.Reader
}

func NewParser(r io.Reader) Parser {
	return &parser{
		reader: r,
	}
}

func NewFileParser(path string) (Parser, error) {

	file, err := os.Open(path)
	defer file.Close()

	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(nil)

	if _, err := io.Copy(buf, file); err != nil {
		return nil, err
	}

	return NewParser(buf), nil
}

func (p *parser) Parse() (actor *Actor, err error) {

	actor = NewActor()

	lex := newLexer(p.reader)
	tkn := newTokeniser()

	tree, lex_err := lex.lex()

	if lex_err != nil {
		return nil, fmt.Errorf("Lexer error: %s", err)
	}

	for _, branch := range tree {
		tokens, err := tkn.tokenise(branch)

		if err != nil {
			return nil, fmt.Errorf("[Line %04d:%02d] %s", branch.line, branch.column, err)
		}

		// Skip comment lines
		if len(tokens) == 0 {
			continue
		}
	}

	return
}
