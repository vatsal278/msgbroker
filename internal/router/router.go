package router

import (
	"fmt"
	"net/http"

	"github.com/vatsal278/msgbroker/internal/handler/controller"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	i := controller.NewController()

	router := mux.NewRouter()
	router.HandleFunc("/register/subscriber", i.RegisterSubscriber()).Methods(http.MethodPost) //Endpoint for inserting
	router.HandleFunc("/register/publisher", i.RegisterPublisher()).Methods(http.MethodPost)
	router.HandleFunc("/publish", i.PublishMessage()).Methods(http.MethodPost)
	http.Handle("/", router)
	fmt.Println("Connected to port " + "8080")

	return router
}
