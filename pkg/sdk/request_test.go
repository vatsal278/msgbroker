package sdk

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/vatsal278/msgbroker/model"
	"github.com/vatsal278/msgbroker/pkg/crypt"
	"github.com/vatsal278/msgbroker/pkg/responseWriter"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func testServer(url string, f func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	router := mux.NewRouter()
	router.HandleFunc(url, f).Methods(http.MethodPost)
	svr := httptest.NewServer(router)
	return svr
}

func Test_RegisterPub(t *testing.T) {
	tests := []struct {
		name              string
		channel           string
		setupFunc         func() *httptest.Server
		mockServerHandler func(w http.ResponseWriter, r *http.Request)
		ValidateFunc      func(uuid string, err error)
		cleanupFunc       func(*httptest.Server)
		expectedResponse  model.Response
	}{
		{
			name: "Success:: Register Publisher",
			setupFunc: func() *httptest.Server {
				svr := testServer("/register/publisher", func(w http.ResponseWriter, r *http.Request) {
					var publisher model.Publisher
					body, err := ioutil.ReadAll(r.Body)
					if err != nil {
						t.Log(err.Error())
						t.Fail()
					}
					err = json.Unmarshal(body, &publisher)
					if err != nil {
						t.Log(err.Error())
						t.Fail()
					}
					if publisher.Channel != "channel1" {
						t.Errorf("Want %v, Got %v", "channel1", publisher.Channel)
					}
					responseWriter.ResponseWriter(w, http.StatusCreated, "", map[string]interface{}{
						"id": "b2ae109d-1382-4b1c-a8ab-5a9d04555e4e",
					}, &model.Response{})
				})
				return svr
			},
			channel: "channel1",
			ValidateFunc: func(uuid string, err error) {
				if err != nil {
					t.Errorf("Want: %v, Got: %v", nil, err.Error())
				}
				if uuid != "b2ae109d-1382-4b1c-a8ab-5a9d04555e4e" {
					t.Errorf("Want: %v, Got: %v", "b2ae109d-1382-4b1c-a8ab-5a9d04555e4e", uuid)
				}
			},
			cleanupFunc: func(svr *httptest.Server) {
				svr.Close()
			},
		},
		{
			name: "Failure:: Register Publisher:: ReadAll Failure",
			setupFunc: func() *httptest.Server {
				svr := testServer("/register/publisher", func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Length", "1")
					responseWriter.ResponseWriter(w, http.StatusCreated, "", map[string]interface{}{
						"id": "b2ae109d-1382-4b1c-a8ab-5a9d04555e4e",
					}, &model.Response{})
				})
				return svr
			},
			channel: "channel1",
			ValidateFunc: func(uuid string, err error) {
				if err.Error() != errors.New("unexpected EOF").Error() {
					t.Errorf("Want: %v, Got: %v", nil, err.Error())
				}
			},
			cleanupFunc: func(svr *httptest.Server) {
				svr.Close()
			},
		},
		{
			name: "Failure:: Register Publisher::Id not found",
			setupFunc: func() *httptest.Server {
				svr := testServer("/register/publisher", func(w http.ResponseWriter, r *http.Request) {
					var publisher model.Publisher
					body, err := ioutil.ReadAll(r.Body)
					if err != nil {
						t.Log(err.Error())
						t.Fail()
					}
					err = json.Unmarshal(body, &publisher)
					if err != nil {
						t.Log(err.Error())
						t.Fail()
					}
					if publisher.Channel != "channel1" {
						t.Errorf("Want %v, Got %v", "channel1", publisher.Channel)
					}
					responseWriter.ResponseWriter(w, http.StatusCreated, "", map[string]interface{}{
						"key": "b2ae109d-1382-4b1c-a8ab-5a9d04555e4e",
					}, &model.Response{})
				})
				return svr
			},
			channel: "channel1",
			ValidateFunc: func(uuid string, err error) {
				expectedErr := fmt.Errorf("id not found")
				if err.Error() != expectedErr.Error() {
					t.Errorf("Want: %v, Got: %v", expectedErr.Error(), err.Error())
				}
				if uuid != "" {
					t.Errorf("Want: %v, Got: %v", "", uuid)
				}
			},
			cleanupFunc: func(svr *httptest.Server) {
				svr.Close()
			},
		},
		{
			name: "Failure:: Register Publisher::Non Succes StatusCode Received from Server",
			setupFunc: func() *httptest.Server {
				svr := testServer("/register/publisher", func(w http.ResponseWriter, r *http.Request) {
					var publisher model.Publisher
					body, err := ioutil.ReadAll(r.Body)
					if err != nil {
						t.Log(err.Error())
						t.Fail()
					}
					err = json.Unmarshal(body, &publisher)
					if err != nil {
						t.Log(err.Error())
						t.Fail()
					}
					if publisher.Channel != "channel1" {
						t.Errorf("Want %v, Got %v", "channel1", publisher.Channel)
					}
					responseWriter.ResponseWriter(w, http.StatusConflict, "", nil, &model.Response{})
				})
				return svr
			},
			channel: "channel1",
			ValidateFunc: func(uuid string, err error) {
				expectedErr := fmt.Errorf("non success status code received : %v", http.StatusConflict)
				if err.Error() != expectedErr.Error() {
					t.Errorf("Want: %v, Got: %v", expectedErr.Error(), err.Error())
				}
				if uuid != "" {
					t.Errorf("Want: %v, Got: %v", "", uuid)
				}
			},
			cleanupFunc: func(svr *httptest.Server) {
				svr.Close()
			},
		},
		{
			name: "Failure:: Register Publisher::Unexpected Response",
			setupFunc: func() *httptest.Server {
				svr := testServer("/register/publisher", func(w http.ResponseWriter, r *http.Request) {
					var publisher model.Publisher
					body, err := ioutil.ReadAll(r.Body)
					if err != nil {
						t.Log(err.Error())
						t.Fail()
					}
					err = json.Unmarshal(body, &publisher)
					if err != nil {
						t.Log(err.Error())
						t.Fail()
					}
					if publisher.Channel != "channel1" {
						t.Errorf("Want %v, Got %v", "channel1", publisher.Channel)
					}
					responseWriter.ResponseWriter(w, http.StatusCreated, "", "Hello World", &model.Response{})
				})
				return svr
			},
			channel: "channel1",
			ValidateFunc: func(uuid string, err error) {
				expectedErr := fmt.Errorf("unexpected response")
				if err.Error() != expectedErr.Error() {
					t.Errorf("Want: %v, Got: %v", expectedErr.Error(), err.Error())
				}
				if uuid != "" {
					t.Errorf("Want: %v, Got: %v", "", uuid)
				}
			},
			cleanupFunc: func(svr *httptest.Server) {
				svr.Close()
			},
		},
		{
			name: "Failure:: Register Publisher::Unmarshalling Error",
			setupFunc: func() *httptest.Server {
				svr := testServer("/register/publisher", func(w http.ResponseWriter, r *http.Request) {
					var publisher model.Publisher
					body, err := ioutil.ReadAll(r.Body)
					if err != nil {
						t.Log(err.Error())
						t.Fail()
					}
					err = json.Unmarshal(body, &publisher)
					if err != nil {
						t.Log(err.Error())
						t.Fail()
					}
					if publisher.Channel != "channel1" {
						t.Errorf("Want %v, Got %v", "channel1", publisher.Channel)
					}
				})
				return svr
			},
			channel: "channel1",
			ValidateFunc: func(uuid string, err error) {
				expectedErr := fmt.Errorf("unexpected end of JSON input")
				if err.Error() != expectedErr.Error() {
					t.Errorf("Want: %v, Got: %v", expectedErr.Error(), err.Error())
				}
				if uuid != "" {
					t.Errorf("Want: %v, Got: %v", "", uuid)
				}
			},
			cleanupFunc: func(svr *httptest.Server) {
				svr.Close()
			},
		},
		{
			name: "Failure:: Register Publisher::Http Call Fail",
			setupFunc: func() *httptest.Server {
				svr := testServer("", func(w http.ResponseWriter, r *http.Request) {})
				svr.Close()
				return svr
			},
			channel: "channel1",
			ValidateFunc: func(uuid string, err error) {
				if err == nil {
					t.Log(err)
					t.Errorf("Want: %v, Got: %v", "not nil", nil)
				}
				if uuid != "" {
					t.Errorf("Want: %v, Got: %v", "", uuid)
				}
			},
			cleanupFunc: func(svr *httptest.Server) {
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svr := tt.setupFunc()
			defer tt.cleanupFunc(svr)

			calls := NewMsgBrokerSvc(svr.URL)
			uuid, err := calls.RegisterPub(tt.channel)

			tt.ValidateFunc(uuid, err)
		})
	}
}

