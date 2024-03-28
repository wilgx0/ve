package ve

import "sort"

type Collection[V any] []V

func NewCollection[V any](arr []V) Collection[V] {
	return arr
}

func (collection Collection[V]) Sort(fn func(i, j int) bool) {
	sort.Slice(collection, fn)
}

func (collection Collection[V]) Reverse() {
	i, j := 0, len(collection)-1
	for i < j {
		collection[i], collection[j] = collection[j], collection[i]
		i++
		j--
	}
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

func (collection Collection[V]) Intersection(other Collection[V], fn func(V) interface{}) (result Collection[V]) {
	return collection.Filter(func(value V) bool {
		return other.Filter(func(otherValue V) bool {
			return fn(value) == fn(otherValue)
		}).Len() > 0
	})
}
