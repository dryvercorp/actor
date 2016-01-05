package actor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CreateActorReturnsNonNil(t *testing.T) {
	assert.NotNil(t, NewActor(), "No new actor supplied")
}
