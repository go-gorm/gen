package field

import "gorm.io/gorm/clause"

// ======================== generic field =======================

func NewField(column string) Field {
	return Field{expr: expr{Col: clause.Column{Name: column}}}
}

// ======================== integer =======================

func NewInt(column string) Int {
	return Int{expr: expr{Col: clause.Column{Name: column}}}
}

func NewInt8(column string) Int8 {
	return Int8{expr: expr{Col: clause.Column{Name: column}}}
}

func NewInt16(column string) Int16 {
	return Int16{expr: expr{Col: clause.Column{Name: column}}}
}

func NewInt32(column string) Int32 {
	return Int32{expr: expr{Col: clause.Column{Name: column}}}
}

func NewInt64(column string) Int64 {
	return Int64{expr: expr{Col: clause.Column{Name: column}}}
}

func NewUint(column string) Uint {
	return Uint{expr: expr{Col: clause.Column{Name: column}}}
}

func NewUint8(column string) Uint8 {
	return Uint8{expr: expr{Col: clause.Column{Name: column}}}
}

func NewUint16(column string) Uint16 {
	return Uint16{expr: expr{Col: clause.Column{Name: column}}}
}

func NewUint32(column string) Uint32 {
	return Uint32{expr: expr{Col: clause.Column{Name: column}}}
}

func NewUint64(column string) Uint64 {
	return Uint64{expr: expr{Col: clause.Column{Name: column}}}
}

// ======================== float =======================

func NewFloat32(column string) Float32 {
	return Float32{expr: expr{Col: clause.Column{Name: column}}}
}

func NewFloat64(column string) Float64 {
	return Float64{expr: expr{Col: clause.Column{Name: column}}}
}

// ======================== string =======================

func NewString(column string) String {
	return String{expr: expr{Col: clause.Column{Name: column}}}
}

func NewBytes(column string) Bytes {
	return Bytes{expr: expr{Col: clause.Column{Name: column}}}
}

// ======================== bool =======================

func NewBool(column string) Bool {
	return Bool{expr: expr{Col: clause.Column{Name: column}}}
}

// ======================== time =======================

func NewTime(column string) Time {
	return Time{expr: expr{Col: clause.Column{Name: column}}}
}
