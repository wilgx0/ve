package ve

import "sort"

func AddRow[T any](table *ETable, arr []T, name ...string) {
	table.AddRow(NewERowByArr(arr, table, name...))
}

func AddRowByFn(table *ETable, fn func(*ECol, int) interface{}, name ...string) {
	table.AddRowByFn(fn, name...)
}

func AddCol[T any](table *ETable, arr []T, name ...string) {
	table.AddCol(NewEColByArr(arr, table, name...))
}

func AddColByFn(table *ETable, fn func(*ERow, int) interface{}, name ...string) {
	table.AddColByFn(fn, name...)
}

type Collection[V any] []V

func NewCollection[V any](arr []V) Collection[V] {
	return arr
}

func (collection Collection[V]) Sort(fn func(i, j int) bool) {
	sort.Slice(collection, fn)
}

func (collection Collection[V]) IsEmpty() bool {
	return len(collection) == 0
}

func (collection Collection[V]) Len() int {
	return len(collection)
}

func (collection Collection[V]) Filter(fn func(V) bool) (result Collection[V]) {
	for _, item := range collection {
		if fn(item) {
			result = append(result, item)
		}
	}
	return
}

func (collection Collection[V]) GetColumn(fn func(V) string) (result []string) {
	for _, value := range collection {
		result = append(result, fn(value))
	}
	return
}

func (collection Collection[V]) SumFloat64(fn func(V) float64) (result float64) {
	for _, value := range collection {
		result += fn(value)
	}
	return
}

func (collection Collection[V]) SumUint64(fn func(V) uint64) (result uint64) {
	for _, value := range collection {
		result += fn(value)
	}
	return
}

func (collection Collection[V]) Count(fn func(V) bool) int {
	return collection.Filter(fn).Len()
}

func (collection Collection[V]) Unique(fn func(V) string) (result Collection[V]) {
	m := map[string]struct{}{}
	for _, value := range collection {
		key := fn(value)
		if _, ok := m[key]; !ok {
			m[key] = struct{}{}
			result = append(result, value)
		}
	}
	return
}

type Map[K comparable, V any] map[K]Collection[V]

func NewMap[K comparable, V any](arr Collection[V], keyFn func(V) K) Map[K, V] {
	m := make(Map[K, V])
	for _, value := range arr {
		m[keyFn(value)] = append(m[keyFn(value)], value)
	}
	return m
}

func (m Map[K, V]) Len() int {
	return len(m)
}

func (m Map[K, V]) GetValue(key K) Collection[V] {
	return m[key]
}

func (m Map[K, V]) GetKeys() []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
