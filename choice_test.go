package strgen

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChoice(t *testing.T) {
	channel, amount, err := GenerateStrings("\\(a|b)")
	assert.NoError(t, err)
	assert.Equal(t, int64(2), amount, "should produce two results")
	assert.Equal(t, "a", <-channel, "should be 'a'")
	assert.Equal(t, "b", <-channel, "should be 'b'")
	_, amount, err = GenerateStrings("\\(|foo|)")
	assert.NoError(t, err)
	assert.Equal(t, int64(3), amount, "should produce three results")
}

func ExampleChoiceIterator() {
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
	assert.Error(t, err)
	_, _, err = GenerateStrings("\\(\\[")
	assert.Error(t, err)
}

func BenchmarkConfigTwoWayChoice(b *testing.B) {
	_BenchmarkConfig("\\(a|b)", b)
}

func BenchmarkConfigTenWayChoice(b *testing.B) {
	_BenchmarkConfig("\\(1|2|3|4|5|6|7|8|9|10)", b)
}
