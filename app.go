package rails

import (
	"log"
	"net/http"
	"os"
)

var (
	router *Router
)

func SetRouter(r *Router) {
	router = r
}

func Run() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Printf("Listening on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
