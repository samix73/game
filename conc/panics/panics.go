package panics

import (
	"fmt"
	"runtime/debug"
	"sync/atomic"
)

type Catcher struct {
	recovered atomic.Pointer[Recovered]
}

// Try executes the provided function and captures any panic that occurs.
func (c *Catcher) Try(f func()) {
	defer c.tryRecover()

	f()
}

// Repanic will panic with the value recovered by the last call to Try.
func (c *Catcher) Repanic() {
	if val := c.Recovered(); val != nil {
		panic(val)
	}
}

func (c *Catcher) tryRecover() {
	if val := recover(); val != nil {
		c.recovered.CompareAndSwap(nil, NewRecovered(val))
	}
}

func (c *Catcher) Recovered() *Recovered {
	return c.recovered.Load()
}

type Recovered struct {
	Value any
	Stack []byte
}

func NewRecovered(value any) *Recovered {
	return &Recovered{
		Value: value,
		Stack: debug.Stack(),
	}
}

func (p *Recovered) String() string {
	return fmt.Sprintf("panic: %v\nstacktrace:\n%s\n", p.Value, p.Stack)
}

func (p *Recovered) AsError() error {
	if p == nil {
		return nil
	}

	return &ErrRecovered{*p}
}

func (p *Recovered) RePanic() {
	if p.Value == nil {
		return
	}

	panic(p.Value)
}

var _ error = (*ErrRecovered)(nil)

type ErrRecovered struct{ Recovered }

func (p *ErrRecovered) Error() string {
	if p == nil {
		return "no panic occurred"
	}

	return "recovered from panic: " + p.String()
}

func (p *ErrRecovered) Unwrap() error {
	if err, ok := p.Value.(error); ok {
		return err
	}

	return nil
}
