package ve

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

type Student struct {
	Name        string `json:"姓名"`   // 姓名
	Class       string `json:"班级"`   // 班级
	Grade       int    `json:"年级"`   // 年级
	Subject     string `json:"科目"`   // 科目
	SubjectType string `json:"科目类别"` // 科目类别
	Score       int    `json:"分数"`   // 分数
	Sex         string `json:"性别"`   // 性别
	Term        string `json:"学期"`   // 学期
	AddScore    int    `json:"加分"`   // 加分
}

func getTestData() []*Student {
	file, err := os.Open("example.json") // 替换为你的JSON文件路径
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil
	}
	var data []*Student
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return nil
	}
	return data
}

func TestNewMap1(t *testing.T) {
	collection := NewCollection(getTestData())
	map1 := NewMap(collection, func(student *Student) string {
		return student.Class
	})

	if map1.Len() != collection.Unique(func(student *Student) string {
		return student.Class
	}).Len() {
		t.Error("map1 length error")
	}

}

func TestNewMap2(t *testing.T) {
	collection := NewCollection(getTestData())

	map2 := NewMap(collection, func(student *Student) int {
		return student.Grade
	})

	if map2.Len() != collection.Unique(func(student *Student) string {
		return String(student.Grade)
	}).Len() {
		t.Error("map2 length error")
	}

	if len(map2.GetKeys()) != collection.Unique(func(student *Student) string {
		return String(student.Grade)
	}).Len() {
		t.Error("map2 keys length error")
	}

	c := map2.GetValue(1)
	if value := c.SumUint64(func(student *Student) uint64 {
		return uint64(student.Score)
	}); value != collection.Filter(func(student *Student) bool {
		return student.Grade == 1
	}).SumUint64(func(student *Student) uint64 {
		return uint64(student.Score)
	}) {
		t.Error("map2 value error")
	}

}

func TestSum1(t *testing.T) {
	collection := NewCollection(getTestData())

	// 按班级分组
	groupByClass := NewMap(collection, func(student *Student) string {
		return student.Class
	})

	et := NewETable()

	// 添加首列
	AddCol(et, groupByClass.GetKeys(), "班级")

	// 添加各行的名称及排序
	firstCol := map[string][]interface{}{
		"一班": {"一班", 1},
		"二班": {"二班", 2},
		"三班": {"三班", 3},
	}

	et.ForRow(func(row *ERow, i int) {
		row.SetName(firstCol[row.FirstCell().String()][0].(string)).Sort = firstCol[row.FirstCell().String()][1].(int)
	})

	// 对行进行排序
	et.SortRow(func(row *ERow, row2 *ERow) bool {
		return row.Sort < row2.Sort
	})
	// 统计列
	groupByClass.VLookup(et, func(row *ERow, i int) string {
		return row.FirstCell().String()
	}, func(c Collection[*Student]) interface{} {
		return c.Filter(func(student *Student) bool {
			return student.Grade == 1
		}).SumUint64(func(student *Student) uint64 {
			return uint64(student.Score)
		})
	}, "一年级")
	groupByClass.VLookup(et, func(row *ERow, i int) string {
		return row.FirstCell().String()
	}, func(c Collection[*Student]) interface{} {
		return c.Filter(func(student *Student) bool {
			return student.Grade == 2
		}).SumUint64(func(student *Student) uint64 {
			return uint64(student.Score)
		})
	}, "二年级")
	groupByClass.VLookup(et, func(row *ERow, i int) string {
		return row.FirstCell().String()
	}, func(c Collection[*Student]) interface{} {
		return c.Filter(func(student *Student) bool {
			return student.Grade == 3
		}).SumUint64(func(student *Student) uint64 {
			return uint64(student.Score)
		})
	}, "三年级")

	// 合计列
	AddColByFn(et, func(row *ERow, i int) interface{} {
		return row.SumUint64(func(cell *ECell, i int) bool {
			return i != 0
		})
	}, "合计")

	// 合计行
	AddRowByFn(et, func(col *ECol, i int) interface{} {
		if i == 0 {
			return nil
		}
		return col.SumUint64()
	}, "合计")

	// 展示数据
	showData := et.ToArr(func(cell *ECell, rN int, cN int) interface{} {
		if cN == 0 {
			return cell.eRow.GetName()
		} else {
			return cell.Val()
		}
	})

	// 生成表头
	showData = append([][]interface{}{et.GetElementByCol(func(col *ECol, i int) interface{} {
		return col.GetName()
	})}, showData...)
	t.Log(showData)

	total := et.GetCell(et.RowNum()-1, et.ColNum()-1)

	if total.Uint64() != collection.SumUint64(func(student *Student) uint64 { return uint64(student.Score) }) {
		t.Error("total error")
	}
}

