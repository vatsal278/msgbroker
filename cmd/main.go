package main

import (
	"log"
	"msgbroker/internal/router"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	r := router.Router()
	port := os.Getenv("PORT")

	log.Fatal(http.ListenAndServe(":"+port, r))

}
