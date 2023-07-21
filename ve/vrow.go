package ve

type VRow []*VCell

func NewVRow(data ...interface{}) VRow {
	if len(data) == 0 {
		return VRow{}
	}
	var vrow VRow
	for _, item := range data {
		vrow = append(vrow, NewVCell(item))
	}
	return vrow
}

func (v VRow) ColNum() int {
	return len(v)
}

func (v VRow) IsEmpty() bool {
	return len(v) == 0
}

func (v VRow) Clone() (result VRow) {
	for _, item := range v {
		result = append(result, NewVCell(item.Val()))
	}
	return
}
