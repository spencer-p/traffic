/*
Package traffic implements an abstract simulation of traffic. A simulation is
comprised of a collection of Edges and Agents. A graph is constructed from the
Edges, which the Agents travel on over time. The traffic library separates
complex path finding and meticulous agent management from the actual edge and
agent implementations. This, and the interfaces for Agent and Edge, allow great
creative freedom and quick development times for any simple traffic simulation.
*/
package traffic

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"runtime"
)

// Simulation is the entry point for the traffic library. A Simulation is
// comprised of a set of Agents and a set of Edges. The Simulation manages all
// this data and manipulates it in the main event loop, Simulate.
type Simulation struct {
	graph          *graph
	agents         []*metaAgent
	agentsByGroup  map[string][]*metaAgent
	currentTime    int
	pathTimeout    int
	finishedAgents map[string]bool
}

// NewSimulation returns a new empty Simulation with a path timeout of 15
// minutes.
func NewSimulation() Simulation {
	// Path timeout is 15 minutes by default
	// TODO add a setter for this
	return Simulation{
		graph:          NewGraph(),
		finishedAgents: make(map[string]bool),
		pathTimeout:    15,
		agentsByGroup:  make(map[string][]*metaAgent),
	}
}

// AddEdge adds and processes the provided Edge to the Simulation.
func (S *Simulation) AddEdge(e Edge) {
	S.graph.AddEdge(e)
}

// AddAgent adds and processes the provided Agent to the Simulation.
func (S *Simulation) AddAgent(a Agent) {
	var ma metaAgent
	ma.agent = a
	ma.pathLastUpdated = -S.pathTimeout
	ma.position = S.graph.nodes[a.Start()] // TODO check Start exists
	ma.timeUntilNextChoice = a.LeaveTime()
	S.agents = append(S.agents, &ma)
	S.agentsByGroup[a.Group()] = append(S.agentsByGroup[a.Group()], &ma)
}

// Simulate is the main event loop for the Simulation. It moves agents along
// optimal paths, updating those paths along the way, until all the agents have
// reached their destination.
func (S *Simulation) Simulate() error {
	for len(S.finishedAgents) < len(S.agents) {
		if err := S.Tick(); err != nil {
			return err
		}
	}
	return nil
}

// Tick processes a single minute of simulation time. It refreshes necessary
// paths and then moves agents. For basic usage, use Simulate.
func (S *Simulation) Tick() error {
	if err := S.RefreshPaths(); err != nil {
		return err
	}
	S.MoveAgents()
	S.currentTime++
	return nil
}

// RefreshPaths updates the paths of any agent that needs their paths updated.
// To save time, it runs Dijkstra's algorithm on the provided edges concurrently
// on every core. For basic usage, see Simulate.
func (S *Simulation) RefreshPaths() error {
	// Error/Done channel - for marking workers as errored or done
	errorDoneCh := make(chan error)

	// Channel to execute searches and update positions of the agents
	searchCh := make(chan *metaAgent)

	// Worker function
	worker := func(jobCh <-chan *metaAgent, errorCh chan<- error) {
		for agent := range jobCh {
			// Perform the shortest spanning tree search
			tree, err := S.graph.Dijkstra(agent.position.Name(), agent.agent.Destination(), agent.agent)
			if err != nil {
				errorCh <- err
				return
			}

			// Set the agent's path
			agent.path, err = tree.Path()
			if err != nil {
				errorCh <- err
				return
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
		for _, agent := range S.agents {
			if !S.finishedAgents[agent.agent.Id()] &&
				S.currentTime-agent.pathLastUpdated >= S.pathTimeout {
				agent := agent // New ref for channel
				searchCh <- agent
			}
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

// MoveAgents moves agents to their next edge if they have completed their last
// edge travel. It additionally saves metadata for output on the way. For basic
// usage, see Simulate.
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

// PrintGroupHistory encodes the full history of group to an io.Writer. It is
// formatted as Choice history, then AverageData for each edge type, and then
// a total cumulative AverageData.
func (S *Simulation) PrintGroupHistory(file io.Writer, group string) {
	fmt.Fprintf(file, "======\nGroup: %s\n======\n", group)
	// First get group
	agents := S.agentsByGroup[group]

	enc := json.NewEncoder(file)

	// Metadata
	var metadata AverageData
	metadataByEdgeType := make(map[string]*AverageData)
	metadata.AgentCount = len(agents)

	// Print history of each agent
	for _, agent := range agents {
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
		enc.Encode(item)

		// Update total metadata
		metadata.MinutesTraveled += item.DeltaTime

		// Update breakdown metadata
		agentRecorded := make(map[string]bool)
		for _, choice := range item.History {
			// Get the avg data, create if doesn't exist
			md, ok := metadataByEdgeType[choice.EdgeType]
			if !ok {
				md = &AverageData{}
				metadataByEdgeType[choice.EdgeType] = md
			}

			// Update it
			if !agentRecorded[choice.EdgeType] {
				metadataByEdgeType[choice.EdgeType].AgentCount++
				agentRecorded[choice.EdgeType] = true
			}
			metadataByEdgeType[choice.EdgeType].AgentCountPerEdge++
			metadataByEdgeType[choice.EdgeType].MinutesTraveled += choice.TravelTime
		}
	}

	// Complete metadata and print it
	metadata.UpdateAverage()
	for _, md := range metadataByEdgeType {
		md.UpdateAverage()
	}
	enc.Encode(metadataByEdgeType)
	enc.Encode(metadata)
}

// PrintHistory calls PrintGroupHistory on every group in the Simulation.
func (S *Simulation) PrintHistory(file io.Writer) {
	for group := range S.agentsByGroup {
		S.PrintGroupHistory(file, group)
	}
}
