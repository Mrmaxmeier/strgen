package strgen

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicText(t *testing.T) {
	channel, amount, err := GenerateStrings("foo bar")
	assert.Nil(t, err)
	assert.Equal(t, int64(1), amount, "should produce one result")
	assert.Equal(t, "foo bar", <-channel, "should produce the correct result")
}

func TestDone(t *testing.T) {
	test := func(after int) {
		g := &Generator{Source: "\\[0..]"}
		assert.Nil(t, g.Configure())
		assert.Equal(t, int64(-1), g.Amount)
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			g.Generate()
			wg.Done()
		}()
		for c := 0; c < after; c++ {
			_, err := g.Next()
			assert.Nil(t, err)
		}
		g.Close()
		_, err := g.Next()
		assert.NotNil(t, err)
		wg.Wait()
	}
	for amount := 1; amount < 20; amount++ {
		test(amount)
	}
}
