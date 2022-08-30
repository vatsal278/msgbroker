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
	fmt.Println("Connected to port " + os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), r))
}
