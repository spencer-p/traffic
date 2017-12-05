package traffic

// Edge is a representation of a mode of transport between two vertices (usually
// intersections). An Edge should internally keep track of the number of agents
// on it (updated through AddAgent and RemoveAgent), and modify its weight and
// travel time accordingly. Additionally, an Edge can return a different weight
// for different agents.
type Edge interface {
	Weight(Agent) float64
	To() string
	From() string
	AddAgent()
	RemoveAgent()
	Time() int
}
