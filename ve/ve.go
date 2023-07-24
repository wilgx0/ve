package ve

type VE struct {
	Body     VTable
	Head     VRow
	Foot     VRow
	LeftCol  VColumn
	RightCol VColumn
}

func NewVE() *VE {
	return &VE{}
}

func (ve *VE) SplitHead() {
	ve.Head = ve.Body.FirstRow()
}

func (ve *VE) SplitFoot() {
	ve.Foot = ve.Body.LastRow()
}

func (ve *VE) SplitLeftCol() {
	ve.LeftCol = ve.Body.FirstCol()
}

func (ve *VE) SplitRightCol() {
	ve.RightCol = ve.Body.LastCol()
}

func (ve *VE) MergeHead() (err error) {
	if ve.Head.IsEmpty() {
		return
	}
	newVt, err := CreateVTableByRow(ve.Head)
	if err != nil {
		return
	}
	_, err = newVt.VerticalMerge(ve.Body)
	return
}

func (ve *VE) MergeFoot() (err error) {
	if ve.Foot.IsEmpty() {
		return
	}
	_, err = ve.Body.AddRow(ve.Foot)
	return
}

func (ve *VE) MergeLeftCol() (err error) {
	if ve.LeftCol.IsEmpty() {
		return
	}
	newVt, err := CreateVTableByCol(ve.LeftCol)
	if err != nil {
		return
	}
	_, err = newVt.HorizontalMerge(ve.Body)
	return
}

func (ve *VE) MergeRightCol() (err error) {
	if ve.RightCol.IsEmpty() {
		return
	}
	_, err = ve.Body.AddColumn(ve.RightCol)
	return
}
