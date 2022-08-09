package model

import (
	"sync"
)

type CallBack struct {
	HttpMethod  string `form:"httpMethod" json:"httpMethod" validate:"required"`
	CallbackUrl string `form:"callbackUrl" json:"callbackUrl" validate:"required"`
}

type Subscriber struct {
	CallBack CallBack `form:"callback" json:"callback" validate:"required"`
	Channel  string   `form:"channel" json:"channel" validate:"required"`
}

type Publisher struct {
	//Name    string `form:"name" json:"name" validate:"required"`
	Id      string `form:"id" json:"id"`
	Channel string `form:"channel" json:"channel" validate:"required"`
}

type Updates struct {
	Publisher Publisher `form:"publisher" json:"publisher" validate:"required"`
	Update    string    `form:"update" json:"update" validate:"required"`
}

type MessageBroker struct {
	PubM map[string]map[string]struct{} // 1st map's key is channel, 2nd map's key is publisher name
	SubM map[string][]Subscriber
	sync.Mutex
}
type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (r *Response) Update(status int, msg string, data interface{}) {
	r.Status = status
	r.Message = msg
	r.Data = data
}
