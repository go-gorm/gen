package helper

// Interface Custom query interface define.
type Interface interface {
	// Name return interface name
	Name() string
	// Methods return interface methods
	Methods() []Method
	// Package return interface packet
	Package() string
}
