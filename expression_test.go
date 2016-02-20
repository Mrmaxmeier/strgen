package strgen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExprBasics(t *testing.T) {
	values, amount, _ := GenerateStrings("\\{3$x^2}")
	assert.Equal(t, int64(3), amount)
	assert.Equal(t, "0", <-values)
	assert.Equal(t, "1", <-values)
	assert.Equal(t, "4", <-values)
}

func TestExprInf(t *testing.T) {
	values, amount, err := GenerateStrings("\\{-1$x}")
	assert.Equal(t, int64(-1), amount)
	assert.NoError(t, err)
	assert.Equal(t, "0", <-values)
	assert.Equal(t, "1", <-values)
}

func TestExprRespectsCycle(t *testing.T) {
	values, amount, _ := GenerateStrings("\\{2$x} \\[0..5]")
	assert.Equal(t, int64(2*6), amount)
	assert.Equal(t, "0 0", <-values)
	assert.Equal(t, "1 0", <-values)
	assert.Equal(t, "0 1", <-values)
	assert.Equal(t, "1 1", <-values)
}

func TestExprFloats(t *testing.T) {
	values, amount, _ := GenerateStrings("\\{3$1/(x+1)}")
	assert.Equal(t, int64(3), amount)
	assert.Equal(t, "1", <-values)
	assert.Equal(t, "0.5", <-values)
	assert.Equal(t, "0.3333333333333333", <-values)
}

func TestExprInvalid(t *testing.T) {
	var err error
	_, _, err = GenerateStrings("\\{a$}")
	assert.Error(t, err)
	_, _, err = GenerateStrings("\\{0.3$x}")
	assert.Error(t, err)
	_, _, err = GenerateStrings("\\{3$x")
	assert.Error(t, err)
}

func BenchmarkConfigExpr(b *testing.B) {
	_BenchmarkConfig("\\{1$x}", b)
}

func BenchmarkReadExpr(b *testing.B) {
	g := Generator{Source: "\\{-1$x}"}
	g.Configure()
	go g.Generate()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		g.Next()
	}
}
