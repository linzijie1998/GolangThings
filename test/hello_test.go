package hello

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSayHello(t *testing.T) {
	output := SayHello()
	expectOutput := "Hello!"
	assert.Equal(t, expectOutput, output)
}
