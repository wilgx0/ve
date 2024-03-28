package ve

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

func (m Map[K, V]) VLookup(table *ETable, keyFn func(*ERow, int) K, ValueFn func(Collection[V]) interface{}, name ...string) {
	AddColByFn(table, func(row *ERow, i int) interface{} {
		return ValueFn(m.GetValue(keyFn(row, i)))
	}, name...)
}
