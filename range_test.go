package strgen

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRangeAmount(t *testing.T) {
	assertAmount := func(s string, amount int64) {
		_, actual, err := GenerateStrings(s)
		assert.NoError(t, err)
		assert.Equal(t, amount, actual, fmt.Sprintf("%v should produce %v results", s, amount))
	}
	// Finite, Integer ranges
	assertAmount("\\[0..3]", 4)
	assertAmount("\\[2..4]", 3)
	assertAmount("\\[5..2]", 4)

	// Finite, FP ranges
	assertAmount("\\[1.5..2]", 1)
	assertAmount("\\[1.5..3]", 2)

	// Finite, Integer ranges w/ step
	assertAmount("\\[0..2..3]", 2)
	assertAmount("\\[-2..0.5..0]", 5)
	assertAmount("\\[5..-1..2]", 4)

	// Infinite
	assertAmount("\\[0..]", -1)
	assertAmount("\\[-42..]", -1)
}

func ExampleRangeIterator() {
	c, _, _ := GenerateStrings("\\[0..0.5..2]")
	for s := range c {
		fmt.Println(s)
	}
	// Output:
	// 0
	// 0.5
	// 1
	// 1.5
	// 2
}

func TestInvalidRange(t *testing.T) {
	var err error
	_, _, err = GenerateStrings("\\[0..bar.\\]foo")
	assert.Error(t, err)
	_, _, err = GenerateStrings("\\[5..1..0]")
	assert.Error(t, err)
	_, _, err = GenerateStrings("\\[0..1..x]")
	assert.Error(t, err)
	_, _, err = GenerateStrings("\\[0..1/2]")
	assert.Error(t, err)
}

func BenchmarkConfigFiniteRange(b *testing.B) {
	_BenchmarkConfig("\\[0..10]", b)
}

func BenchmarkConfigInfiniteRange(b *testing.B) {
	_BenchmarkConfig("\\[0..]", b)
}

func BenchmarkConfigStepRange(b *testing.B) {
	_BenchmarkConfig("\\[10..-1..0]", b)
}
