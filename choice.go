package traffic

// Choice represents a choice an agent made while travelling. Slices of Choices
// are encoded in the output of a simulation.
type Choice struct {
	To         string
	EdgeType   string
	Timestamp  int
	TravelTime int
}

// AverageData aggregates agent count and total time sums, then computes average
// travel times for simulation output analysis.
type AverageData struct {
	AgentCount             int
	AgentCountPerEdge      int
	MinutesTraveled        int
	AverageDeltaTravelTime float64
}

func (ad *AverageData) UpdateAverage() {
	ad.AverageDeltaTravelTime = float64(ad.MinutesTraveled) / float64(ad.AgentCount)
}
