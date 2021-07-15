package field

type Bool Field

func (field Bool) Not() Expr {
	return field.not()
}

func (field Bool) Is(value bool) Expr {
	return field.is(value)
}

func (field Bool) And(value bool) Expr {
	return field.and(value)
}

func (field Bool) Or(value bool) Expr {
	return field.or(value)
}

func (field Bool) Xor(value bool) Expr {
	return field.xor(value)
}

func (field Bool) BitXor(value bool) Expr {
	return field.bitXor(value)
}

func (field Bool) BitAnd(value bool) Expr {
	return field.bitAnd(value)
}

func (field Bool) BitOr(value bool) Expr {
	return field.bitOr(value)
}
