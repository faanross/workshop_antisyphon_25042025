package router

import (
	"log"
	"net/http"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {

	log.Printf("Endpoint has been hit: %s\n", r.URL.Path)

	w.Write([]byte("I'm Mister Derp!"))

}
