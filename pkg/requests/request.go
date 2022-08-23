package requests

import (
	"bytes"
	"encoding/json"
	"github.com/vatsal278/msgbroker/internal/model"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type msgBrokerUrl struct {
	msgbrokerUrl string
}

func NewController(url string) ApiCalls {
	return &msgBrokerUrl{
		msgbrokerUrl: url,
	}
}

type ApiCalls interface {
	RegisterSub(string, string, string, string) error
	RegisterPub(string) (string, error)
	UpdateSubs(string, string, string) error
}
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

func (m *msgBrokerUrl) RegisterSub(method string, callbackUrl string, publicKey string, channel string) error {
	sub := Subscriber{CallBack: CallBack{HttpMethod: method, CallbackUrl: callbackUrl, PublicKey: publicKey},
		Channel: channel,
	}
	reqBody, err := json.Marshal(sub)
	if err != nil {
		log.Print(err.Error())
		return err
	}
	client := http.Client{
		Timeout: time.Duration(2 * time.Second),
	}
	r, err := client.Post(m.msgbrokerUrl+"/register/subscriber", "application/json", bytes.NewBuffer(reqBody))
	log.Println(r.Status)
	if err != nil {
		log.Print(err.Error())
		return err
	}

	return nil
}

func (m *msgBrokerUrl) RegisterPub(channel string) (string, error) {
	var response model.Response
	pub := Publisher{Channel: channel}
	reqBody, err := json.Marshal(pub)
	client := http.Client{
		Timeout: time.Duration(2 * time.Second),
	}
	//request, err := http.NewRequest("POST","http://localhost:9090/register/subscriber",bytes.NewBuffer(reqBody))
	r, err := client.Post(m.msgbrokerUrl+"/register/publisher", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Print(err.Error())
		return "", err
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print(err.Error())
		return "", err
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Print(err.Error())
		return "", err
	}
	data := response.Data.(map[string]interface{})
	id := data["id"]
	return id.(string), nil

}

func (m *msgBrokerUrl) UpdateSubs(msg string, key string, channel string) error {
	var update = Updates{
		Update: msg,
		Publisher: Publisher{
			Id:      key,
			Channel: channel,
		},
	}
	reqBody, err := json.Marshal(update)
	if err != nil {
		log.Print(err.Error())
		return err
	}
	client := http.Client{
		Timeout: time.Duration(2 * time.Second),
	}
	r, err := client.Post(m.msgbrokerUrl+"/publish", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Print(err.Error())
		return err
	}
	log.Print(r.Header)
	return nil
}
