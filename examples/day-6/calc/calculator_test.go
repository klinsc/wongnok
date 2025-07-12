package calc_test

import (
	"testing"

	"github.com/klins/devpool/go-day6/wongnok/examples/day-6/calc"
	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	result := calc.Add(2, 3)
	assert.Equal(t, 5, result, "2 + 3 should equal 5")

	result = calc.Add(-1, 1)
	assert.Equal(t, 0, result, "-1 + 1 should equal 0")

	result = calc.Add(0, 0)
	assert.Equal(t, 0, result, "0 + 0 should equal 0")
}
