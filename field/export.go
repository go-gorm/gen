package field

import "gorm.io/gorm/clause"

// ======================== 通用类型 =======================

func NewField(column string) Field {
	return Field{expr: expr{Col: clause.Column{Name: column}}}
}

// ======================== 整数类型 =======================

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

// ======================== 浮点数类型 =======================

func NewFloat32(column string) Float32 {
	return Float32{expr: expr{Col: clause.Column{Name: column}}}
}

func NewFloat64(column string) Float64 {
	return Float64{expr: expr{Col: clause.Column{Name: column}}}
}

// ======================== 字符串类型 =======================

func NewString(column string) String {
	return String{expr: expr{Col: clause.Column{Name: column}}}
}

func NewBytes(column string) Bytes {
	return Bytes{expr: expr{Col: clause.Column{Name: column}}}
}

// ======================== 布尔类型 =======================

func NewBool(column string) Bool {
	return Bool{expr: expr{Col: clause.Column{Name: column}}}
}

// ======================== 时间类型 =======================

func NewTime(column string) Time {
	return Time{expr: expr{Col: clause.Column{Name: column}}}
}
