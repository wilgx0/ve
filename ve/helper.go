package ve

import (
	"encoding/json"
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

func AddColByTrie[K comparable, V any](table *ETable, trie *Trie[K, V], name ...string) *ECol {
	col := NewEColByTrie(trie, table, name...)
	table.AddCol(col)
	return col
}

type CalculateByTrieOpts[V any] struct {
	GetVal    func(Collection[V]) interface{}
	GetFnName func() string
}

func CalculateByTrie[K comparable, V any](table *ETable, trie *Trie[K, V], fns ...CalculateByTrieOpts[V]) {
	fnCount := len(fns)
	if table.IsEmpty() || fnCount == 0 {
		return
	}
	bottomTrieArr := trie.Bottom()
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

				cell := NewECell(fn(rtRie.List.Intersection(bTire.List, func(v V) interface{} {
					return v
				})), table, row, eCol)
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

func CreateRowHeaderByECol[K comparable, V any](et *ETable, fn func(*Trie[K, V]) string, fillColNames ...string) [][]interface{} {
	result := et.GetElementByCol(func(col *ECol, i int) interface{} {
		name := col.GetName()
		if itemTrie, ok := col.Trie.(*Trie[K, V]); ok {
			f := GetFields(NewCollection(itemTrie.Ancestor()), func(t *Trie[K, V]) string {
				return fn(t)
			})
			name = Join(f, "/")
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

func CreateColHeaderByColCell[K comparable, V any](cell *ECell, fn func(*Trie[K, V]) string) string {
	name := cell.eRow.GetName()
	if itemTrie, ok := cell.Trie.(*Trie[K, V]); ok {
		arr := GetFields(NewCollection(itemTrie.Ancestor()), fn)
		name = Join(arr, "/")
	}
	return name
}

// 生成列表头
func CreateTreeColHeader[K comparable, V any](et *ETable) (result [][]interface{}, colNames []string) {
	var temp [][]string
	et.ForRow(func(row *ERow, i int) {
		if itemTrie, ok := row.Trie.(*Trie[K, V]); ok {
			arr := GetFields(NewCollection(itemTrie.Ancestor()), func(t *Trie[K, V]) string {
				return t.GetKey()
			})
			if len(arr) > 2 {
				temp = append(temp, arr[1:len(arr)-1])
				if colNames == nil {
					colNames = GetFields(NewCollection(itemTrie.Ancestor()), func(t *Trie[K, V]) string {
						return t.Name
					})
					colNames = colNames[0 : len(colNames)-2]
				}
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
