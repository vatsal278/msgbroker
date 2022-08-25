package requests

import (
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"github.com/vatsal278/msgbroker/internal/model"
	"github.com/vatsal278/msgbroker/pkg/crypt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type msgBrokerUrl struct {
	msgbrokerUrl string
}

func NewController(url string) MsgBrokerCalls {
	return &msgBrokerUrl{
		msgbrokerUrl: url,
	}
}

type MsgBrokerCalls interface {
	UpdateMsgCalls
	RegisterCalls
	ReceiveMsgCalls
}
type UpdateMsgCalls interface {
	UpdateSubs(string, string, string) (*http.Response, error)
}
type RegisterCalls interface {
	RegisterSub(string, string, string, string) (*http.Response, error)
	RegisterPub(string) (string, *http.Response, error)
}
type ReceiveMsgCalls interface {
	ReceiveMsg(closer io.ReadCloser, key *rsa.PrivateKey) func() (string, error)
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

func (m *msgBrokerUrl) RegisterSub(method string, callbackUrl string, publicKey string, channel string) (*http.Response, error) {
	sub := Subscriber{CallBack: CallBack{HttpMethod: method, CallbackUrl: callbackUrl, PublicKey: publicKey},
		Channel: channel,
	}
	reqBody, err := json.Marshal(sub)
	if err != nil {
		log.Print(err.Error())
		return nil, err
	}
	client := http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	request, err := http.NewRequest("POST", m.msgbrokerUrl+"/register/subscriber", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	log.Printf("%+v \n", *request)
	r, err := client.Do(request)
	//_, err = client.Post(m.msgbrokerUrl+"/register/subscriber", "application/json", bytes.NewBuffer(reqBody))
	//log.Println(r.Status)
	if err != nil {
		log.Print(err.Error())
		return nil, err
	}
	log.Print(r.Status)

	return r, nil
}

func (m *msgBrokerUrl) RegisterPub(channel string) (string, *http.Response, error) {
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
		return "", nil, err
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print(err.Error())
		return "", r, err
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Print(err.Error())
		return "", r, err
	}
	data := response.Data.(map[string]interface{})
	id := data["id"]
	return id.(string), r, nil

}

func (m *msgBrokerUrl) UpdateSubs(msg string, key string, channel string) (*http.Response, error) {
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
		return nil, err
	}
	client := http.Client{
		Timeout: time.Duration(2 * time.Second),
	}
	r, err := client.Post(m.msgbrokerUrl+"/publish", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Print(err.Error())
		return nil, err
	}
	log.Print(r.Header)
	return r, nil
}
func (m *msgBrokerUrl) ReceiveMsg(closer io.ReadCloser, key *rsa.PrivateKey) func() (string, error) {
	return func() (string, error) {
		body, err := ioutil.ReadAll(closer)
		if err != nil {
			log.Print(err)
			return "", err
		}
		defer closer.Close()
		res, err := crypt.RsaOaepDecrypt(string(body), *key)
		if err != nil {
			log.Print(err)
			return "", err
		}
		log.Printf("%+v", res)
		return res, nil
	}
}
