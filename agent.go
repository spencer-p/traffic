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
	timeUntilNextChoice int
	history             []Choice
	position            Node
}
