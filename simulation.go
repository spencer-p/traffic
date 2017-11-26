package traffic

import (
	"encoding/json"
	"io"
	"reflect"
	"runtime"
)

type Simulation struct {
	graph          *Graph
	agents         []*metaAgent
	currentTime    int
	pathTimeout    int
	finishedAgents map[string]bool
}

func NewSimulation() Simulation {
	// Path timeout is 15 minutes by default
	// TODO add a setter for this
	return Simulation{
		graph:          NewGraph(),
		finishedAgents: make(map[string]bool),
		pathTimeout:    15}
}

func (S *Simulation) AddEdge(e Edge) {
	S.graph.AddEdge(e)
}

func (S *Simulation) AddAgent(a Agent) {
	var ma metaAgent
	ma.agent = a
	ma.pathLastUpdated = -S.pathTimeout
	ma.position = S.graph.nodes[a.Start()] // TODO check Start exists
	ma.timeUntilNextChoice = a.LeaveTime()
	S.agents = append(S.agents, &ma)
}

func (S *Simulation) Simulate() error {
	for len(S.finishedAgents) < len(S.agents) {
		if err := S.Tick(); err != nil {
			return err
		}
	}
	return nil
}

func (S *Simulation) Tick() error {
	if err := S.RefreshPaths(); err != nil {
		return err
	}
	S.MoveAgents()
	S.currentTime++
	return nil
}

func (S *Simulation) RefreshPaths() error {
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

	// Error/Done channel - for marking workers as errored or done
	errorDoneCh := make(chan error)

	// Channel to execute searches and update positions of the agents
	searchCh := make(chan []*metaAgent)

	// Worker function
	worker := func(jobCh <-chan []*metaAgent, errorCh chan<- error) {
		for agentList := range jobCh {
			// Perform the shortest spanning tree search
			_, err := S.graph.Dijkstra(agentList[0].position.Name())
			if err != nil {
				errorCh <- err
				return
			}

			// Set each agent's path
			for _, agent := range agentList {
				agent.path, err = S.graph.Path(agent.position.Name(), agent.agent.Destination())
				if err != nil {
					errorCh <- err
					return
				}
			}
		}

		// Done
		errorCh <- nil
	}

	// Create worker goroutines
	var workerCount int
	for i := 0; i < runtime.NumCPU(); i++ {
		go worker(searchCh, errorDoneCh)
		workerCount++
	}

	// Send jobs to the workers
	go func() {
		for _, agentList := range jobs {
			agentList := agentList // Create new reference for goroutine worker
			searchCh <- agentList
		}
		close(searchCh)
	}()

	// Wait for all jobs to finish, or break with error early
	for err := range errorDoneCh {
		workerCount--
		if err != nil {
			return err
		} else if workerCount == 0 {
			break
		}
	}

	// Success
	return nil
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
					To:         agent.path[len(agent.path)-1].edge.To(),
					EdgeType:   reflect.TypeOf(agent.path[len(agent.path)-1].edge).Elem().Name(),
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

	// Metadata
	var metadata struct {
		AgentCount             int
		MinutesTraveled        int
		AverageDeltaTravelTime float64
	}
	metadata.AgentCount = len(S.agents)

	// Print history of each agent
	for _, agent := range S.agents {
		var item struct {
			Agent     string
			History   []Choice
			StartTime int
			EndTime   int
			DeltaTime int
		}
		item.Agent = agent.agent.Id()
		item.History = agent.history
		item.StartTime = agent.history[0].Timestamp
		item.EndTime = agent.history[len(agent.history)-1].Timestamp +
			agent.history[len(agent.history)-1].TravelTime
		item.DeltaTime = item.EndTime - item.StartTime
		metadata.MinutesTraveled += item.DeltaTime
		enc.Encode(item)
	}

	// Complete metadata and print it
	metadata.AverageDeltaTravelTime = float64(metadata.MinutesTraveled) / float64(metadata.AgentCount)
	enc.Encode(metadata)
}
