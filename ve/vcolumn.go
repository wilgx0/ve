package ve

type VColumn []*VCell

func NewVColumn(data ...interface{}) VColumn {
	if len(data) == 0 {
		return VColumn{}
	}
	var vc VColumn
	for _, item := range data {
		vc = append(vc, NewVCell(item))
	}
	return vc
}

func (v VColumn) Each(fn func(int, *VCell)) {
	for index, item := range v {
		fn(index, item)
	}
}

func (v VColumn) RowNum() int {
	return len(v)
}

func (v VColumn) IsEmpty() bool {
	return len(v) == 0
}

func (v VColumn) Clone() (result VColumn) {
	for _, item := range v {
		result = append(result, NewVCell(item.Val()))
	}
	return
}

func (v VColumn) ValidRowNum(rowNum int) bool {
	return rowNum >= 0 && rowNum < v.RowNum()
}

func (v *VColumn) DelRow(rowNum int) VColumn {
	if v.IsEmpty() || !v.ValidRowNum(rowNum) {
		return *v
	}
	copy((*v)[rowNum:], (*v)[rowNum+1:])
	(*v)[len(*v)-1] = nil
	*v = (*v)[:len(*v)-1]
	return *v
}
