package ui

import (
	"github.com/agentsdance/agentx/internal/agent"
)

// AgentStatus represents the status of an agent
type AgentStatus struct {
	Agent     agent.Agent
	Installed bool
	Exists    bool
	Error     error
}

// Model is the TUI state model
type Model struct {
	agents   []AgentStatus
	cursor   int
	message  string
	quitting bool
}

// NewModel creates a new TUI model
func NewModel() Model {
	agents := agent.GetAllAgents()
	statuses := make([]AgentStatus, len(agents))

	for i, a := range agents {
		installed, err := a.HasPlaywright()
		statuses[i] = AgentStatus{
			Agent:     a,
			Installed: installed,
			Exists:    a.Exists(),
			Error:     err,
		}
	}

	return Model{
		agents: statuses,
		cursor: 0,
	}
}

// refreshStatus refreshes the status of all agents
func (m *Model) refreshStatus() {
	for i := range m.agents {
		installed, err := m.agents[i].Agent.HasPlaywright()
		m.agents[i].Installed = installed
		m.agents[i].Exists = m.agents[i].Agent.Exists()
		m.agents[i].Error = err
	}
}
