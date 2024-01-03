package greeter_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vidarandrebo/once-tree/pkg/greeter"
)

var greetingText = "Hello there"

func TestGreeting(t *testing.T) {

	g := greeter.NewGreeting()
	assert.Equal(t, g.Text, greetingText)
}
