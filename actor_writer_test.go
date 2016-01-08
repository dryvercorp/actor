package actor

import (
	"bytes"
	"io/ioutil"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ItCanWriteAnActorToAWriter(t *testing.T) {

	actor := newMockActor()
	writer := &bytes.Buffer{}
	expected := `@tag1 @tag2
Actor: Mock actor
    Blurb line 1
    BLurb line 2

    @tag3 @tag4
    Goal: Goal 1

    Goals:
        Goal 2
        Goal 3
`

	assert.Nil(t, actor.Write(writer))
	assert.Equal(t, expected, writer.String())
}

func Test_ItCanReadAnActorAfterWriting(t *testing.T) {

	actor := newMockActor()
	buf := &bytes.Buffer{}

	assert.Nil(t, actor.Write(buf))

	parser := NewParser(buf)

	read_actor, err := parser.Parse()
	assert.Nil(t, err)

	compareActors(t, read_actor, actor)
}

func Test_ItCanWriteToAFile(t *testing.T) {

	actor := newMockActor()
	f, err := ioutil.TempFile("", "go-actor-tesr")

	assert.Nil(t, err)

	name := f.Name()

	defer syscall.Unlink(name)

	assert.Nil(t, actor.WriteToFile(name))

	parser, err := NewFileParser(name)
	assert.Nil(t, err)

	read_actor, err := parser.Parse()
	assert.Nil(t, err)

	compareActors(t, read_actor, actor)
}
