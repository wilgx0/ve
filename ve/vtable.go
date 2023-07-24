package ve

import "errors"

type VTable []VRow

func NewVTable() VTable {
	return VTable{}
}

func CreateVTableByCol(col VColumn) (vt VTable, err error) {
	vt = NewVTable()
	_, err = vt.AddColumn(col)
	return
}

func CreateVTableByRow(row VRow) (vt VTable, err error) {
	vt = NewVTable()
	_, err = vt.AddRow(row)
	return
}

func CreateVTable(row, col int) (result VTable) {
	result = make(VTable, 0)
	for i := 0; i < row; i++ {
		result = append(result, make(VRow, col))
	}
	return
}

// 获取并移除第一行
func (v VTable) FirstRow() (result VRow) {
	if v.IsEmpty() {
		return
	}
	result = v.GetRow(0)
	v.DelRow(0)
	return
}

// 获取并移除最后一行
func (v VTable) LastRow() (result VRow) {
	if v.IsEmpty() {
		return
	}
	lastIndex := v.RowNum() - 1
	result = v.GetRow(lastIndex)
	v.DelRow(lastIndex)
	return
}

// 获取并移除第一列
func (v VTable) FirstCol() (result VColumn) {
	if v.IsEmpty() {
		return
	}
	result = v.GetColumn(0)
	v.DelColumn(0)
	return
}

// 获取并移除最后一列
func (v VTable) LastCol() (result VColumn) {
	if v.IsEmpty() {
		return
	}
	lastIndex := v.ColNum() - 1
	result = v.GetColumn(lastIndex)
	v.DelRow(lastIndex)
	return
}

// 装换成 [][]interface{}
func (v VTable) ToInterface() (result [][]interface{}) {
	for _, item := range v {
		var row []interface{}
		for _, el := range item {
			row = append(row, el.Val())
		}
		result = append(result, row)
	}
	return
}

func (v VTable) ColNum() int {
	if len(v) > 0 {
		return len(v[0])
	} else {
		return 0
	}
}

func (v VTable) RowNum() int {
	return len(v)
}

func (v VTable) IsEmpty() bool {
	return len(v) == 0
}

// 追加行
// rows不传则新增一行空行
// rows里的列数和v里的列数必须相同
// i 返回新添加的行首所在的行数， 行数从0开始
func (v *VTable) AddRow(rows ...VRow) (i int, err error) {
	if len(rows) == 0 {
		if len(*v) == 0 {
			*v = CreateVTable(1, 1)
			return 0, nil
		}
		row := make(VRow, len((*v)[0]))
		*v = append(*v, row)
		return len(*v) - 1, nil
	} else {
		if len(*v) == 0 {
			for _, row := range rows {
				*v = append(*v, row)
			}
			return 0, nil
		}
		colNum := v.ColNum()
		for _, row := range rows {
			if colNum != row.ColNum() {
				return -1, errors.New("VTable.AddRow error : column count error")
			}
		}
		i = len(*v)
		for _, row := range rows {
			*v = append(*v, row)
		}
		return i, nil
	}
}

// 追加列
// columns不传则新增一列空列
// columns里列的行数和v里的行数必须相同
// i 返回新添加的列首所在的列数， 列数从0开始
func (v *VTable) AddColumn(columns ...VColumn) (i int, err error) {
	if len(columns) == 0 {
		if len(*v) == 0 {
			*v = CreateVTable(1, 1)
			return 0, nil
		}
		rowNum := v.RowNum()
		for j := 0; j < rowNum; j++ {
			(*v)[i] = append((*v)[i], nil)
		}
		return v.ColNum() - 1, nil
	} else {
		if len(*v) == 0 {
			rowNum := columns[0].RowNum()
			for _, col := range columns {
				if col.RowNum() != rowNum {
					return -1, errors.New("VTable.AddColumn error : row count error")
				}
			}
			for j := 0; j < rowNum; j++ {
				var mergeRows VRow
				for _, col := range columns {
					mergeRows = append(mergeRows, col[j])
				}
				*v = append(*v, mergeRows)
			}
			return 0, nil
		}

		rowNum := v.RowNum()
		for _, col := range columns {
			if col.RowNum() != rowNum {
				return -1, errors.New("VTable.AddColumn error : row count error")
			}
		}
		i = v.ColNum()
		for j := 0; j < rowNum; j++ {
			var mergeRow VRow
			for _, col := range columns {
				mergeRow = append(mergeRow, col[j])
			}
			(*v)[j] = append((*v)[j], mergeRow...)
		}
		return
	}
}

// 垂直合并
// vts 里的行数和v里的行数必须相同
// i 返回新添加的列首所在的列数， 列数从0开始
func (v *VTable) VerticalMerge(vts ...VTable) (i int, err error) {
	var j *int
	for _, vt := range vts {
		if j == nil {
			*j, err = v.AddRow(vt...)
		} else {
			_, err = v.AddRow(vt...)
		}

		if err != nil {
			return
		}
	}
	if j != nil {
		i = *j
	}
	return
}

