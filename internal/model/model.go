package model

import "sync"

type Subscriber struct {
	Subscriber struct {
		HttpMethod  string
		CallbackUrl string
	}
	Channel string `form:"channel" json:"channel" validate:"required"`
}

type Publisher struct {
	Name    map[string]string `form:"name" json:"name" validate:"required"`
	Channel string            `form:"channel" json:"channel" validate:"required"`
}

type Updates struct {
	Publisher string `form:"publisher" json:"publisher" validate:"required"`
	Update    string `form:"update" json:"update" validate:"required"`
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

type Temp_struct struct {
	Id      interface{} `json:"status"`
	Title   interface{} `json:"message"`
	Content interface{}
	Author  interface{}
}
