package model

import (
	"sync"
)

// CallBack represents a callback object with the necessary fields
type CallBack struct {
	HttpMethod  string `form:"httpMethod" json:"httpMethod" validate:"required"`
	CallbackUrl string `form:"callbackUrl" json:"callbackUrl" validate:"required"`
	PublicKey   string `form:"key" json:"key"`
}

// Subscriber represents a subscriber object with a callback and channel
type Subscriber struct {
	CallBack CallBack `form:"callback" json:"callback" validate:"required"`
	Channel  string   `form:"channel" json:"channel" validate:"required"`
}

// Publisher represents a publisher object with a unique id and channel
type Publisher struct {
	Id      string `form:"id" json:"id"`
	Channel string `form:"channel" json:"channel" validate:"required"`
}

// Updates represents an update object with a publisher and update string
type Updates struct {
	Publisher Publisher `form:"publisher" json:"publisher" validate:"required"`
	Update    string    `form:"update" json:"update" validate:"required"`
}

// MessageBroker is the main object which holds the subscribers and publishers
type MessageBroker struct {
	PubM map[string]map[string]struct{} // 1st map's key is channel, 2nd map's key is publisher name
	SubM map[string][]Subscriber
	sync.Mutex
}

// Response is a struct for responses with a status code, message, and data
type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Update is a method on the Response struct for updating its fields
func (r *Response) Update(status int, msg string, data interface{}) {
	r.Status = status
	r.Message = msg
	r.Data = data
}
