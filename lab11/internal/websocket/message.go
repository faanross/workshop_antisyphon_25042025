package websocket

type MessageType string

const (
	CommandMessage  MessageType = "command"
	ResponseMessage MessageType = "response"
)

type Message struct {
	Type      MessageType `json:"type"`
	Command   string      `json:"command,omitempty"`
	Output    string      `json:"output,omitempty"`
	AgentUUID string      `json:"agentUUID,omitempty"`
}