func Test_RegisterSub(t *testing.T) {
	type args struct {
		httpMethod string
		callBack   string
		publicKey  string
		channel    string
	}
	tests := []struct {
		name              string
		reqBody           args
		setupFunc         func(args) *httptest.Server
		mockServerHandler func(w http.ResponseWriter, r *http.Request)
		ValidateFunc      func(err error)
		cleanupFunc       func(*httptest.Server)
		expectedResponse  model.Response
	}{
		{
			name: "Success:: Register Subscriber",
			setupFunc: func(a args) *httptest.Server {
				svr := testServer("/register/subscriber", func(w http.ResponseWriter, r *http.Request) {
					var subscriber model.Subscriber
					body, err := ioutil.ReadAll(r.Body)
					if err != nil {
						t.Log(err.Error())
						t.Fail()
					}
					err = json.Unmarshal(body, &subscriber)
					if err != nil {
						t.Log(err.Error())
						t.Fail()
					}
					if subscriber.Channel != "channel1" {
						t.Errorf("Want %v, Got %v", "channel1", subscriber.Channel)
					}
					x := model.Subscriber{
						CallBack: model.CallBack{
							HttpMethod:  a.httpMethod,
							CallbackUrl: a.callBack,
							PublicKey:   a.publicKey,
						},
						Channel: a.channel,
					}
					if !reflect.DeepEqual(x, subscriber) {
						t.Errorf("Want %v, Got %v", x, subscriber)
					}
					responseWriter.ResponseWriter(w, http.StatusCreated, "", nil, &model.Response{})
				})
				return svr
			},
			reqBody: args{
				httpMethod: "GET",
				callBack:   "http://localhost:8086/pong",
				publicKey:  "Hello World",
				channel:    "channel1",
			},
			ValidateFunc: func(err error) {
				if err != nil {
					t.Errorf("Want: %v, Got: %v", nil, err.Error())
				}
			},
			cleanupFunc: func(svr *httptest.Server) {
				svr.Close()
			},
		},
		{
			name: "Failure:: Register Subscriber::Non Success StatusCode Received from Server",
			setupFunc: func(a args) *httptest.Server {
				svr := testServer("/register/subscriber", func(w http.ResponseWriter, r *http.Request) {
					var subscriber model.Subscriber
					body, err := ioutil.ReadAll(r.Body)
					if err != nil {
						t.Log(err.Error())
						t.Fail()
					}
					err = json.Unmarshal(body, &subscriber)
					if err != nil {
						t.Log(err.Error())
						t.Fail()
					}
					x := model.Subscriber{
						CallBack: model.CallBack{
							HttpMethod:  a.httpMethod,
							CallbackUrl: a.callBack,
							PublicKey:   a.publicKey,
						},
						Channel: a.channel,
					}
					if !reflect.DeepEqual(x, subscriber) {
						t.Errorf("Want %v, Got %v", x, subscriber)
					}
					responseWriter.ResponseWriter(w, http.StatusConflict, "", nil, &model.Response{})
				})
				return svr
			},
			reqBody: args{
				httpMethod: "GET",
				callBack:   "http://localhost:8086/pong",
				publicKey:  "Hello World",
				channel:    "channel1",
			},
			ValidateFunc: func(err error) {
				expectedErr := fmt.Errorf("non success status code received : %v", http.StatusConflict)
				if err.Error() != expectedErr.Error() {
					t.Errorf("Want: %v, Got: %v", expectedErr.Error(), err.Error())
				}
			},
			cleanupFunc: func(svr *httptest.Server) {
				svr.Close()
			},
		},
		{
			name: "Failure:: Register Subscriber::Http Call Fail",
			setupFunc: func(a args) *httptest.Server {
				svr := testServer("", func(w http.ResponseWriter, r *http.Request) {})
				svr.Close()
				return svr
			},
			reqBody: args{
				httpMethod: "GET",
				callBack:   "http://localhost:8086/pong",
				publicKey:  "Hello World",
				channel:    "channel1",
			},
			ValidateFunc: func(err error) {
				t.Log(err)
				if err == nil {
					t.Errorf("Want: %v, Got: %v", "not nil", nil)
				}
			},
			cleanupFunc: func(svr *httptest.Server) {
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			svr := tt.setupFunc(tt.reqBody)
			defer tt.cleanupFunc(svr)

			calls := NewMsgBrokerSvc(svr.URL)
			err := calls.RegisterSub(tt.reqBody.httpMethod, tt.reqBody.callBack, tt.reqBody.publicKey, tt.reqBody.channel)

			tt.ValidateFunc(err)
		})
	}
}

