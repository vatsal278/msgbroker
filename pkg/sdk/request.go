// Package sdk provides a suite of wrapped functions to make the direct calls to msgbroker service
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

// msgBrokerSvc struct represents the message broker service with a given URL.
type msgBrokerSvc struct {
	msgbrokerUrl string
}

// NewMsgBrokerSvc creates and returns a new instance of MsgBrokerSvcI with the given message broker service URL.
func NewMsgBrokerSvc(url string) MsgBrokerSvcI {
	return &msgBrokerSvc{
		msgbrokerUrl: url,
	}
}

// MsgBrokerSvcI interface provides a combined interface for SubscriberSvc and PublisherSvc.
type MsgBrokerSvcI interface {
	SubscriberSvc
	PublisherSvc
}

// ExtractMsgSvc interface provides a method to extract messages using the given private key.
type ExtractMsgSvc interface {
	ExtractMsg(key *rsa.PrivateKey) func(io.ReadCloser) (string, error)
}

// SubscriberSvc interface provides methods to register a subscriber and extract messages.
type SubscriberSvc interface {
	RegisterSub(string, string, string, string) error
	ExtractMsgSvc
}

// PushSvc interface provides a method to push messages to the message broker.
type PushSvc interface {
	PushMsg(string, string, string) error
}

// PublisherSvc interface provides methods to register a publisher and push messages.
type PublisherSvc interface {
	RegisterPub(string) (string, error)
	PushSvc
}

// RegisterSub registers a subscriber with the message broker service.
func (m *msgBrokerSvc) RegisterSub(method string, callbackUrl string, publicKey string, channel string) error {
	sub := model.Subscriber{
		CallBack: model.CallBack{
			HttpMethod:  method,
			CallbackUrl: callbackUrl,
			PublicKey:   publicKey,
		},
		Channel: channel,
	}
	reqBody, err := json.Marshal(sub)
	if err != nil {
		return err
	}
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	r, err := client.Post(m.msgbrokerUrl+"/register/subscriber", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	if r.StatusCode < 200 || r.StatusCode > 299 {
		return fmt.Errorf("non success status code received : %v", r.StatusCode)
	}
	return nil
}

// RegisterPub registers a publisher with the message broker service using the provided channel. It returns the ID of the registered publisher and an error if any occurred.
func (m *msgBrokerSvc) RegisterPub(channel string) (string, error) {
	// Create a new publisher object with the given channel
	pub := model.Publisher{Channel: channel}
	// Convert the publisher object to a JSON byte array
	reqBody, err := json.Marshal(pub)
	if err != nil {
		return "", err
	}
	// Create a new HTTP client with a 2 second timeout
	client := http.Client{
		Timeout: 2 * time.Second,
	}
	// Send a POST request to the message broker's /register/publisher endpoint with the JSON request body
	r, err := client.Post(m.msgbrokerUrl+"/register/publisher", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}
	// If the HTTP response status code is not in the 200-299 range, return an error
	if r.StatusCode < 200 || r.StatusCode > 299 {
		return "", fmt.Errorf("non success status code received : %v", r.StatusCode)
	}
	// Read the response body into a byte array
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	// Unmarshal the JSON response body into a Response object
	var response model.Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}
	// Extract the data field from the Response object and convert it to a map[string]interface{}
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("unexpected response")
	}
	// Extract the ID field from the map and convert it to a string
	id, ok := data["id"].(string)
	if !ok || len(id) == 0 {
		return "", fmt.Errorf("id not found")
	}
	// Return the ID of the newly registered publisher
	return id, nil
}

// PushMsg sends a message to the message broker service
func (m *msgBrokerSvc) PushMsg(msg string, uuid string, channel string) error {
	// create a new message update with the message, publisher ID and channel
	var update = model.Updates{
		Update: msg,
		Publisher: model.Publisher{
			Id:      uuid,
			Channel: channel,
		},
	}
	// marshal the update to JSON
	reqBody, err := json.Marshal(update)
	if err != nil {
		return err
	}
	// send a POST request to the message broker service's /publish endpoint with the update in the request body
	client := http.Client{
		Timeout: 2 * time.Second,
	}
	r, err := client.Post(m.msgbrokerUrl+"/publish", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	// check that the response status code indicates success
	if r.StatusCode < 200 || r.StatusCode > 299 {
		return fmt.Errorf("non success status code received : %v", r.StatusCode)
	}
	return nil
}

// ExtractMsg takes in a pointer to an RSA private key as a parameter and returns a function that accepts a ReadCloser and returns a string and an error.
// This returned function reads the message from the ReadCloser and decrypts it using the provided private key if it's not nil, or returns the message as is if the private key is nil.
func (m *msgBrokerSvc) ExtractMsg(key *rsa.PrivateKey) func(closer io.ReadCloser) (string, error) {
	// return a function that takes a ReadCloser and returns the decrypted message or an error
	return func(source io.ReadCloser) (string, error) {
		if source == nil {
			return "", fmt.Errorf("source cannot be nil")
		}
		// read the message from the source
		body, err := ioutil.ReadAll(source)
		if err != nil {
			return "", err
		}
		// if no private key is provided, return the message as is
		if key == nil {
			return string(body), nil
		}
		// otherwise, decrypt the message using the private key
		res, err := crypt.RsaOaepDecrypt(string(body), *key)
		if err != nil {
			return "", err
		}
		return res, nil
	}
}