func TestSum2(t *testing.T) {
	collection := NewCollection(getTestData())
	groupByClass := NewMap(collection, func(student *Student) int {
		return student.Grade
	})

	et := NewETable()

	// 添加首列
	AddCol(et, groupByClass.GetKeys(), "年级")

	// 添加各行的名称及排序
	firstCol := map[int][]interface{}{
		1: {"一年级", 1},
		2: {"二年级", 2},
		3: {"三年级", 3},
	}
	for _, cell := range et.GetCol(0).Cells() {
		cell.GetERow().SetName(firstCol[cell.Int()][0].(string)).Sort = firstCol[cell.Int()][1].(int)
	}

	//  填充汇总各列的数据
	AddColByFn(et, func(row *ERow, i int) interface{} {
		return groupByClass.GetValue(row.FirstCell().Int()).Filter(func(student *Student) bool {
			return student.Class == "一班"
		}).SumUint64(func(student *Student) uint64 {
			return uint64(student.Score)
		})
	}, "一班")

	AddColByFn(et, func(row *ERow, i int) interface{} {
		return groupByClass.GetValue(row.FirstCell().Int()).Filter(func(student *Student) bool {
			return student.Class == "二班"
		}).SumUint64(func(student *Student) uint64 {
			return uint64(student.Score)
		})
	}, "二班")

	AddColByFn(et, func(row *ERow, i int) interface{} {
		return groupByClass.GetValue(row.FirstCell().Int()).Filter(func(student *Student) bool {
			return student.Class == "三班"
		}).SumUint64(func(student *Student) uint64 {
			return uint64(student.Score)
		})
	}, "三班")

	// 最后一列的合计
	AddColByFn(et, func(row *ERow, i int) interface{} {
		return row.SumUint64(func(cell *ECell, i int) bool {
			return i != 0
		})
	}, "合计")

	// 对行进行排序
	et.SortRow(func(row *ERow, row2 *ERow) bool {
		return row.Sort < row2.Sort
	})

	// 最后一行的合计
	AddRowByFn(et, func(col *ECol, i int) interface{} {
		if i == 0 {
			return nil
		}
		return col.SumUint64()
	}, "合计")

	showData := et.ToArr(func(cell *ECell, rN int, cN int) interface{} {
		if cN == 0 {
			return cell.eRow.GetName()
		} else {
			return cell.Val()
		}
	})
	// 生成表头
	showData = append([][]interface{}{et.GetElementByCol(func(col *ECol, i int) interface{} {
		return col.GetName()
	})}, showData...)
	t.Log(showData)
	jsonStr, _ := ToJson(showData)
	t.Log(jsonStr)

	total := et.GetCell(et.RowNum()-1, et.ColNum()-1)

	if total.Uint64() != collection.SumUint64(func(student *Student) uint64 { return uint64(student.Score) }) {
		t.Error("total error")
	}
}

// 表格合并
func TestMergeTable(t *testing.T) {
	et1 := NewETable()
	AddRow(et1, []interface{}{"一年级", 5, 11, 0, 16}, "一年级")
	AddRow(et1, []interface{}{"二年级", 0, 7, 1, 8}, "二年级")
	AddRow(et1, []interface{}{"三年级", 0, 0, 2, 2}, "三年级")
	AddRow(et1, []interface{}{"总计", 5, 18, 3, 26}, "总计")

	et2 := NewETable()
	AddRow(et2, []interface{}{"一班", 5, 0, 0, 5}, "一班")
	AddRow(et2, []interface{}{"二班", 11, 7, 0, 18}, "二班")
	AddRow(et2, []interface{}{"三班", 0, 1, 2, 3}, "三班")
	AddRow(et2, []interface{}{"总计", 16, 8, 2, 26}, "总计")
	et1.Merge(et2)
	// t.Log(et1)
	// 计算总计
	AddRowByFn(et1, func(col *ECol, i int) interface{} {
		if i == 0 {
			return nil
		}
		return col.SumUint64(func(cell *ECell, i int) bool {
			return cell.GetERow().GetName() == "总计"
		})
	}, "总合计")

	total := et1.GetCell(et1.RowNum()-1, et1.ColNum()-1)
	if total.Uint64() != 52 {
		t.Error("total error")
	}

	showData := et1.ToArr(func(cell *ECell, rN int, cN int) interface{} {
		if cN == 0 {
			return cell.eRow.GetName()
		} else {
			return cell.Val()
		}
	})
	t.Log(showData)
}

