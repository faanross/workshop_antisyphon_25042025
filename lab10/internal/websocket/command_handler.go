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