func Test_PushMsg(t *testing.T) {
	type args struct {
		msg     string
		uuid    string
		channel string
	}
	tests := []struct {
		name              string
		reqBody           args
		setupFunc         func(args) *httptest.Server
		mockServerHandler func(w http.ResponseWriter, r *http.Request)
		ValidateFunc      func(err error)
		cleanupFunc       func(*httptest.Server)
		expectedResponse  model.Response
	}{
		{
			name: "Success:: PushMsg",
			setupFunc: func(a args) *httptest.Server {
				svr := testServer("/publish", func(w http.ResponseWriter, r *http.Request) {
					var update model.Updates
					body, err := ioutil.ReadAll(r.Body)
					if err != nil {
						t.Log(err.Error())
						t.Fail()
					}
					err = json.Unmarshal(body, &update)
					if err != nil {
						t.Log(err.Error())
						t.Fail()
					}

					x := model.Updates{
						Publisher: model.Publisher{
							Id:      a.uuid,
							Channel: a.channel,
						},
						Update: a.msg,
					}
					if !reflect.DeepEqual(x, update) {
						t.Errorf("Want %v, Got %v", x, update)
					}
					responseWriter.ResponseWriter(w, http.StatusOK, "", nil, &model.Response{})
				})
				return svr
			},
			reqBody: args{
				msg:     "Hello World",
				uuid:    "b2ae109d-1382-4b1c-a8ab-5a9d04555e4e",
				channel: "channel1",
			},
			ValidateFunc: func(err error) {
				if err != nil {
					t.Errorf("Want: %v, Got: %v", nil, err.Error())
				}
			},
			cleanupFunc: func(svr *httptest.Server) {
				svr.Close()
			},
		},
		{
			name: "Failure:: PushMsg ::Non Success StatusCode Received from Server",
			setupFunc: func(a args) *httptest.Server {
				svr := testServer("/publish", func(w http.ResponseWriter, r *http.Request) {
					var update model.Updates
					body, err := ioutil.ReadAll(r.Body)
					if err != nil {
						t.Log(err.Error())
						t.Fail()
					}
					err = json.Unmarshal(body, &update)
					if err != nil {
						t.Log(err.Error())
						t.Fail()
					}

					x := model.Updates{
						Publisher: model.Publisher{
							Id:      a.uuid,
							Channel: a.channel,
						},
						Update: a.msg,
					}
					if !reflect.DeepEqual(x, update) {
						t.Errorf("Want %v, Got %v", x, update)
					}
					responseWriter.ResponseWriter(w, http.StatusConflict, "", nil, &model.Response{})
				})
				return svr
			},
			reqBody: args{
				msg:     "Hello World",
				uuid:    "b2ae109d-1382-4b1c-a8ab-5a9d04555e4e",
				channel: "channel1",
			},
			ValidateFunc: func(err error) {
				expectedErr := fmt.Errorf("non success status code received : %v", http.StatusConflict)
				if err.Error() != expectedErr.Error() {
					t.Errorf("Want: %v, Got: %v", expectedErr.Error(), err.Error())
				}
			},
			cleanupFunc: func(svr *httptest.Server) {
				svr.Close()
			},
		},

		{
			name: "Failure:: PushMsg ::Http Call Fail",
			setupFunc: func(a args) *httptest.Server {
				svr := testServer("", func(w http.ResponseWriter, r *http.Request) {})
				svr.Close()
				return svr
			},
			reqBody: args{
				msg:     "Hello World",
				uuid:    "b2ae109d-1382-4b1c-a8ab-5a9d04555e4e",
				channel: "channel1",
			},
			ValidateFunc: func(err error) {
				t.Log(err)
				if err == nil {
					t.Errorf("Want: %v, Got: %v", "not nil", nil)
				}
			},
			cleanupFunc: func(svr *httptest.Server) {
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			svr := tt.setupFunc(tt.reqBody)
			defer tt.cleanupFunc(svr)

			calls := NewMsgBrokerSvc(svr.URL)
			err := calls.PushMsg(tt.reqBody.msg, tt.reqBody.uuid, tt.reqBody.channel)

			tt.ValidateFunc(err)
		})
	}
}

