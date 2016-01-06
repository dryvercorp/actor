package actor

import gherkin "github.com/cucumber/gherkin-go"

type Actor struct {
	gherkin.Node
	Tags  []*gherkin.Tag `json:"tags"`
	Name  string         `json:"name"`
	Blurb []string       `json:"blurb,omitempty"`
	Goals []*Goal        `json:"goals,omitempty"`
}

type Goal struct {
	gherkin.Node
	Tags []*gherkin.Tag `json:"tags"`
	Name string         `json:"name"`
}

func NewActor() *Actor {
	actor := Actor{}

	actor.Tags = make([]*gherkin.Tag, 0)
	actor.Blurb = make([]string, 0)
	actor.Goals = make([]*Goal, 0)

	return &actor
}
