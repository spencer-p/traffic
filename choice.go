package traffic

type Choice struct {
	To         string
	EdgeType   string
	Timestamp  int
	TravelTime int
}

type AverageData struct {
	AgentCount             int
	MinutesTraveled        int
	AverageDeltaTravelTime float64
}

func (ad *AverageData) UpdateAverage() {
	ad.AverageDeltaTravelTime = float64(ad.MinutesTraveled) / float64(ad.AgentCount)
}
