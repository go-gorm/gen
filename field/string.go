package field

// String string type field
type String struct {
	GenericsString[string]
}

// Bytes []byte type field
type Bytes struct {
	GenericsString[[]byte]
}
