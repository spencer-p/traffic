package traffic

type Edge interface {
	Weight() float64
	To() string
	From() string
}

type SimpleEdge struct {
	weight   float64
	from, to string
}

func NewSimpleEdge(weight float64, from, to string) Edge {
	return SimpleEdge{weight: weight, to: to, from: from}
}

func (e SimpleEdge) Weight() float64 {
	return e.weight
}

func (e SimpleEdge) To() string {
	return e.to
}

func (e SimpleEdge) From() string {
	return e.from
}
