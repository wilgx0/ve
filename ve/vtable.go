package ve

import "errors"

type VTable []VRow

func NewVTable() VTable {
	return VTable{}
}

// row将被转置成列
func CreateVTableByCol(row VRow) VTable {
	return nil
}

func CreateVTableByRow(row VRow) VTable {
	return nil
}

func CreateVTable(row, col int) (result VTable) {
	result = make(VTable, 0)
	for i := 0; i < row; i++ {
		result = append(result, make(VRow, col))
	}
	return
}

func (v VTable) FirstRow() (result VRow) {
	return
}

func (v VTable) LastRow() (result VRow) {
	return
}

func (v VTable) FirstCol() (result VRow) {
	return
}

func (v VTable) LastCol() (result VRow) {
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

func (v *VTable) AddColumn(colum VColumn) (i int, err error) {
	return
}

// 水平追加
// vts 里的行数和v里的行数必须相同
// i 返回新添加的列首所在的列数， 列数从0开始
func (v *VTable) HorizontalMerge(vts ...VTable) (i int, err error) {
	if len(vts) == 0 {
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

// 根据列编号创建新的VTable
// 列编号从0开始
func (v VTable) NewVTableByCol(colNums ...int) (result VTable, err error) {
	if len(colNums) == 0 || v.IsEmpty() {
		return
	}
	for _, row := range v {
		var mergeRow VRow
		for _, colNum := range colNums {
			if colNum < 0 || colNum >= v.ColNum() {
				err = errors.New("VTable.NewVTableByCol error : colNums error")
				return
			}
			mergeRow = append(mergeRow, row[colNum])
		}
		result = append(result, mergeRow)
	}
	return
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

func (v VTable) FillColumn(colNum int) VTable {

	return v
}

func (v VTable) FillRow(rowNum int) {
	return
}

func (v VTable) GetColumn(colNum int) {

	return
}

func (v VTable) GetRow(rowNum int) (result VRow) {
	if v.IsEmpty() || rowNum < 0 || rowNum >= v.RowNum() {
		return
	}
	return v[rowNum]
}
