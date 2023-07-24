package ve

type StatGroup map[interface{}]*VCell

func NewStatGroup() StatGroup {
	return make(StatGroup)
}

func (sg StatGroup) IsEmpty() bool {
	return len(sg) == 0
}

func (sg StatGroup) VLookup(findColumn VColumn) (result VColumn) {
	if sg.IsEmpty() {
		return
	}

	for _, item := range findColumn {
		if value, ok := sg[item.Val()]; ok {
			result = append(result, value)
		} else {
			result = append(result, nil)
		}
	}

	return
}
