package router

import (
	"github.com/faanross/orlokC2/internal/middleware"
	"log"
	"net/http"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {

	agentUUID, _ := r.Context().Value(middleware.AgentUUIDKey).(string)

	log.Printf("Endpoint %s has been hit by agent %s: \n", r.URL.Path, agentUUID)

	w.Write([]byte("I'm Mister Derp!"))

}
