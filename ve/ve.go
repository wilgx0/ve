package ve

type VE struct {
	Body     VTable
	Head     VRow
	Foot     VRow
	FirstCol VTable
	LastCol  VTable
}

func NewVE() *VE {
	return &VE{}
}
