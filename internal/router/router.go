package router

import (
	"net/http"

	"github.com/vatsal278/msgbroker/internal/handler/controller"

	"github.com/gorilla/mux"
)

//Router returns a new instance of a mux.Router and sets up routes on it.
func Router() *mux.Router {
	i := controller.NewController()

	router := mux.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(i.NoRouteFound())
	router.HandleFunc("/register/subscriber", i.RegisterSubscriber()).Methods(http.MethodPost) //Endpoint for inserting
	router.HandleFunc("/register/publisher", i.RegisterPublisher()).Methods(http.MethodPost)
	router.HandleFunc("/publish", i.PublishMessage()).Methods(http.MethodPost)
	http.Handle("/", router)

	return router
}
