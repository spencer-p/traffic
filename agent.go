package traffic

type Agent interface {
	Id() string
	Start() string
	Destination() string
	LeaveTime() int
}

type metaAgent struct {
	agent               Agent
	path                []Step
	pathLastUpdated     int
	lastEdge            Edge
	timeUntilNextChoice int
	history             []Choice
	position            Node
}
