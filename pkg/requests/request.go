package requests

import (
	"bytes"
	"encoding/json"
	"github.com/vatsal278/msgbroker/internal/model"
	"io/ioutil"
	"log"
	"net/http"
)

type CallBack struct {
	HttpMethod  string `form:"httpMethod" json:"httpMethod" validate:"required"`
	CallbackUrl string `form:"callbackUrl" json:"callbackUrl" validate:"required"`
	PublicKey   string `form:"key" json:"key"`
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

func RegisterSub(v Subscriber) error {
	x, err := json.Marshal(v)

	reqBody := []byte(x)
	_, err = http.Post("http://localhost:9090/register/subscriber", "application/json", bytes.NewBuffer(reqBody))

	return err
}

func RegisterPub(url string, v Publisher) (string, error) {
	var response model.Response
	x, err := json.Marshal(v)
	reqBody := []byte(x)
	r, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	log.Print(r.Body)
	body, err := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &response)
	data := response.Data.(map[string]interface{})
	id := data["id"]
	return id.(string), err
}

func UpdateSubs(url string, msg string, key string, channel string) error {
	var update model.Updates
	update.Update = msg
	update.Publisher.Id = key
	update.Publisher.Channel = channel
	//var response model.Response
	x, err := json.Marshal(update)
	reqBody := []byte(x)
	r, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	log.Print(r.Body)

	return err
}
