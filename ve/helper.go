package ve

import (
	"encoding/json"
	"math"
	"strings"
)

func AddRow[T any](table *ETable, arr []T, name ...string) *ERow {
	row := NewERowByArr(arr, table, name...)
	table.AddRow(row)
	return row
}

func AddRowByFn(table *ETable, fn func(*ECol, int) interface{}, name ...string) *ERow {
	return table.AddRowByFn(fn, name...)
}

func AddCol[T any](table *ETable, arr []T, name ...string) *ECol {
	col := NewEColByArr(arr, table, name...)
	table.AddCol(col)
	return col
}

func AddColByFn(table *ETable, fn func(*ERow, int) interface{}, name ...string) *ECol {
	return table.AddColByFn(fn, name...)
}

// trie会被添加到每一行
func AddColByTrie[K comparable, V any](table *ETable, trie *Trie[K, V], name ...string) *ECol {
	col := NewEColByTrie(trie, table, name...)
	table.AddCol(col)
	return col
}

type CalculateByTrieOpts[V any] struct {
	GetVal    func(Collection[V]) interface{}
	GetFnName func() string
}

// 数据区的统计
func CalculateByTrie3[K comparable, V any](table *ETable, c Collection[V], fns ...CalculateByTrieOpts[V]) {
	fnCount := len(fns)
	if !table.IsEmpty() || fnCount == 0 {
		return
	}
	eRow := NewERow(table)
	for _, opt := range fns {
		fn, fn2 := opt.GetVal, opt.GetFnName
		var eCol *ECol
		eCol = NewECol(table)
		cell := NewECell(fn(c), table, eRow, eCol)
		cell.payload = c
		eCol.AddCell(cell)
		eCol.SetName(fn2())
		eCol.SetFnName(fn2())
		eRow.AddCell(cell)
	}
	table.eRows = append(table.eRows, eRow)
}

// 列及数据区域的统计
func CalculateByTrie2[K comparable, V any](table *ETable, rowTrie *Trie[K, V], fns ...CalculateByTrieOpts[V]) {
	fnCount := len(fns)
	if !table.IsEmpty() {
		return
	}
	bottomTrieArr := rowTrie.Bottom()
	if len(bottomTrieArr) == 0 || fnCount == 0 {
		return
	}
	eRow := NewERow(table)
	for _, bTire := range bottomTrieArr {
		for _, opt := range fns {
			fn, fn2 := opt.GetVal, opt.GetFnName
			var eCol *ECol
			eCol = NewECol(table)
			cell := NewECell(fn(bTire.List), table, eRow, eCol)
			cell.payload = bTire.List
			eCol.AddCell(cell)
			eCol.Trie = bTire
			eCol.SetName(bTire.GetKey())
			eCol.SetFnName(fn2())
			eRow.AddCell(cell)
		}
	}
	table.eRows = append(table.eRows, eRow)
}

// 行及数据区域的统计
func CalculateByTrie1[K comparable, V any](table *ETable, fns ...CalculateByTrieOpts[V]) {
	fnCount := len(fns)
	if table.IsEmpty() || fnCount == 0 {
		return
	}
	firstRwo := make([]*ECol, fnCount)
	table.ForRow(func(row *ERow, index int) {
		rtRie, ok := row.Trie.(*Trie[K, V])
		if !ok {
			return
		}
		for j, opt := range fns {
			fn, fn2 := opt.GetVal, opt.GetFnName
			var eCol *ECol
			if index == 0 {
				eCol = NewECol(table)
				firstRwo[j] = eCol
			} else {
				eCol = firstRwo[j]
			}
			cell := NewECell(fn(rtRie.List), table, row, eCol)
			cell.payload = rtRie.List
			eCol.AddCell(cell)
			eCol.SetName(fn2())
			eCol.SetFnName(fn2())
			row.AddCell(cell)
		}

	})

}

func Calculate[K comparable, V any](cTrie *Trie[K, V], rTrie *Trie[K, V], c Collection[V], fns ...CalculateByTrieOpts[V]) (et *ETable) {
	et = NewETable()
	if len(fns) == 0 {
		return
	}
	// 行、列、数据区域
	if !cTrie.IsEmpty() && !rTrie.IsEmpty() {
		AddColByTrie(et, cTrie, cTrie.RootName)
		CalculateByTrie[K, V](et, rTrie, fns...)
	} else if !cTrie.IsEmpty() {
		// 列、数据区域
		AddColByTrie(et, cTrie, cTrie.RootName)
		CalculateByTrie1[K, V](et, fns...)
	} else if !rTrie.IsEmpty() {
		// 行及数据区域
		CalculateByTrie2[K, V](et, rTrie, fns...)
	} else {
		// 数据区域的统计
		CalculateByTrie3[K, V](et, c, fns...)
	}

	return
}

