package traffic

type Edge interface {
	Weight(Agent) float64
	To() string
	From() string
	AddAgent()
	RemoveAgent()
	Time() int
}
