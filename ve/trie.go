package ve

type Trie[K comparable, V any] struct {
	Children map[K]*Trie[K, V]
	Last     *Trie[K, V]
	List     Collection[V]
	Key      K
	Name     string
	RootName string
}

func NewTrie[K comparable, V any](key K) *Trie[K, V] {
	return &Trie[K, V]{
		Key:      key,
		Children: make(map[K]*Trie[K, V]),
	}
}

func (t *Trie[K, V]) IsEmpty() bool {
	return t == nil || len(t.Children) == 0
}

func (t *Trie[K, V]) GetKey() string {
	return String(t.Key)
}

func (t *Trie[K, V]) GetChildren() (result []*Trie[K, V]) {
	for _, child := range t.Children {
		result = append(result, child)
	}
	return
}

func (t *Trie[K, V]) insert(list Collection[V], opt TrieInsertOpt[K, V]) {
	fn, getName := opt.Fn, opt.GetName
	for _, value := range list {
		node, ok := t.Children[fn(value)]
		if !ok {
			node = NewTrie[K, V](fn(value))
			t.Children[fn(value)] = node
		}
		node.Name = getName(value)
		node.List = append(node.List, value)
	}
}

type TrieInsertOpt[K comparable, V any] struct {
	Fn      func(V) K
	GetName func(V) string
}

func (t *Trie[K, V]) Insert(list Collection[V], fns ...TrieInsertOpt[K, V]) {
	if list.IsEmpty() || len(fns) == 0 {
		return
	}
	queue := []*Trie[K, V]{t}
	for i, opt := range fns {
		for ql := len(queue); ql > 0; ql-- {
			var node *Trie[K, V]
			node, queue = queue[0], queue[1:]
			if i == 0 {
				node.insert(list, opt)
			} else {
				node.insert(node.List, opt)
			}
			for _, child := range node.Children {
				child.Last = node
				queue = append(queue, child)
			}
		}
	}
}

func (t *Trie[K, V]) Bottom() (result []*Trie[K, V]) {
	queue := []*Trie[K, V]{t}
	for len(queue) > 0 {
		for ql := len(queue); ql > 0; ql-- {
			var node *Trie[K, V]
			node, queue = queue[0], queue[1:]
			if len(node.Children) > 0 {
				queue = append(queue, node.GetChildren()...)
			} else {
				result = append(result, node)
			}
		}
	}

	return
}

func (t *Trie[K, V]) Ancestor() (result []*Trie[K, V]) {
	node := t
	for node != nil {
		result = append(result, node)
		node = node.Last
	}
	if len(result) > 0 {
		NewCollection(result).Reverse()
	}
	return
}