// 行、列及数据区域的统计
func CalculateByTrie[K comparable, V any](table *ETable, rowTrie *Trie[K, V], fns ...CalculateByTrieOpts[V]) {
	fnCount := len(fns)
	if table.IsEmpty() || fnCount == 0 {
		return
	}
	bottomTrieArr := rowTrie.Bottom()
	firstRwo := make([]*ECol, fnCount*len(bottomTrieArr))
	table.ForRow(func(row *ERow, index int) {
		rtRie, ok := row.Trie.(*Trie[K, V])
		if !ok {
			return
		}
		for i, bTire := range bottomTrieArr {
			for j, opt := range fns {
				fn, fn2 := opt.GetVal, opt.GetFnName
				var eCol *ECol
				if index == 0 {
					eCol = NewECol(table)
					firstRwo[i*fnCount+j] = eCol
				} else {
					eCol = firstRwo[i*fnCount+j]
				}

				intersection := rtRie.List.Intersection(bTire.List, func(v V) interface{} {
					return v
				})
				cell := NewECell(fn(intersection), table, row, eCol)
				cell.payload = intersection
				eCol.AddCell(cell)
				eCol.Trie = bTire
				eCol.SetName(bTire.GetKey())
				eCol.SetFnName(fn2())
				row.AddCell(cell)
			}
		}
	})
}

func GetFields[T any, V any](c Collection[V], fn func(V) T) (result []T) {
	for _, item := range c {
		result = append(result, fn(item))
	}
	return
}

func Join(arr []string, sep string) string {
	return strings.Trim(strings.Join(arr, sep), sep)
}

func ToJson(value interface{}) (string, error) {
	b, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// 生成行表头
func CreateRowHeaderByECol(et *ETable, fn func(*ECol, int) string, fillColNames ...string) [][]interface{} {
	result := et.GetElementByCol(func(col *ECol, i int) interface{} {
		return fn(col, i)
	})
	if len(fillColNames) > 0 {
		temp := make([]interface{}, len(fillColNames))
		for i := 0; i < len(fillColNames); i++ {
			temp[i] = fillColNames[i]
		}
		result = append(temp, result...)
	}
	return [][]interface{}{result}
}

// 生成行表头
func CreateRowHeaderByEColTrie[K comparable, V any](et *ETable, fn func(*Trie[K, V]) string, fillColNames ...string) [][]interface{} {
	result := et.GetElementByCol(func(col *ECol, i int) interface{} {
		name := col.GetName()
		if itemTrie, ok := col.Trie.(*Trie[K, V]); ok {
			strArr := GetFields(NewCollection(itemTrie.Ancestor()), func(t *Trie[K, V]) string {
				return fn(t)
			})
			name = Join(strArr, "/")
		}
		if col.GetFnName() != "" {
			return name + "/" + col.GetFnName()
		}
		return name
	})
	if len(fillColNames) > 0 {
		temp := make([]interface{}, len(fillColNames))
		for i := 0; i < len(fillColNames); i++ {
			temp[i] = fillColNames[i]
		}
		result = append(temp, result...)
	}
	return [][]interface{}{result}
}

// 获取列字段的名称
func CreateColHeaderByColCell[K comparable, V any](cell *ECell, fn func(*Trie[K, V]) string) string {
	name := cell.eRow.GetName()
	if itemTrie, ok := cell.Trie.(*Trie[K, V]); ok {
		arr := GetFields(NewCollection(itemTrie.Ancestor()), fn)
		name = Join(arr, "/")
	}
	return name
}

// 生成多级列表头
func CreateTreeColHeader[K comparable, V any](et *ETable, fn func(*Trie[K, V]) string) (result [][]interface{}) {
	var temp [][]string
	et.ForRow(func(row *ERow, i int) {
		if itemTrie, ok := row.Trie.(*Trie[K, V]); ok {
			arr := GetFields(NewCollection(itemTrie.Ancestor()), func(t *Trie[K, V]) string {
				return fn(t)
			})
			if len(arr) > 2 {
				temp = append(temp, arr[1:len(arr)-1])
			}
		}
	})

	if len(temp) > 0 {
		result = make([][]interface{}, len(temp))
		for i, arr := range temp {
			result[i] = make([]interface{}, len(arr))
			for j, item := range arr {
				result[i][j] = item
			}
		}

		if len(result) < et.RowNum() {
			for i := et.RowNum() - len(result); i > 0; i-- {
				result = append(result, make([]interface{}, len(result[0])))
			}
		}
	}
	return
}

func Round(f float64, n int) float64 {
	// 计算需要乘以的10的幂次
	pow := math.Pow(10, float64(n))
	// 四舍五入
	rounded := math.Round(f*pow) / pow
	return rounded
}
