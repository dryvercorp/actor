package actor

import (
	"testing"

	gherkin "github.com/cucumber/gherkin-go"

	"github.com/stretchr/testify/assert"
)

func Test_CreateActorReturnsNonNil(t *testing.T) {
	assert.NotNil(t, NewActor(), "No new actor supplied")
}

func compareActors(t *testing.T, a1, a2 *Actor) {
	assert.NotNil(t, a1)
	assert.NotNil(t, a2)
	assert.Equal(t, a1.Name, a2.Name)

	// Check general tags
	assert.Equal(t, len(a1.Tags), len(a2.Tags))

	for i := 0; i < len(a2.Tags); i++ {
		assert.Equal(t, a1.Tags[i].Name, a2.Tags[i].Name)
	}

	// Check blurb
	assert.Equal(t, a1.Blurb, a2.Blurb)

	// Check goals
	assert.Equal(t, len(a1.Goals), len(a2.Goals))

	for i := 0; i < len(a2.Goals); i++ {
		assert.Equal(t, a1.Goals[i].Name, a2.Goals[i].Name)
		assert.Equal(t, len(a1.Goals[i].Tags), len(a2.Goals[i].Tags))

		for j := 0; j < len(a1.Goals[i].Tags); j++ {
			assert.Equal(t, a1.Goals[i].Tags[j].Name, a2.Goals[i].Tags[j].Name)
		}
	}
}

func newMockActor() *Actor {
	actor := NewActor()

	actor.Name = "Mock actor"

	actor.Tags = []*gherkin.Tag{
		{Name: "tag1"},
		{Name: "tag2"},
	}

	actor.Blurb = []string{
		"Blurb line 1",
		"BLurb line 2",
	}

	actor.Goals = []*Goal{
		{
			Name: "Goal 1",
			Tags: []*gherkin.Tag{
				{Name: "tag3"},
				{Name: "tag4"},
			},
		},
		{
			Name: "Goal 2",
		},
		{
			Name: "Goal 3",
		},
	}

	return actor
}
