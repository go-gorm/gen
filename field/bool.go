package field

type Bool Field

func (field Bool) Not() Bool {
	return Bool{field.not()}
}

func (field Bool) Is(value bool) Expr {
	return field.is(value)
}

func (field Bool) And(value bool) Expr {
	return Bool{field.and(value)}
}

func (field Bool) Or(value bool) Expr {
	return Bool{field.or(value)}
}

func (field Bool) Xor(value bool) Expr {
	return Bool{field.xor(value)}
}

func (field Bool) BitXor(value bool) Expr {
	return Bool{field.bitXor(value)}
}

func (field Bool) BitAnd(value bool) Expr {
	return Bool{field.bitAnd(value)}
}

func (field Bool) BitOr(value bool) Expr {
	return Bool{field.bitOr(value)}
}

func (field Bool) Value(value bool) AssignExpr {
	return field.value(value)
}

func (field Bool) Zero() AssignExpr {
	return field.value(false)
}
