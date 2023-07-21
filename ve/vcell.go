package ve

import (
	"encoding/json"
	"time"
)

type VCell struct {
	value interface{}
}

func NewVCell(value interface{}) *VCell {
	v := VCell{}
	v.value = value
	return &v
}

func (v *VCell) Clone() *VCell {
	return NewVCell(v.Val())
}

func (v *VCell) Set(value interface{}) (old interface{}) {
	old = v.value
	v.value = value
	return
}

func (v *VCell) Val() interface{} {
	if v == nil {
		return nil
	}
	return v.value
}

// Interface is alias of Val.
func (v *VCell) Interface() interface{} {
	return v.Val()
}

// Bytes converts and returns `v` as []byte.
func (v *VCell) Bytes() []byte {
	return Bytes(v.Val())
}

// String converts and returns `v` as string.
func (v *VCell) String() string {
	return String(v.Val())
}

// Bool converts and returns `v` as bool.
func (v *VCell) Bool() bool {
	return Bool(v.Val())
}

// Int converts and returns `v` as int.
func (v *VCell) Int() int {
	return Int(v.Val())
}

// Int8 converts and returns `v` as int8.
func (v *VCell) Int8() int8 {
	return Int8(v.Val())
}

// Int16 converts and returns `v` as int16.
func (v *VCell) Int16() int16 {
	return Int16(v.Val())
}

// Int32 converts and returns `v` as int32.
func (v *VCell) Int32() int32 {
	return Int32(v.Val())
}

// Int64 converts and returns `v` as int64.
func (v *VCell) Int64() int64 {
	return Int64(v.Val())
}

// Uint converts and returns `v` as uint.
func (v *VCell) Uint() uint {
	return Uint(v.Val())
}

// Uint8 converts and returns `v` as uint8.
func (v *VCell) Uint8() uint8 {
	return Uint8(v.Val())
}

// Uint16 converts and returns `v` as uint16.
func (v *VCell) Uint16() uint16 {
	return Uint16(v.Val())
}

// Uint32 converts and returns `v` as uint32.
func (v *VCell) Uint32() uint32 {
	return Uint32(v.Val())
}

// Uint64 converts and returns `v` as uint64.
func (v *VCell) Uint64() uint64 {
	return Uint64(v.Val())
}

// Float32 converts and returns `v` as float32.
func (v *VCell) Float32() float32 {
	return Float32(v.Val())
}

// Float64 converts and returns `v` as float64.
func (v *VCell) Float64() float64 {
	return Float64(v.Val())
}

// Time converts and returns `v` as time.Time.
// The parameter `format` specifies the format of the time string using gtime,
// eg: Y-m-d H:i:s.
func (v *VCell) Time(format ...string) time.Time {
	return Time(v.Val(), format...)
}

// Duration converts and returns `v` as time.Duration.
// If value of `v` is string, then it uses time.ParseDuration for conversion.
func (v *VCell) Duration() time.Duration {
	return Duration(v.Val())
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (v *VCell) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.Val())
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (v *VCell) UnmarshalJSON(b []byte) error {
	var i interface{}
	err := UnmarshalUseNumber(b, &i)
	if err != nil {
		return err
	}
	v.Set(i)
	return nil
}

// UnmarshalValue is an interface implement which sets any type of value for Var.
func (v *VCell) UnmarshalValue(value interface{}) error {
	v.Set(value)
	return nil
}
