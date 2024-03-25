package ve

// ETable represents the table itself.
type ETable struct {
	eRows []*ERow
}

// NewETable creates a new ETable instance.
func NewETable() *ETable {
	return &ETable{}
}

func (table *ETable) FirstRow() *ERow {
	if table.IsEmpty() {
		return nil
	}
	return table.eRows[0]
}

func (table *ETable) FirstCol() *ECol {
	return table.FirstRow().FirstCell().GetECol()
}

func (table *ETable) IsEmpty() bool {
	return table == nil || len(table.eRows) == 0
}

func (table *ETable) RowNum() int {
	if table.IsEmpty() {
		return 0
	}
	return len(table.eRows)
}

func (table *ETable) ColNum() int {
	if table.IsEmpty() {
		return 0
	}
	return table.FirstRow().Len()

}

func (table *ETable) GetCell(rowNum, colNum int) *ECell {
	if table.IsEmpty() {
		return nil
	}
	if rowNum < 0 || rowNum >= table.RowNum() || colNum < 0 || colNum >= table.ColNum() {
		return nil
	}
	return table.eRows[rowNum].GetCell(colNum)
}

func (table *ETable) LastRow() *ERow {
	if table.IsEmpty() {
		return nil
	}

	return table.eRows[len(table.eRows)-1]
}

func (table *ETable) LastCol() *ECol {
	return table.LastRow().LastCell().GetECol()
}

// AddRow adds a row to the table.
func (table *ETable) AddRow(row *ERow) {
	table.eRows = append(table.eRows, row)
	for _, cell := range row.Cells() {
		cell.eCol.AddCell(cell)
	}

}

// AddCol adds a column to the table.
func (table *ETable) AddCol(col *ECol) {
	if table.IsEmpty() {
		for _, cell := range col.Cells() {
			cell.eRow.AddCell(cell)
			table.eRows = append(table.eRows, cell.eRow)
		}
	} else {
		for _, cell := range col.Cells() {
			cell.eRow.AddCell(cell)
		}
	}
}

func (table *ETable) AddColByFn(fn func(*ERow, int) interface{}, name ...string) {

	if table.IsEmpty() {
		return
	}
	eCol := NewECol(table)
	if len(name) > 0 {
		eCol.SetName(name[0])
	}
	for index, eRow := range table.eRows {
		cell := NewECell(fn(eRow, index), table, eRow, eCol)
		eCol.AddCell(cell)
		eRow.AddCell(cell)
	}
}

func (table *ETable) AddRowByFn(fn func(*ECol, int) interface{}, name ...string) {
	if table.IsEmpty() {
		return
	}
	eRow := NewERow(table)
	if len(name) > 0 {
		eRow.SetName(name[0])
	}
	for i := 0; i < table.ColNum(); i++ {
		eCol := table.GetCol(i)
		cell := NewECell(fn(eCol, i), table, eRow, eCol)
		eCol.AddCell(cell)
		eRow.AddCell(cell)
	}
	table.eRows = append(table.eRows, eRow)
}

func (table *ETable) IndexOf(row *ERow) int {
	for i := 0; i < len(table.eRows); i++ {
		if table.eRows[i] == row {
			return i
		}
	}
	return -1
}

// GetRow returns a row at the specified index.
func (table *ETable) GetRow(index int) *ERow {
	if table.IsEmpty() {
		return nil
	}
	if index < 0 || index >= table.RowNum() {
		return nil
	}
	return table.eRows[index]
}

func (table *ETable) GetCol(index int) *ECol {
	if index < 0 || index >= table.ColNum() {
		return nil
	}

	return table.FirstRow().GetCell(index).GetECol()
}

func (table *ETable) ForCol(fn func(*ECol, int)) {
	if table.IsEmpty() {
		return
	}
	for i := 0; i < table.ColNum(); i++ {
		fn(table.GetCol(i), i)
	}
}

func (table *ETable) ToArr(fn func(*ECell, int, int) interface{}) (result [][]interface{}) {

	for rNum, row := range table.eRows {
		var eRow = make([]interface{}, 0, table.ColNum())
		for cNum, cell := range row.Cells() {
			eRow = append(eRow, fn(cell, rNum, cNum))
		}
		result = append(result, eRow)
	}
	return result
}

func (table *ETable) SortRow(fn func(*ERow, *ERow) bool) {
	NewCollection(table.eRows).Sort(func(i, j int) bool {
		return fn(table.eRows[i], table.eRows[j])
	})
}
