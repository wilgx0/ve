package ve

type Grid [][]interface{}

func NewGrid(data [][]interface{}) Grid {
	return data
}

func (g Grid) IsEmpty() bool {
	return len(g) == 0
}

func (g Grid) RowNum() int {
	return len(g)
}

func (g Grid) ColNum() int {
	if len(g) == 0 {
		return 0
	}
	return len(g[0])
}

func (g Grid) MergeLeft(other Grid) (result Grid) {
	if other.IsEmpty() {
		return g
	}
	for i, item := range other {
		item = append(item, g[i]...)
		result = append(result, item)
	}
	return
}

func (g Grid) MergeTop(other Grid) (result Grid) {
	for _, item := range other {
		result = append(result, item)
	}
	result = append(result, g...)
	return
}
