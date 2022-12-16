package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/vatsal278/msgbroker/internal/router"
)

func main() {
	r := router.Router()
	port := os.Getenv("PORT")
	if port == "" {
		port = "9091"
	}
	fmt.Println("Connected to port " + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
