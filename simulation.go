package traffic

import (
	"encoding/json"
	"io"
	"log"
	"runtime"
	"sync"
)

type Simulation struct {
	graph          *Graph
	agents         []*metaAgent
	currentTime    int
	pathTimeout    int
	finishedAgents map[string]bool
}

func NewSimulation() Simulation {
	return Simulation{graph: NewGraph(), finishedAgents: make(map[string]bool)}
}

func (S *Simulation) AddEdge(e Edge) {
	S.graph.AddEdge(e)
}

func (S *Simulation) AddAgent(a Agent) {
	var ma metaAgent
	ma.agent = a
	ma.pathLastUpdated = -S.pathTimeout
	ma.position = S.graph.nodes[a.Start()]
	ma.timeUntilNextChoice = a.LeaveTime()
	S.agents = append(S.agents, &ma)
}

func (S *Simulation) Simulate() {
	for len(S.finishedAgents) < len(S.agents) {
		S.Tick()
	}
}

func (S *Simulation) Tick() {
	S.RefreshPaths()
	S.MoveAgents()
	S.currentTime++
}

func (S *Simulation) RefreshPaths() {
	// Map of nodes to start at and agents that start there
	jobs := make(map[Node][]*metaAgent)
	for _, agent := range S.agents {
		if !S.finishedAgents[agent.agent.Id()] &&
			S.currentTime-agent.pathLastUpdated >= S.pathTimeout {
			// If timed out, append agent to the corresponding node list
			agent := agent
			jobs[agent.position] = append(jobs[agent.position], agent)
		}
	}

	// Channel to execute searches and update positions of the agents
	searchCh := make(chan []*metaAgent)
	var searchWG sync.WaitGroup

	// Worker function
	worker := func(jobCh <-chan []*metaAgent) {
		for agentList := range jobCh {
			// Perform the shortest spanning tree search
			S.graph.Dijkstra(agentList[0].position.Name())

			// Set each agent's path
			for _, agent := range agentList {
				var err error
				agent.path, err = S.graph.Path(agent.position.Name(), agent.agent.Destination())
				if err != nil {
					log.Fatal(err)
				}
			}
		}
		searchWG.Done()
	}

	// Create worker goroutines
	for i := 0; i < runtime.NumCPU(); i++ {
		go worker(searchCh)
		searchWG.Add(1)
	}

	// Send jobs to the workers
	for _, agentList := range jobs {
		agentList := agentList // Create new reference for goroutine worker
		searchCh <- agentList
	}
	close(searchCh)

	searchWG.Wait()
}

func (S *Simulation) MoveAgents() {
	for _, agent := range S.agents {
		if !S.finishedAgents[agent.agent.Id()] {
			agent.timeUntilNextChoice--
			if agent.timeUntilNextChoice <= 0 {
				// Move position and set time
				agent.position = agent.path[len(agent.path)-1].node
				agent.timeUntilNextChoice = agent.path[len(agent.path)-1].edge.Time()

				// Update agent counts on edges this node is leaving/moving on
				if agent.lastEdge != nil {
					agent.lastEdge.RemoveAgent()
				}
				agent.lastEdge = agent.path[len(agent.path)-1].edge
				agent.lastEdge.AddAgent()

				// Remember choice
				agent.history = append(agent.history, Choice{
					Edge:       agent.path[len(agent.path)-1].edge.To(),
					Timestamp:  S.currentTime,
					TravelTime: agent.timeUntilNextChoice})

				// Shorten the path
				agent.path = agent.path[:len(agent.path)-1]
			}
			if len(agent.path) == 0 {
				// Mark it as finished if it reached its destination
				S.finishedAgents[agent.agent.Id()] = true
			}
		}
	}
}

func (S *Simulation) PrintHistory(file io.Writer) {
	enc := json.NewEncoder(file)
	for _, agent := range S.agents {
		var item struct {
			Agent   string
			History []Choice
		}
		item.Agent = agent.agent.Id()
		item.History = agent.history
		enc.Encode(item)
	}
}
