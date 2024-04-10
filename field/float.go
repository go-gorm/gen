package field

// Float64 float64 type field
type Float64 struct {
	Number[float64]
}

// Floor implement floor method
func (field Float64) Floor() Number[int] { return newNumber[int](field.floor()) }

// Float32 float32 type field
type Float32 struct {
	Number[float32]
}

// Floor implement floor method
func (field Float32) Floor() Number[int] { return newNumber[int](field.floor()) }
