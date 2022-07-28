package controllerInterface

import "net/http"

type IController interface {
	RegisterPublisher() func(w http.ResponseWriter, r *http.Request)
	RegisterSubscriber() func(w http.ResponseWriter, r *http.Request)
	PublishMessage() func(w http.ResponseWriter, r *http.Request)
	NoRouteFound() func(w http.ResponseWriter, r *http.Request)
}
