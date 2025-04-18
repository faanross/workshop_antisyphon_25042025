package websocket

import (
	"log"
	"sync"
)

type CommandQueue struct {
	PendingCommands []string
	mu              sync.Mutex
}

var AgentCommands = CommandQueue{
	PendingCommands: make([]string, 0),
}

func (cq *CommandQueue) QueueCommand(command string) {
	cq.mu.Lock()
	defer cq.mu.Unlock()

	cq.PendingCommands = append(cq.PendingCommands, command)

	log.Printf("WebSocket command queue: %s", command)
}

func (cq *CommandQueue) GetCommand() (string, bool) {
	cq.mu.Lock()
	defer cq.mu.Unlock()

	if len(cq.PendingCommands) == 0 {
		return "", false
	}

	cmd := cq.PendingCommands[0]

	cq.PendingCommands = cq.PendingCommands[1:]

	log.Printf("Command retrieved: %s\n", cmd)

	return cmd, true
}