func TestTrie1(t *testing.T) {
	collection := NewCollection(getTestData())
	trie := NewTrie[string, *Student]("")

	trie.Insert(collection, TrieInsertOpt[string, *Student]{
		Fn: func(student *Student) string {
			return student.Class
		},
		GetName: func() string {
			return "班级"
		},
	}, TrieInsertOpt[string, *Student]{
		Fn: func(student *Student) string {
			return student.Name
		},
		GetName: func() string {
			return "姓名"
		},
	})

	et := NewETable()

	AddColByTrie(et, trie, "班级/姓名")
	AddColByFn(et, func(row *ERow, i int) interface{} {
		if itemTrie, ok := row.FirstCell().Trie.(*Trie[string, *Student]); ok {
			return itemTrie.List.SumUint64(func(student *Student) uint64 {
				return uint64(student.Score)
			})
		} else {
			return uint64(0)
		}
	}, "分数")

	AddRowByFn(et, func(col *ECol, i int) interface{} {
		if i == 0 {
			return nil
		}
		return col.SumUint64()
	}, "合计")

	showData := et.ToArr(func(cell *ECell, rN int, cN int) interface{} {
		if cN == 0 {
			if itemTrie, ok := cell.Trie.(*Trie[string, *Student]); ok {
				f := GetFields(NewCollection(itemTrie.Ancestor()), func(t *Trie[string, *Student]) string {
					return t.Key
				})
				return Join(f, "/")
			} else {
				return ""
			}
		} else {
			return cell.Val()
		}
	})

	// 生成表头
	showData = append([][]interface{}{et.GetElementByCol(func(col *ECol, i int) interface{} {
		return col.GetName()
	})}, showData...)
	t.Log(showData)
	jsonStr, _ := ToJson(showData)
	t.Log(jsonStr)
	total := et.GetCell(et.RowNum()-1, et.ColNum()-1)

	if total.Uint64() != collection.SumUint64(func(student *Student) uint64 { return uint64(student.Score) }) {
		t.Error("total error")
	}
}

