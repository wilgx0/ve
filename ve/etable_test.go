package ve

import "testing"

type Student struct {
	Name  string
	Class string
	Grade int
	Score int
}

func getTestData() []*Student {
	return []*Student{
		{
			Name:  "小明",
			Class: "一班",
			Grade: 1,
			Score: 2,
		},
		{
			Name:  "小张",
			Class: "一班",
			Grade: 1,
			Score: 3,
		},
		{
			Name:  "小王",
			Class: "二班",
			Grade: 1,
			Score: 5,
		},
		{
			Name:  "小李",
			Class: "二班",
			Grade: 1,
			Score: 6,
		},
		{
			Name:  "小赵",
			Class: "二班",
			Grade: 2,
			Score: 7,
		},
		{
			Name:  "小刘",
			Class: "三班",
			Grade: 2,
			Score: 1,
		},
	}

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
	groupByClass := NewMap(collection, func(student *Student) string {
		return student.Class
	})

	et := NewETable()
	AddCol(et, groupByClass.GetKeys(), "class")
	AddColByFn(et, func(row *ERow, i int) interface{} {
		return groupByClass.GetValue(row.FirstCell().String()).Filter(func(student *Student) bool {
			return student.Grade == 1
		}).SumUint64(func(student *Student) uint64 {
			return uint64(student.Score)
		})
	}, "grade1")

	AddColByFn(et, func(row *ERow, i int) interface{} {
		return groupByClass.GetValue(row.FirstCell().String()).Filter(func(student *Student) bool {
			return student.Grade == 2
		}).SumUint64(func(student *Student) uint64 {
			return uint64(student.Score)
		})
	}, "grade2")

	AddColByFn(et, func(row *ERow, i int) interface{} {
		return row.SumUint64()
	}, "total")

	AddRowByFn(et, func(col *ECol, i int) interface{} {
		if i == 0 {
			return nil
		}
		return col.SumUint64()
	})

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

	AddCol(et, groupByClass.GetKeys(), "grade")

	AddColByFn(et, func(row *ERow, i int) interface{} {
		return groupByClass.GetValue(row.FirstCell().Int()).Filter(func(student *Student) bool {
			return student.Class == "一班"
		}).SumUint64(func(student *Student) uint64 {
			return uint64(student.Score)
		})
	}, "class1")

	AddColByFn(et, func(row *ERow, i int) interface{} {
		return groupByClass.GetValue(row.FirstCell().Int()).Filter(func(student *Student) bool {
			return student.Class == "二班"
		}).SumUint64(func(student *Student) uint64 {
			return uint64(student.Score)
		})
	}, "class2")

	AddColByFn(et, func(row *ERow, i int) interface{} {
		return groupByClass.GetValue(row.FirstCell().Int()).Filter(func(student *Student) bool {
			return student.Class == "三班"
		}).SumUint64(func(student *Student) uint64 {
			return uint64(student.Score)
		})
	}, "class3")

	AddColByFn(et, func(row *ERow, i int) interface{} {
		return row.SumUint64(func(cell *ECell, i int) bool {
			return i != 0
		})
	}, "total")

	AddRowByFn(et, func(col *ECol, i int) interface{} {
		if i == 0 {
			return nil
		}
		return col.SumUint64()
	})

	total := et.GetCell(et.RowNum()-1, et.ColNum()-1)

	if total.Uint64() != collection.SumUint64(func(student *Student) uint64 { return uint64(student.Score) }) {
		t.Error("total error")
	}
}
