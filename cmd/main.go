package main

import (
	"log"
	"net/http"

	"github.com/vatsal278/msgbroker/internal/router"
)

func main() {
	/*go func() {
		r := router.TempRouter()
		log.Fatal(http.ListenAndServe(":"+"8081", r))
	}()*/
	r := router.Router()
	log.Fatal(http.ListenAndServe(":"+"9090", r))

}

//to dos: 1. implement diff. go routines for reducing latency at various end points
//2. use maps for storing and retrieving sub and pub details and for publishing msg
//cleanup unused code
