package router

import (
	"fmt"
	"net/http"

	article_controller "github.com/vatsal278/msgbroker/internal/handler/controller"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	var e article_controller.IController
	//Initialised the router
	router := mux.NewRouter()
	//router.NotFoundHandler = http.HandlerFunc(e.NoRouteFound())
	router.HandleFunc("/subscriber", e.RegisterSubscriber()).Methods(http.MethodPost) //Endpoint for inserting
	router.HandleFunc("/publisher", e.RegisterPublisher()).Methods(http.MethodPost)
	router.HandleFunc("/update", e.PublishMessage()).Methods(http.MethodPost)
	http.Handle("/", router)
	fmt.Println("Connected to port " + "8080")

	return router
}

func TempRouter() *mux.Router {
	var e article_controller.TController
	router := mux.NewRouter()
	router.HandleFunc("/notifications", e.NotifySubscriber()).Methods(http.MethodPost)
	http.Handle("/", router)
	fmt.Println("Connected to port " + "8081")
	return router
}
