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
	g := &Generator{source: "\\[0..]"}
	assert.Nil(t, g.configure())
	assert.Equal(t, int64(-1), g.amount)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		g.generate()
		wg.Done()
	}()
	var counter int
	for _ = range g.results {
		counter++
		if counter == 42 {
			g.kill()
		}
	}
	assert.Equal(t, 42, counter, "should terminate after 42 results")
	wg.Wait()
}
