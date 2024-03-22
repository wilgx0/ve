package ve

import "time"

// ECell represents a cell in the table.
type ECell struct {
	eTable  *ETable
	eRow    *ERow
	eCol    *ECol
	value   interface{}
	payload interface{}
}

// NewECell creates a new ECell instance.
func NewECell(value interface{}, et *ETable, er *ERow, ec *ECol) *ECell {
	return &ECell{
		eTable: et,
		eRow:   er,
		eCol:   ec,
		value:  value,
	}
}

func (cell *ECell) IsEmpty() bool {
	return cell == nil || cell.value == nil
}

func (cell *ECell) GetERow() *ERow {
	if cell.IsEmpty() {
		return nil
	}

	return cell.eRow
}

func (cell *ECell) GetECol() *ECol {
	if cell.IsEmpty() {
		return nil
	}
	return cell.eCol
}

// GetColNum returns the column number of the cell.
func (cell *ECell) GetColNum() int {
	if cell.IsEmpty() {
		return -1
	}
	if cell.eRow != nil {
		return cell.eRow.IndexOf(cell)
	}
	return -1
}

// GetRowNum returns the row number of the cell.
func (cell *ECell) GetRowNum() int {
	if cell.IsEmpty() {
		return -1
	}

	if cell.eCol != nil {
		return cell.eCol.IndexOf(cell)
	}

	return -1
}

func (cell *ECell) Clone() *ECell {
	return NewECell(cell.Val(), cell.eTable, cell.eRow, cell.eCol)
}

func (cell *ECell) Set(value interface{}) (old interface{}) {
	old = cell.value
	cell.value = value
	return
}

func (cell *ECell) Val() interface{} {
	if cell.IsEmpty() {
		return nil
	}
	return cell.value
}

// Interface is alias of Val.
func (cell *ECell) Interface() interface{} {
	return cell.Val()
}

// Bytes converts and returns `v` as []byte.
func (cell *ECell) Bytes() []byte {
	return Bytes(cell.Val())
}

// String converts and returns `v` as string.
func (cell *ECell) String() string {
	return String(cell.Val())
}

// Bool converts and returns `v` as bool.
func (cell *ECell) Bool() bool {
	return Bool(cell.Val())
}

// Int converts and returns `v` as int.
func (cell *ECell) Int() int {
	return Int(cell.Val())
}

// Int8 converts and returns `v` as int8.
func (cell *ECell) Int8() int8 {
	return Int8(cell.Val())
}

// Int16 converts and returns `v` as int16.
func (cell *ECell) Int16() int16 {
	return Int16(cell.Val())
}

// Int32 converts and returns `v` as int32.
func (cell *ECell) Int32() int32 {
	return Int32(cell.Val())
}

// Int64 converts and returns `v` as int64.
func (cell *ECell) Int64() int64 {
	return Int64(cell.Val())
}

// Uint converts and returns `v` as uint.
func (cell *ECell) Uint() uint {
	return Uint(cell.Val())
}

// Uint8 converts and returns `v` as uint8.
func (cell *ECell) Uint8() uint8 {
	return Uint8(cell.Val())
}

// Uint16 converts and returns `v` as uint16.
func (cell *ECell) Uint16() uint16 {
	return Uint16(cell.Val())
}

// Uint32 converts and returns `v` as uint32.
func (cell *ECell) Uint32() uint32 {
	return Uint32(cell.Val())
}

// Uint64 converts and returns `v` as uint64.
func (cell *ECell) Uint64() uint64 {
	return Uint64(cell.Val())
}

// Float32 converts and returns `v` as float32.
func (cell *ECell) Float32() float32 {
	return Float32(cell.Val())
}

// Float64 converts and returns `v` as float64.
func (cell *ECell) Float64() float64 {
	return Float64(cell.Val())
}

// Time converts and returns `v` as time.Time.
// The parameter `format` specifies the format of the time string using gtime,
// eg: Y-m-d H:i:s.
func (cell *ECell) Time(format ...string) time.Time {
	return Time(cell.Val(), format...)
}

// Duration converts and returns `v` as time.Duration.
// If value of `v` is string, then it uses time.ParseDuration for conversion.
func (cell *ECell) Duration() time.Duration {
	return Duration(cell.Val())
}
