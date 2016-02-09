package strgen

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChoice(t *testing.T) {
	channel, amount, err := GenerateStrings("\\(a|b)")
	assert.Nil(t, err)
	assert.Equal(t, int64(2), amount, "should produce two results")
	assert.Equal(t, "a", <-channel, "should be 'a'")
	assert.Equal(t, "b", <-channel, "should be 'b'")
	_, amount, err = GenerateStrings("\\(|foo|)")
	assert.Nil(t, err)
	assert.Equal(t, int64(3), amount, "should produce three results")
}

func ExampleChoice() {
	c, _, _ := GenerateStrings("\\(foo|bar|baz)")
	for s := range c {
		fmt.Println(s)
	}
	// Output:
	// foo
	// bar
	// baz
}

func TestInvalidChoice(t *testing.T) {
	var err error
	_, _, err = GenerateStrings("\\(a|b")
	assert.NotNil(t, err)
	_, _, err = GenerateStrings("\\(\\[")
	assert.NotNil(t, err)
}
