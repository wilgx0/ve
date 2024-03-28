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

func NewEColByTrie[K comparable, V any](trie *Trie[K, V], et *ETable, name ...string) *ECol {
	eCol := NewECol(et)
	trieArr := trie.Bottom()
	for i := 0; i < len(trieArr); i++ {
		var row *ERow
		if et.IsEmpty() {
			row = NewERow(et)
		} else {
			row = et.GetRow(i)
		}
		if row == nil {
			break
		}
		ec := NewECell(trieArr[i].Key, et, row, eCol)
		ec.Trie = trieArr[i]
		row.Trie = trieArr[i]
		row.SetName(trieArr[i].GetKey())
		eCol.AddCell(ec)
	}
	if len(name) > 0 {
		eCol.SetName(name[0])
	}
	return eCol
}
