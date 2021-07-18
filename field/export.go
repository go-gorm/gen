package field

import "gorm.io/gorm/clause"

type FieldOption func(clause.Column) clause.Column

// TODO implement validator options

// ======================== generic field =======================

func NewField(table, column string, opts ...FieldOption) Field {
	return Field{expr: expr{Col: toColumn(table, column, opts...)}}
}

// ======================== integer =======================

func NewInt(table, column string, opts ...FieldOption) Int {
	return Int{expr: expr{Col: toColumn(table, column, opts...)}}
}

func NewInt8(table, column string, opts ...FieldOption) Int8 {
	return Int8{expr: expr{Col: toColumn(table, column, opts...)}}
}

func NewInt16(table, column string, opts ...FieldOption) Int16 {
	return Int16{expr: expr{Col: toColumn(table, column, opts...)}}
}

func NewInt32(table, column string, opts ...FieldOption) Int32 {
	return Int32{expr: expr{Col: toColumn(table, column, opts...)}}
}

func NewInt64(table, column string, opts ...FieldOption) Int64 {
	return Int64{expr: expr{Col: toColumn(table, column, opts...)}}
}

func NewUint(table, column string, opts ...FieldOption) Uint {
	return Uint{expr: expr{Col: toColumn(table, column, opts...)}}
}

func NewUint8(table, column string, opts ...FieldOption) Uint8 {
	return Uint8{expr: expr{Col: toColumn(table, column, opts...)}}
}

func NewUint16(table, column string, opts ...FieldOption) Uint16 {
	return Uint16{expr: expr{Col: toColumn(table, column, opts...)}}
}

func NewUint32(table, column string, opts ...FieldOption) Uint32 {
	return Uint32{expr: expr{Col: toColumn(table, column, opts...)}}
}

func NewUint64(table, column string, opts ...FieldOption) Uint64 {
	return Uint64{expr: expr{Col: toColumn(table, column, opts...)}}
}

// ======================== float =======================

func NewFloat32(table, column string, opts ...FieldOption) Float32 {
	return Float32{expr: expr{Col: toColumn(table, column, opts...)}}
}

func NewFloat64(table, column string, opts ...FieldOption) Float64 {
	return Float64{expr: expr{Col: toColumn(table, column, opts...)}}
}

// ======================== string =======================

func NewString(table, column string, opts ...FieldOption) String {
	return String{expr: expr{Col: toColumn(table, column, opts...)}}
}

func NewBytes(table, column string, opts ...FieldOption) Bytes {
	return Bytes{expr: expr{Col: toColumn(table, column, opts...)}}
}

// ======================== bool =======================

func NewBool(table, column string, opts ...FieldOption) Bool {
	return Bool{expr: expr{Col: toColumn(table, column, opts...)}}
}

// ======================== time =======================

func NewTime(table, column string, opts ...FieldOption) Time {
	return Time{expr: expr{Col: toColumn(table, column, opts...)}}
}

func toColumn(table, column string, opts ...FieldOption) clause.Column {
	col := clause.Column{Table: table, Name: column}
	for _, opt := range opts {
		col = opt(col)
	}
	return col
}