// 横向合并
// vts 里的行数和v里的行数必须相同
// i 返回新添加的列首所在的列数， 列数从0开始
func (v *VTable) HorizontalMerge(vts ...VTable) (i int, err error) {
	if len(vts) == 0 {
		return 0, nil
	} else {
		if len(*v) == 0 {
			rowNum := vts[0].RowNum()
			for _, table := range vts {
				if table.RowNum() != rowNum {
					return -1, errors.New("VTable.AddColumn error : row count error")
				}
			}
			for j := 0; j < rowNum; j++ {
				var mergeRows VRow
				for _, table := range vts {
					mergeRows = append(mergeRows, table[j]...)
				}
				*v = append(*v, mergeRows)
			}
			return 0, nil
		}
		rowNum := v.RowNum()
		for _, table := range vts {
			if table.RowNum() != rowNum {
				return -1, errors.New("VTable.AddColumn error : row count error")
			}
		}
		i = v.ColNum()
		for j := 0; j < rowNum; j++ {
			var mergeRow VRow
			for _, table := range vts {
				mergeRow = append(mergeRow, table[j]...)
			}
			(*v)[j] = append((*v)[j], mergeRow...)
		}
		return
	}
}

func (v VTable) FilterColumn(fn func(index int, col VColumn) bool) (result VTable) {
	colNum := v.ColNum()
	for i := 0; i < colNum; i++ {
		column := v.GetColumn(i)
		if fn(i, column) {
			result.AddColumn(column)
		}
	}
	return
}

func (v VTable) FilterRow(fn func(index int, row VRow) bool) (result VTable) {
	for index, item := range v {
		if fn(index, item) {
			result.AddRow(item)
		}
	}
	return
}

func (v VTable) EachColumn(fn func(index int, col VColumn)) {
	colNum := v.ColNum()
	for i := 0; i < colNum; i++ {
		column := v.GetColumn(i)
		fn(i, column)
	}
}

func (v VTable) EachRow(fn func(index int, row VRow)) {
	for index, item := range v {
		fn(index, item)
	}
}

func (v *VTable) Clear() {
	*v = (*v)[:0]
}

func (v VTable) Clone() (result VTable) {
	for _, row := range v {
		result = append(result, row.Clone())
	}
	return
}

func (v VTable) FillColumn(colNum int, col VColumn) (err error) {
	if v.IsEmpty() || !v.ValidColNum(colNum) {
		return
	}

	if v.RowNum() != col.RowNum() {
		err = errors.New("VTable.FillColumn error : 'col' param error")
		return
	}

	for index, row := range v {
		row[colNum] = col[index]
	}
	return
}

func (v VTable) FillRow(rowNum int, row VRow) (err error) {
	if v.IsEmpty() || !v.ValidRowNum(rowNum) {
		return
	}
	if v.ColNum() != row.ColNum() {
		err = errors.New("VTable.FillRow error : 'row' param error")
		return
	}
	v[rowNum] = row
	return
}

func (v VTable) ValidColNum(colNum int) bool {
	return colNum >= 0 && colNum < v.ColNum()
}

func (v VTable) ValidRowNum(rowNum int) bool {
	return rowNum < 0 && rowNum < v.RowNum()
}

// colNum 从0开始
func (v VTable) GetColumn(colNum int) (result VColumn) {
	if v.IsEmpty() || !v.ValidColNum(colNum) {
		return
	}
	for i := 0; i < len(v); i++ {
		result = append(result, v[i][colNum])
	}
	return
}

// rowNum 从0开始
func (v VTable) GetRow(rowNum int) (result VRow) {
	if v.IsEmpty() || !v.ValidRowNum(rowNum) {
		return
	}
	return v[rowNum]
}

// rowNum及colNum从0开始
func (v VTable) GetCell(rowNum int, colNum int) (result *VCell) {
	if v.IsEmpty() {
		return
	}
	if v.ValidRowNum(rowNum) && v.ValidColNum(colNum) {
		return v[rowNum][colNum]
	}
	return
}

func (v VTable) DelColumn(colNum int) {
	if v.IsEmpty() || !v.ValidColNum(colNum) {
		return
	}
	for i := 0; i < v.RowNum(); i++ {
		v[i] = v[i].DelCol(colNum)
	}
}

func (v *VTable) DelRow(rowNum int) {
	if v.IsEmpty() || !v.ValidRowNum(rowNum) {
		return
	}
	copy((*v)[rowNum:], (*v)[rowNum+1:])
	(*v)[len(*v)-1] = nil
	*v = (*v)[:len(*v)-1]
}
