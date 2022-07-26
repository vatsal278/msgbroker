package router

import (
	"fmt"
	"net/http"

	article_controller "github.com/vatsal278/msgbroker/internal/handler/controller"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/register/subscriber", article_controller.RegisterSubscriber()).Methods(http.MethodPost) //Endpoint for inserting
	router.HandleFunc("/register/publisher", article_controller.RegisterPublisher()).Methods(http.MethodPost)
	router.HandleFunc("/publish", article_controller.PublishMessage()).Methods(http.MethodPost)
	http.Handle("/", router)
	fmt.Println("Connected to port " + "8080")

	return router
}
