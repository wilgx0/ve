package ve

// ECol represents a column in the table.
type ECol struct {
	ERow
}

func NewECol(et *ETable) *ECol {
	return &ECol{
		ERow: ERow{
			eTable: et,
		},
	}
}

func NewEColByArr[T any](arr []T, et *ETable, name ...string) *ECol {
	eCol := NewECol(et)
	for i := 0; i < len(arr); i++ {
		var row *ERow
		if et.IsEmpty() {
			row = NewERow(et)
		} else {
			row = et.GetRow(i)
		}
		if row == nil {
			break
		}
		ec := NewECell(arr[i], et, row, eCol)
		eCol.AddCell(ec)
	}
	if len(name) > 0 {
		eCol.SetName(name[0])
	}
	return eCol
}
