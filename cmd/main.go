package main

import (
	"log"
	"net/http"

	"github.com/vatsal278/msgbroker/internal/router"
)

func main() {
	r := router.Router()
	log.Fatal(http.ListenAndServe(":"+"9090", r))

}
