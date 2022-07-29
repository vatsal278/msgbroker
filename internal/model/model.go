package model

import "sync"

type CallBack struct {
	HttpMethod  string `validate:"required"`
	CallbackUrl string `validate:"required"`
}

type TempPublisher struct {
	Name    interface{} `validate:"required"`
	Channel string      `form:"channel" json:"channel" validate:"required"`
}
type Subscriber struct {
	CallBack CallBack
	Channel  string `form:"channel" json:"channel" validate:"required"`
}
type TempSubscriber struct {
	CallBack CallBack
	Channel  int `form:"channel" json:"channel" validate:"required"`
}

type Publisher struct {
	Name    string `form:"name" json:"name" validate:"required"`
	Channel string `form:"channel" json:"channel" validate:"required"`
}

type Updates struct {
	Publisher Publisher `form:"publisher" json:"publisher" validate:"required"`
	Update    string    `form:"update" json:"update" validate:"required"`
}
type TempUpdates struct {
	Publisher TempPublisher `form:"publisher" json:"publisher" validate:"required"`
	Update    int           `form:"update" json:"update" validate:"required"`
}

type MessageBroker struct {
	PubM map[string]map[string]struct{} // 1st map's key is channel, 2nd map's key is publisher name
	SubM map[string][]Subscriber
	sync.Mutex
}
type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    interface{}
}
type TempResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    interface{}
}

func (r *Response) Update(status int, msg string, data interface{}) {
	r.Status = status
	r.Message = msg
	r.Data = data
}
