package communication

type Server interface {
	Agents() []ServerAgent
	UpdateAgentsStatus()
}

type ServerAgent interface {
	ID() int
	SendHeartbeat() error
}
