package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRandomWordFromList(t *testing.T) {
	word := RandomWordFromList()
	assert.Contains(t, "computer house airplane helicopter table elephant giraffe software hardware france", word)
}

func TestCreateInitGuess(t *testing.T) {
	guess := CreateInitGuess("europe")
	assert.Equal(t, "______", guess)
}

func TestPlayGame(t *testing.T) {
	g := Game{12, "Joe", "desk", "d___"}

	g, ok := PlayGame(g, "s")
	assert.Equal(t, "d_s_", g.Guess)
	assert.Equal(t, ok, true)
}

func TestPlayGameMultipleMatches(t *testing.T) {
	g := Game{12, "Joe", "banana", "b_____"}

	g, ok := PlayGame(g, "a")
	assert.Equal(t, "ba_a_a", g.Guess)
	assert.Equal(t, ok, true)
}

func TestPlayGameWithnotExistentChar(t *testing.T) {
	g := Game{12, "Joe", "desk", "d___"}

	g, ok := PlayGame(g, "a")
	assert.Equal(t, "d___", g.Guess)
	assert.Equal(t, ok, false)
}
