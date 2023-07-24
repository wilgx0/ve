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

func (v VRow) Each(fn func(int, *VCell)) {
	for index, item := range v {
		fn(index, item)
	}
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

func (v VRow) ValidColNum(colNum int) bool {
	return colNum >= 0 && colNum < v.ColNum()
}

func (v *VRow) DelCol(colNum int) VRow {
	if v.IsEmpty() || !v.ValidColNum(colNum) {
		return *v
	}
	copy((*v)[colNum:], (*v)[colNum+1:])
	(*v)[len(*v)-1] = nil
	*v = (*v)[:len(*v)-1]
	return *v
}
