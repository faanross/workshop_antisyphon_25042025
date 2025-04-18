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

func ResultHandler(w http.ResponseWriter, r *http.Request) {

	agentUUID, _ := r.Context().Value(middleware.AgentUUIDKey).(string)

	log.Printf("Endpoint %s has been hit by agent %s: \n", r.URL.Path, agentUUID)

	var result websocket.Message

	err := json.NewDecoder(r.Body).Decode(&result)

	if err != nil {
		http.Error(w, "Request body invalid", http.StatusBadRequest)
	}

	log.Printf("Result from agent: %s\n Command: %s\n Output: %s\n", agentUUID, result.Command, result.Output)

	// SEND TO CLIENT
	if websocket.GlobalWSServer != nil {
		result.Type = websocket.ResponseMessage

		result.AgentUUID = agentUUID

		websocket.GlobalWSServer.Broadcast(result)

		log.Printf("Sending result from agent %s to ALL clients\n", agentUUID)
	}

	w.WriteHeader(http.StatusOK)

}
