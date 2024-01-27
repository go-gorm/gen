package field

// Float64 float64 type field
type Float64 struct {
	GenericsInt[float64]
}

// Floor ...
func (field Float64) Floor() Int {
	return Int{NewGenericsInt[int](field.floor())}
}

// Float32 float32 type field
type Float32 struct {
	GenericsInt[float32]
}

// Floor ...
func (field Float32) Floor() Int {
	return Int{NewGenericsInt[int](field.floor())}
}
