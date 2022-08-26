package sdk

import (
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"github.com/vatsal278/msgbroker/model"
	"github.com/vatsal278/msgbroker/pkg/crypt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type msgBrokerSvc struct {
	msgbrokerUrl string
}

func NewController(url string) MsgBrokerSvcI {
	return &msgBrokerSvc{
		msgbrokerUrl: url,
	}
}

type MsgBrokerSvcI interface {
	SubscriberSvc
	PublisherSvc
}
type ExtractMsgSvc interface {
	ExtractMsg(key *rsa.PrivateKey) func(io.ReadCloser) (string, error)
}
type SubscriberSvc interface {
	RegisterSub(string, string, string, string) error
	ExtractMsgSvc
}

type PushSvc interface {
	PushMsg(string, string, string) error
}
type PublisherSvc interface {
	RegisterPub(string) (string, error)
	PushSvc
}

func (m *msgBrokerSvc) RegisterSub(method string, callbackUrl string, publicKey string, channel string) error {
	sub := model.Subscriber{CallBack: model.CallBack{HttpMethod: method, CallbackUrl: callbackUrl, PublicKey: publicKey},
		Channel: channel,
	}
	reqBody, err := json.Marshal(sub)
	if err != nil {
		return err
	}
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	request, err := http.NewRequest("POST", m.msgbrokerUrl+"/register/subscriber", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	r, err := client.Do(request)
	if r.StatusCode < 200 || r.StatusCode > 299 {
		return fmt.Errorf("non success status code received : %v", r.StatusCode)
	}
	//_, err = client.Post(m.msgbrokerUrl+"/register/subscriber", "application/json", bytes.NewBuffer(reqBody))
	//log.Println(r.Status)
	if err != nil {
		return err
	}
	return nil
}

func (m *msgBrokerSvc) RegisterPub(channel string) (string, error) {
	var response model.Response
	pub := model.Publisher{Channel: channel}
	reqBody, err := json.Marshal(pub)
	if err != nil {
		return "", err
	}
	client := http.Client{
		Timeout: 2 * time.Second,
	}
	//request, err := http.NewRequest("POST","http://localhost:9090/register/subscriber",bytes.NewBuffer(reqBody))
	r, err := client.Post(m.msgbrokerUrl+"/register/publisher", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}
	if r.StatusCode < 200 || r.StatusCode > 299 {
		return "", fmt.Errorf("non success status code received : %v", r.StatusCode)
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("unexpected response")
	}
	id, ok := data["id"].(string)
	if !ok || len(id) == 0 {
		return "", fmt.Errorf("id not found")
	}

	return id, nil

}

func (m *msgBrokerSvc) PushMsg(msg string, uuid string, channel string) error {
	var update = model.Updates{
		Update: msg,
		Publisher: model.Publisher{
			Id:      uuid,
			Channel: channel,
		},
	}
	reqBody, err := json.Marshal(update)
	if err != nil {
		return err
	}
	client := http.Client{
		Timeout: 2 * time.Second,
	}
	r, err := client.Post(m.msgbrokerUrl+"/publish", "application/json", bytes.NewBuffer(reqBody))
	if r.StatusCode < 200 || r.StatusCode > 299 {
		return fmt.Errorf("non success status code received : %v", r.StatusCode)
	}
	if err != nil {
		return err
	}
	return nil
}
func (m *msgBrokerSvc) ExtractMsg(key *rsa.PrivateKey) func(closer io.ReadCloser) (string, error) {
	return func(source io.ReadCloser) (string, error) {
		//check if source is not nill
		body, err := ioutil.ReadAll(source)
		if err != nil {
			return "", err
		}
		if key == nil {
			return string(body), nil
		}
		res, err := crypt.RsaOaepDecrypt(string(body), *key)
		if err != nil {
			return "", err
		}
		return res, nil
	}
}