func Test_ExtractMessage(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  io.ReadCloser
		setupFunc    func() (io.ReadCloser, *rsa.PrivateKey)
		ValidateFunc func(string, error)
	}{
		{
			name: "Success:: ExtractMsg",
			setupFunc: func() (io.ReadCloser, *rsa.PrivateKey) {
				privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
				if err != nil {
					t.Log(err.Error())
					t.Fail()
					return nil, nil
				}
				publicKey := privateKey.PublicKey
				cipherMsg, err := crypt.RsaOaepEncrypt("Hello, world!", publicKey)
				if err != nil {
					t.Log(err.Error())
					t.Fail()
					return nil, nil
				}
				return io.NopCloser(bytes.NewReader([]byte(cipherMsg))), privateKey
			},
			ValidateFunc: func(msg string, err error) {
				if err != nil {
					t.Errorf("Want: %v, Got: %v", nil, err.Error())
				}
				if msg != "Hello, world!" {
					t.Errorf("Want: %v, Got: %v", "Hello, world!", msg)
				}
			},
		},
		{
			name: "Failure:: ExtractMsg:: Decryption Fail",
			setupFunc: func() (io.ReadCloser, *rsa.PrivateKey) {
				privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
				if err != nil {
					t.Log(err.Error())
					t.Fail()
					return nil, nil
				}
				publicKey := privateKey.PublicKey
				cipherMsg, err := crypt.RsaOaepEncrypt("Hello, world!", publicKey)
				if err != nil {
					t.Log(err.Error())
					t.Fail()
					return nil, nil
				}
				privateKey.E = 123
				return io.NopCloser(bytes.NewReader([]byte(cipherMsg))), privateKey
			},
			ValidateFunc: func(msg string, err error) {
				t.Log(err)
				if err == nil {
					t.Errorf("Want: %v, Got: %v", "not nil", nil)
				}
				if msg != "" {
					t.Errorf("Want: %v, Got: %v", "", msg)
				}
			},
		},
		{
			name: "Success:: ExtractMsg:: Without Encryption",
			setupFunc: func() (io.ReadCloser, *rsa.PrivateKey) {
				cipherMsg := "Hello world!"
				return io.NopCloser(bytes.NewReader([]byte(cipherMsg))), nil
			},
			ValidateFunc: func(msg string, err error) {
				if err != nil {
					t.Errorf("Want: %v, Got: %v", nil, err.Error())
				}
				if msg != "Hello world!" {
					t.Errorf("Want: %v, Got: %v", "Hello, world!", msg)
				}
			},
		},
		{
			name: "Failure:: ExtractMsg:: Source nil",
			setupFunc: func() (io.ReadCloser, *rsa.PrivateKey) {
				return nil, nil
			},
			ValidateFunc: func(msg string, err error) {
				expectedErr := fmt.Errorf("source cannot be nil")
				if err.Error() != expectedErr.Error() {
					t.Errorf("Want: %v, Got: %v", expectedErr.Error(), err.Error())
				}
				if msg != "" {
					t.Errorf("Want: %v, Got: %v", "", msg)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, key := tt.setupFunc()

			calls := NewMsgBrokerSvc("")
			extractMsg := calls.ExtractMsg(key)
			s, err := extractMsg(msg)

			tt.ValidateFunc(s, err)
		})
	}

}
