package main

import (
	"log"
	"net/http"
	"os"

	"github.com/vatsal278/msgbroker/internal/router"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	r := router.Router()
	port := os.Getenv("PORT")

	log.Fatal(http.ListenAndServe(":"+port, r))

}
