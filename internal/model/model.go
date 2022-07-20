package model

type Subscriber struct {
	Name    string `form:"name" json:"name" validate:"required"`
	Channel string `form:"channel" json:"channel" validate:"required"`
}

type Publisher struct {
	Name    string `form:"name" json:"name" validate:"required"`
	Channel string `form:"channel" json:"channel" validate:"required"`
}

type Updates struct {
	Publisher string `form:"publisher" json:"publisher" validate:"required"`
	Update    string `form:"update" json:"update" validate:"required"`
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
