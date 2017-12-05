package traffic

// Agent is an agent that travels from a start to a destination at its leave
// time. It is identified by id and group id.
type Agent interface {
	Id() string
	Group() string
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
