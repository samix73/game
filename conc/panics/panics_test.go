package panics_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gitlab.com/abios/gohelpers/conc/panics"
)

func TestCatcher_Try_NoPanic(t *testing.T) {
	t.Parallel()

	var c panics.Catcher
	executed := false

	c.Try(func() {
		executed = true
	})

	assert.True(t, executed)
	assert.Nil(t, c.Recovered(), "expected no panic to be recovered")
}

func TestCatcher_Try_WithPanic(t *testing.T) {
	t.Parallel()

	var c panics.Catcher
	panicValue := "test panic"

	c.Try(func() {
		panic(panicValue)
	})

	recovered := c.Recovered()
	assert.NotNil(t, recovered, "expected panic to be recovered")
	assert.Equal(t, recovered.Value, panicValue, "expected panic to match")
	assert.NotEmpty(t, recovered.Stack, "expected stack trace to be captured")
}

func TestCatcher_Try_MultiplePanics(t *testing.T) {
	t.Parallel()

	var c panics.Catcher

	// First panic
	c.Try(func() {
		panic("first panic")
	})

	firstRecovered := c.Recovered()
	if firstRecovered == nil || firstRecovered.Value != "first panic" {
		t.Error("first panic not recovered correctly")
	}

	// Second panic should not overwrite the first
	c.Try(func() {
		panic("second panic")
	})

	recovered := c.Recovered()
	assert.Equal(t, recovered.Value, "first panic")
}

func TestCatcher_Repanic(t *testing.T) {
	t.Parallel()

	var c panics.Catcher
	panicValue := "test panic for repanic"

	c.Try(func() {
		panic(panicValue)
	})

	assert.Panics(t, func() {
		c.Repanic()
	})
}
