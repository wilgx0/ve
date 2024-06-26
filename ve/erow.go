package ve

// ERow represents a row in the table.
type ERow struct {
	eTable  *ETable
	eCells  []*ECell
	name    string
	payload interface{}
	Sort    int
	Trie    interface{}
	fnName  string
}

func NewERow(et *ETable) *ERow {
	return &ERow{
		eTable: et,
	}
}

func NewERowByArr[T any](arr []T, et *ETable, name ...string) *ERow {
	eRow := NewERow(et)
	for i := 0; i < len(arr); i++ {
		var eCol *ECol
		if et.IsEmpty() {
			eCol = NewECol(et)
		} else {
			eCol = et.GetCol(i)
		}
		if eCol == nil {
			break
		}
		ec := NewECell(arr[i], et, eRow, eCol)
		eRow.AddCell(ec)
	}
	if len(name) > 0 {
		eRow.SetName(name[0])
	}
	return eRow
}

func (row *ERow) SetName(name string) *ERow {
	row.name = name
	return row
}

func (row *ERow) GetName() string {
	return row.name
}

func (row *ERow) SetFnName(name string) *ERow {
	row.fnName = name
	return row
}

func (row *ERow) GetFnName() string {
	return row.fnName
}

func (row *ERow) Cells() []*ECell {
	if row.IsEmpty() {
		return nil
	}
	return row.eCells
}

func (row *ERow) GetECellByName(name string) *ECell {
	for _, cell := range row.eCells {
		if cell.eCol.name == name {
			return cell
		}
	}
	return nil
}

// AddCell adds a cell to the row.
func (row *ERow) AddCell(cell *ECell) {
	row.eCells = append(row.eCells, cell)
}

// GetCell returns a cell at the specified index.
func (row *ERow) GetCell(index int) *ECell {
	if row.IsEmpty() {
		return nil
	}
	if index < 0 || index >= len(row.eCells) {
		return nil
	}

	return row.eCells[index]
}

// Len returns the number of cells in the row.
func (row *ERow) Len() int {
	if row == nil {
		return 0
	}
	return len(row.eCells)
}

// IsEmpty checks if the row is empty.
func (row *ERow) IsEmpty() bool {
	return row == nil || len(row.eCells) == 0
}

// LastCell returns the last cell in the row.
func (row *ERow) LastCell() *ECell {
	return row.GetCell(row.Len() - 1)
}

// FirstCell returns the first cell in the row.
func (row *ERow) FirstCell() *ECell {
	return row.GetCell(0)
}

// IndexOf returns the index of the given cell in the row.
func (row *ERow) IndexOf(cell *ECell) int {
	for i, c := range row.eCells {
		if c == cell {
			return i
		}
	}
	return -1
}

func (row *ERow) SumUint64(fns ...func(*ECell, int) bool) uint64 {
	var sum uint64

	fn := func(*ECell, int) bool {
		return true
	}
	if len(fns) > 0 {
		fn = fns[0]
	}

	for i, cell := range row.eCells {
		if !fn(cell, i) {
			continue
		}
		sum += cell.Uint64()
	}
	return sum
}

func (row *ERow) SumFloat64(fns ...func(*ECell, int) bool) float64 {
	var sum float64

	fn := func(*ECell, int) bool {
		return true
	}
	if len(fns) > 0 {
		fn = fns[0]
	}

	for i, cell := range row.eCells {
		if !fn(cell, i) {
			continue
		}
		sum += cell.Float64()
	}
	return sum
}

func (row *ERow) ToArr(fn func(*ECell, int) interface{}) []interface{} {
	var arr []interface{}
	for index, cell := range row.eCells {
		arr = append(arr, fn(cell, index))
	}
	return arr
}

func (row *ERow) SortCell(fn func(*ECell, *ECell) bool) {
	NewCollection(row.eCells).Sort(func(i, j int) bool {
		return fn(row.eCells[i], row.eCells[j])
	})
}
