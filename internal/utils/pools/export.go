// Package pools provides bounded goroutine pools used during code generation.
package pools

// NewPool return a new pool
func NewPool(size int) Pool {
	var p pool
	p.Init(size)
	return &p
}
