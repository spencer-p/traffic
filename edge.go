package traffic

type Edge interface {
	Weight() float64
	To() string
	From() string
	AddAgent()
	RemoveAgent()
	Time() int
}
