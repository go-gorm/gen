package field

// Bool boolean type field
type Bool Field

// Not ...
func (field Bool) Not() Bool {
	return Bool{field.not()}
}

// Is ...
func (field Bool) Is(value bool) Expr {
	return field.is(value)
}

// And boolean and
func (field Bool) And(value bool) Expr {
	return Bool{field.and(value)}
}

// Or boolean or
func (field Bool) Or(value bool) Expr {
	return Bool{field.or(value)}
}

// Xor ...
func (field Bool) Xor(value bool) Expr {
	return Bool{field.xor(value)}
}

// BitXor ...
func (field Bool) BitXor(value bool) Expr {
	return Bool{field.bitXor(value)}
}

// BitAnd ...
func (field Bool) BitAnd(value bool) Expr {
	return Bool{field.bitAnd(value)}
}

// BitOr ...
func (field Bool) BitOr(value bool) Expr {
	return Bool{field.bitOr(value)}
}

// Value ...
func (field Bool) Value(value bool) AssignExpr {
	return field.value(value)
}

// Zero ...
func (field Bool) Zero() AssignExpr {
	return field.value(false)
}
