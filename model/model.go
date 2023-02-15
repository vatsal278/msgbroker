// Package model is responsible for defining the structures or types that represent the data entities in the application.
package model

import (
	"sync"
)

// CallBack represents a callback object with the necessary fields for a subscriber to receive updates
type CallBack struct {
	// HTTP method for the callback (e.g. GET, POST)
	HttpMethod string `form:"httpMethod" json:"httpMethod" validate:"required"`
	// URL for the callback (i.e. the endpoint that the subscriber will receive updates on)
	CallbackUrl string `form:"callbackUrl" json:"callbackUrl" validate:"required"`
	// Public key for authentication (if needed)
	PublicKey string `form:"key" json:"key"`
}

// Subscriber represents a subscriber object with a callback and channel
// Subscriber is a client that subscribes to a channel to receive updates via its CallbackUrl endpoint.
type Subscriber struct {
	// CallBack field contains the details for the callback, including the HTTP method, callback URL, and public key for authentication (if needed).
	CallBack CallBack `form:"callback" json:"callback" validate:"required"`
	// Name of the channel to subscribe to (i.e. the topic that the subscriber is interested in)
	Channel string `form:"channel" json:"channel" validate:"required"`
}

// Publisher represents a publisher object that contains a unique identifier and the name of the channel it publishes to.
type Publisher struct {
	// Unique identifier for the publisher
	Id string `form:"id" json:"id"`
	// Name of the channel to publish updates to
	Channel string `form:"channel" json:"channel" validate:"required"`
}

// Updates represents an update object, containing information about the publisher
// and the message to be published.
type Updates struct {
	Publisher Publisher `form:"publisher" json:"publisher" validate:"required"` // Publisher object containing information about the publisher
	Update    string    `form:"update" json:"update" validate:"required"`       // The message to be published
}

// MessageBroker is the main object which holds the subscribers and publishers
type MessageBroker struct {
	// PubM holds the mapping of publishers for each channel
	// The key of the outer map is the channel name
	// The key of the inner map is the publisher name
	PubM map[string]map[string]struct{}

	// SubM holds the mapping of subscribers for each channel
	// The key is the channel name
	// The value is an array of Subscriber objects that subscribed to the channel
	SubM map[string][]Subscriber

	// Mutex is used to synchronize access to the message broker object
	sync.Mutex
}

// Response is a struct for responses with a status code, message, and data
type Response struct {
	Status  int         `json:"status"`  // HTTP status code
	Message string      `json:"message"` // Response message
	Data    interface{} `json:"data"`    // Response data
}

// Update is a method on the Response struct for updating its fields
func (r *Response) Update(status int, msg string, data interface{}) {
	r.Status = status
	r.Message = msg
	r.Data = data
}
