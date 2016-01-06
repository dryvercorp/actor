package actor

import (
	"bytes"
	"fmt"
	"io"
	"os"

	gherkin "github.com/cucumber/gherkin-go"
)

type Parser interface {
	Parse() (*Actor, error)
}

type parser struct {
	reader      io.Reader
	actor       *Actor
	pendingTags []*gherkin.Tag
}

func NewParser(r io.Reader) Parser {
	return &parser{
		reader:      r,
		pendingTags: make([]*gherkin.Tag, 0),
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

	p.resetTags()

	lex := newLexer(p.reader)
	tkn := newTokeniser()

	tree, lex_err := lex.lex()

	if lex_err != nil {
		return nil, fmt.Errorf("Lexer error: %s", err)
	}

	if err := p.parseTree(tree, tkn); err != nil {
		return nil, err
	}

	return p.actor, nil
}

func (p *parser) resetTags() {
	p.pendingTags = make([]*gherkin.Tag, 0)
}

func (p *parser) addTag(l *line, name string) {

	p.pendingTags = append(
		p.pendingTags,
		&gherkin.Tag{
			Location: &gherkin.Location{
				Line:   l.line,
				Column: l.column,
			},
			Name: name,
		},
	)
}

func (p *parser) addPendingTagsToList(list *[]*gherkin.Tag) {
	*list = append(*list, p.pendingTags...)
	p.resetTags()
}

func (p *parser) err(branch *line, e string, args ...interface{}) error {
	return fmt.Errorf("[Line %04d:%02d] %s", branch.line, branch.column, fmt.Sprintf(e, args...))
}

func (p *parser) parseTree(tree lexerTree, tkn *tokeniser) error {
	for _, branch := range tree {

		tokens, err := tkn.tokenise(branch)

		if err != nil {
			return p.err(branch, err.Error())
		}

		for _, token := range tokens {

			switch token.kind {

			case token_tag:
				p.addTag(branch, token.content)

			case token_actorDefinition:
				if err := p.parseActorDefinition(branch, token, tkn); err != nil {
					return err
				}

			case token_goal:
				if err := p.parseGoal(branch, token, tkn); err != nil {
					return err
				}

			case token_goals:
				if err := p.parseGoals(branch, token, tkn); err != nil {
					return err
				}

			case token_text:
				if err := p.parseText(branch, token, tkn); err != nil {
					return err
				}

			default:
				return p.err(branch, "Parse error: %s", branch.content)
			}
		}
	}

	return nil
}

func (p *parser) parseActorDefinition(branch *line, t token, tkn *tokeniser) error {

	if t.content == "" {
		return p.err(branch, "Actor keyword must be followed by an actor name")
	}

	if p.actor != nil {
		return p.err(branch, "Only one actor definition is permitted per file (other actor '%s' : [Line %04d:%02d])", p.actor.Name, p.actor.Location.Line, p.actor.Location.Column)
	}

	p.actor = NewActor()
	p.actor.Name = t.content

	p.actor.Location = &gherkin.Location{
		Line:   branch.line,
		Column: branch.column,
	}

	p.addPendingTagsToList(&p.actor.Tags)

	return p.parseTree(branch.children, tkn)
}

func (p *parser) parseGoal(branch *line, t token, tkn *tokeniser) error {

	if p.actor == nil {
		return p.err(branch, "Goal keyword outside of actor context")
	}

	if t.content == "" {
		return p.err(branch, "Goal keyword must be followed by a goal name")
	}

	goal := &Goal{Name: t.content}
	goal.Location = &gherkin.Location{
		Line:   branch.line,
		Column: branch.column,
	}

	p.addPendingTagsToList(&goal.Tags)

	p.actor.Goals = append(p.actor.Goals, goal)

	return p.parseTree(branch.children, tkn)
}

func (p *parser) parseGoals(branch *line, t token, tkn *tokeniser) error {

	if p.actor == nil {
		return p.err(branch, "Goals keyword outside of actor context")
	}

	for _, goalDef := range branch.children {

		tokens, err := tkn.tokenise(goalDef)

		if err != nil {
			return p.err(goalDef, err.Error())
		}

		for _, t := range tokens {

			if t.kind != token_text {
				return p.err(goalDef, "Unexpected %s in goal list", t.kind)
			}

			goal := &Goal{Name: t.content}
			goal.Location = &gherkin.Location{
				Line:   branch.line,
				Column: branch.column,
			}

			goal.Tags = append(goal.Tags, p.pendingTags...)

			p.actor.Goals = append(p.actor.Goals, goal)
		}
	}

	p.resetTags()

	return nil
}

func (p *parser) parseText(branch *line, t token, tkn *tokeniser) error {

	if p.actor == nil {
		return p.err(branch, "Blurb text outside of actor context")
	}

	p.actor.Blurb = append(p.actor.Blurb, t.content)

	return p.parseTree(branch.children, tkn)
}