func TestTrie2(t *testing.T) {
	collection := NewCollection(getTestData())
	cTrie := NewTrie[string, *Student]("")

	//列： 班级 / 姓名
	cTrie.Insert(collection, TrieInsertOpt[string, *Student]{
		Fn: func(student *Student) string {
			return student.Class
		},
		GetName: func() string {
			return "班级"
		},
	}, TrieInsertOpt[string, *Student]{
		Fn: func(student *Student) string {
			return student.Name
		},
		GetName: func() string {
			return "姓名"
		},
	})

	//行: 科目类别 / 科目
	rTrie := NewTrie[string, *Student]("")
	rTrie.Insert(collection, TrieInsertOpt[string, *Student]{
		Fn: func(student *Student) string {
			return student.SubjectType
		},
		GetName: func() string {
			return "科目类别"
		},
	}, TrieInsertOpt[string, *Student]{
		Fn: func(student *Student) string {
			return student.Subject
		},
		GetName: func() string {
			return "科目"
		},
	})

	et := NewETable()

	// 添加首列
	AddColByTrie(et, cTrie, "班级/姓名")

	// 计算
	//CalculateByTrie(et, rTrie, func(c Collection[*Student]) interface{} {
	//	return c.SumUint64(func(student *Student) uint64 {
	//		return uint64(student.Score)
	//	})
	//})
	CalculateByTrie(et, rTrie, CalculateByTrieOpts[*Student]{
		GetVal: func(c Collection[*Student]) interface{} {
			return c.SumUint64(func(student *Student) uint64 {
				return uint64(student.Score)
			})
		},
		GetFnName: func() string {
			return "分数"
		},
	})
	// 最后一列的合计
	AddColByFn(et, func(row *ERow, i int) interface{} {
		return row.SumUint64(func(cell *ECell, i int) bool {
			return i != 0
		})
	}, "合计")

	// 最后一行的合计
	AddRowByFn(et, func(col *ECol, i int) interface{} {
		if i == 0 {
			return nil
		}
		return col.SumUint64()
	}, "合计")

	showData := et.ToArr(func(cell *ECell, rN int, cN int) interface{} {
		if cN == 0 {
			//  生成列头
			name := cell.eRow.GetName()
			if itemTrie, ok := cell.Trie.(*Trie[string, *Student]); ok {
				f := GetFields(NewCollection(itemTrie.Ancestor()), func(t *Trie[string, *Student]) string {
					return t.Key
				})
				name = Join(f, "/")
			}
			return name
		} else {
			return cell.Val()
		}
	})

	// 生成表头
	showData = append([][]interface{}{et.GetElementByCol(func(col *ECol, i int) interface{} {
		name := col.GetName()
		if itemTrie, ok := col.Trie.(*Trie[string, *Student]); ok {
			f := GetFields(NewCollection(itemTrie.Ancestor()), func(t *Trie[string, *Student]) string {
				return t.Key
			})
			name = Join(f, "/")
		}
		return name
	})}, showData...)

	//t.Log(showData)
	//jsonStr, _ := ToJson(showData)
	//t.Log(jsonStr)

	total := et.GetCell(et.RowNum()-1, et.ColNum()-1)

	if total.Uint64() != collection.SumUint64(func(student *Student) uint64 { return uint64(student.Score) }) {
		t.Error("total error")
	}
}

func TestTrie3(t *testing.T) {
	collection := NewCollection(getTestData())
	cTrie := NewTrie[string, *Student]("")

	//列: 学期 / 班级 / 姓名
	cTrie.Insert(collection, TrieInsertOpt[string, *Student]{
		Fn: func(student *Student) string {
			return student.Term
		},
		GetName: func() string {
			return "学期"
		},
	}, TrieInsertOpt[string, *Student]{
		Fn: func(student *Student) string {
			return student.Class
		},
		GetName: func() string {
			return "班级"
		},
	}, TrieInsertOpt[string, *Student]{
		Fn: func(student *Student) string {
			return student.Name
		},
		GetName: func() string {
			return "姓名"
		},
	})

	//行: 科目类别 / 科目
	rTrie := NewTrie[string, *Student]("")
	rTrie.Insert(collection, TrieInsertOpt[string, *Student]{
		Fn: func(student *Student) string {
			return student.SubjectType
		},
		GetName: func() string {
			return "科目类别"
		},
	}, TrieInsertOpt[string, *Student]{
		Fn: func(student *Student) string {
			return student.Subject
		},
		GetName: func() string {
			return "科目"
		},
	})

	et := NewETable()

	// 添加首列
	AddColByTrie(et, cTrie, "学期/班级/姓名")

	// 计算值 : 分数 / 加分项
	CalculateByTrie(et, rTrie, CalculateByTrieOpts[*Student]{
		GetVal: func(c Collection[*Student]) interface{} {
			return c.SumUint64(func(student *Student) uint64 {
				return uint64(student.Score)
			})
		},
		GetFnName: func() string {
			return "分数"
		},
	}, CalculateByTrieOpts[*Student]{
		GetVal: func(c Collection[*Student]) interface{} {
			return c.SumUint64(func(student *Student) uint64 {
				return uint64(student.AddScore)
			})
		},
		GetFnName: func() string {
			return "加分"
		},
	})

	// 最后一列的合计
	AddColByFn(et, func(row *ERow, i int) interface{} {
		return row.SumUint64(func(cell *ECell, i int) bool {
			return i != 0
		})
	}, "合计")

	// 最后一行的合计
	AddRowByFn(et, func(col *ECol, i int) interface{} {
		if i == 0 {
			return nil
		}
		return col.SumUint64()
	}, "合计")

	showData := NewGrid(et.ToArr(func(cell *ECell, rN int, cN int) interface{} {
		if cN == 0 {
			//  生成列头
			return CreateColHeaderByColCell(cell, func(t *Trie[string, *Student]) string {
				return t.Key
			})
		} else {
			return cell.Val()
		}
	}))

	// 生成行表头
	showData = showData.MergeTop(NewGrid(CreateRowHeaderByECol(et, func(t *Trie[string, *Student]) string {
		return t.Key
	})))

	//t.Log(showData)
	//jsonStr, _ := ToJson(showData)
	//t.Log(jsonStr)

	total := et.GetCell(et.RowNum()-1, et.ColNum()-1)
	t.Logf("total: %d", total.Uint64())
	if total.Uint64() != collection.SumUint64(func(student *Student) uint64 { return uint64(student.Score + student.AddScore) }) {
		t.Error("total error")
	}
}

