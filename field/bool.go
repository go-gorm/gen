package field

// Bool boolean type field
type Bool struct {
	genericsField[bool]
}

// Not ...
func (field Bool) Not() Bool {
	return Bool{genericsField[bool]{field.not()}}
}

// Is ...
func (field Bool) Is(value bool) Expr {
	return field.is(value)
}

// And boolean and
func (field Bool) And(value bool) Expr {
	return Bool{genericsField[bool]{field.and(value)}}
}

// Or boolean or
func (field Bool) Or(value bool) Expr {
	return Bool{genericsField[bool]{field.or(value)}}
}

// Xor ...
func (field Bool) Xor(value bool) Expr {
	return Bool{genericsField[bool]{field.xor(value)}}
}

// BitXor ...
func (field Bool) BitXor(value bool) Expr {
	return Bool{genericsField[bool]{field.bitXor(value)}}
}

// BitAnd ...
func (field Bool) BitAnd(value bool) Expr {
	return Bool{genericsField[bool]{field.bitAnd(value)}}
}

// BitOr ...
func (field Bool) BitOr(value bool) Expr {
	return Bool{genericsField[bool]{field.bitOr(value)}}
}

// Zero ...
func (field Bool) Zero() AssignExpr {
	return field.value(false)
}
