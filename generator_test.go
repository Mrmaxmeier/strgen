package strgen

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerator(t *testing.T) {
	g := Generator{Source: "teststr \\(a|b)"}
	assert.NoError(t, g.Configure())
	assert.Equal(t, int64(2), g.Amount, "should produce two results")
	assert.NotNil(t, g.DoneCh)
	assert.NotNil(t, g.Results)
	assert.NoError(t, g.Err)
	assert.False(t, g.Infinite)
	assert.NotEmpty(t, g.Iterators)
}

func TestRangeAmountAlive(t *testing.T) {
	g := Generator{Source: "blah \\(a|b) \\[1..5] \\(foo|bar)"}
	g.Configure()
	go g.Generate()
	assert.Equal(t, int64(20), g.Amount)
	for i := 0; i < 20; i++ {
		g.Next()
		assert.True(t, g.Alive(), fmt.Sprintf("cycle %d", i))
	}
	g.Next()
	assert.False(t, g.Alive())
}

func ExampleGenerator() {
	g := Generator{Source: "teststr"}
	g.Configure()
	go g.Generate()
	str, err := g.Next()
	fmt.Printf("string: '%s' err: %v\n", str, err)
	str, err = g.Next()
	fmt.Printf("err: %v\n", err)
	// Output:
	// string: 'teststr' err: <nil>
	// err: channel closed
}

func _BenchmarkConfig(source string, b *testing.B) {
	g := Generator{Source: source}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		g.Configure()
	}
}

func BenchmarkConfigEmpty(b *testing.B) {
	_BenchmarkConfig("", b)
}

func BenchmarkConfigSimpleText(b *testing.B) {
	_BenchmarkConfig("foobar", b)
}

func BenchmarkReadValues(b *testing.B) {
	g := Generator{Source: "\\[0..]"}
	g.Configure()
	go g.Generate()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		g.Next()
	}
}