// 生成表头
func TestHeaderByTrie(t *testing.T) {
	collection := NewCollection(getTestData())
	cTrie := NewTrie[string, *Student]("")

	//列: 学期 / 班级 / 姓名
	cTrie.Insert(collection, TrieInsertOpt[string, *Student]{
		Fn: func(student *Student) string {
			return student.Term
		},
		GetName: func() string {
			return "学期"
		},
	}, TrieInsertOpt[string, *Student]{
		Fn: func(student *Student) string {
			return student.Class
		},
		GetName: func() string {
			return "班级"
		},
	}, TrieInsertOpt[string, *Student]{
		Fn: func(student *Student) string {
			return student.Name
		},
		GetName: func() string {
			return "姓名"
		},
	})

	//行: 科目类别 / 科目
	rTrie := NewTrie[string, *Student]("")
	rTrie.Insert(collection, TrieInsertOpt[string, *Student]{
		Fn: func(student *Student) string {
			return student.SubjectType
		},
		GetName: func() string {
			return "科目类别"
		},
	}, TrieInsertOpt[string, *Student]{
		Fn: func(student *Student) string {
			return student.Subject
		},
		GetName: func() string {
			return "科目"
		},
	})

	et := NewETable()

	// 添加首列
	//学期/班级/姓名
	AddColByTrie(et, cTrie, "姓名")

	// 计算值 : 分数 / 加分项
	CalculateByTrie(et, rTrie, CalculateByTrieOpts[*Student]{
		GetVal: func(c Collection[*Student]) interface{} {
			return c.SumUint64(func(student *Student) uint64 {
				return uint64(student.Score)
			})
		},
		GetFnName: func() string {
			return "分数"
		},
	}, CalculateByTrieOpts[*Student]{
		GetVal: func(c Collection[*Student]) interface{} {
			return c.SumUint64(func(student *Student) uint64 {
				return uint64(student.AddScore)
			})
		},
		GetFnName: func() string {
			return "加分"
		},
	})

	// 最后一列的合计
	AddColByFn(et, func(row *ERow, i int) interface{} {
		return row.SumUint64(func(cell *ECell, i int) bool {
			return i != 0
		})
	}, "合计")

	// 最后一行的合计
	AddRowByFn(et, func(col *ECol, i int) interface{} {
		if i == 0 {
			return nil
		}
		return col.SumUint64()
	}, "合计")

	showData := NewGrid(et.ToArr(func(cell *ECell, rN int, cN int) interface{} {
		if cN == 0 {
			return cell.GetERow().GetName()
		} else {
			return cell.Val()
		}
	}))

	// 生成并组合列表头
	treeColHeader, treeColHeaderNames := CreateTreeColHeader[string, *Student](et)
	showData = showData.MergeLeft(NewGrid(treeColHeader))

	// 生成并组合行表头
	showData = showData.MergeTop(NewGrid(CreateRowHeaderByECol(et, func(t *Trie[string, *Student]) string {
		return t.Key
	}, treeColHeaderNames...)))

	//t.Log(showData)
	//jsonStr, _ := ToJson(showData)
	//t.Log(jsonStr)

	total := et.GetCell(et.RowNum()-1, et.ColNum()-1)
	t.Logf("total: %d", total.Uint64())
	if total.Uint64() != collection.SumUint64(func(student *Student) uint64 { return uint64(student.Score + student.AddScore) }) {
		t.Error("total error")
	}

}
