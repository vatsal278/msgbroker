package sdk

import (
	"crypto/rand"
	"crypto/rsa"
	"github.com/gorilla/mux"
	"github.com/vatsal278/msgbroker/internal/constants"
	parseRequest "github.com/vatsal278/msgbroker/internal/pkg/parser"
	"github.com/vatsal278/msgbroker/model"
	"github.com/vatsal278/msgbroker/pkg/crypt"
	"github.com/vatsal278/msgbroker/pkg/responseWriter"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func testServer(url string, f func(w http.ResponseWriter, r *http.Request)) (*mux.Router, string) {
	router := mux.NewRouter()
	router.HandleFunc(url, f).Methods(http.MethodPost)
	svr := httptest.NewServer(router)
	//url = svr.URL + url

	return router, svr.URL

}
func Test_RegiterSub(t *testing.T) {
	type tempStruct struct {
		method      string
		callbackUrl string
		publicKey   string
		channel     string
	}

	tests := []struct {
		name              string
		requestBody       tempStruct
		mockServerHandler func(http.ResponseWriter, *http.Request)
		ValidateFunc      func(err error)
		expectedResponse  model.Response
	}{
		{
			name:        "Success:: Register Subscriber",
			requestBody: tempStruct{method: "GET", callbackUrl: "http://localhost:8083/pong", publicKey: "", channel: "c1"},
			mockServerHandler: func(w http.ResponseWriter, r *http.Request) {

				//defer sGrp.Done()
				log.Print("HIT")
				var subscriber model.Subscriber
				err := parseRequest.ParseAndValidateRequest(r.Body, &subscriber)
				if err != nil {
					t.Errorf(err.Error())
					return
				}
				responseWriter.ResponseWriter(w, 200, "", "", &model.Response{})
			},
			ValidateFunc: func(err error) {
				if err != nil {
					t.Errorf("Want: %v, Got: %v", nil, err.Error())
				}
			},
		},
		{
			name:        "Failure:: Register Subscriber :: Incorrect Method",
			requestBody: tempStruct{method: "", callbackUrl: "http://localhost:8083/pong", publicKey: ""},
			mockServerHandler: func(w http.ResponseWriter, r *http.Request) {

				responseWriter.ResponseWriter(w, 400, "", "", &model.Response{})
			},
			ValidateFunc: func(err error) {
				if err == nil {
					t.Errorf("Want: %v, Got: %v", "Key: 'Publisher.Channel' Error:Field validation for 'Channel' failed on the 'required' tag", nil)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//sGrp.Add(1)
			_, url := testServer("/register/subscriber", tt.mockServerHandler)
			//defer log.Fatal(router)
			calls := NewController(url)
			//reqBody := tt.requestBody
			//err := calls.RegisterPub("channel")
			err := calls.RegisterSub(tt.requestBody.method, tt.requestBody.callbackUrl, tt.requestBody.publicKey, tt.requestBody.channel)
			t.Log(err)
			tt.ValidateFunc(err)
		})
	}

}

func Test_RegiterPub(t *testing.T) {
	//sGrp := &sync.WaitGroup{}
	tests := []struct {
		name              string
		requestBody       map[string]string
		mockServerHandler func(w http.ResponseWriter, r *http.Request)
		ValidateFunc      func(err error)
		expectedResponse  model.Response
	}{
		{
			name:        "Success:: Register Publisher",
			requestBody: map[string]string{"channel": "c1"},
			mockServerHandler: func(w http.ResponseWriter, r *http.Request) {

				//defer sGrp.Done()
				log.Print("HIT")
				var publisher model.Publisher
				err := parseRequest.ParseAndValidateRequest(r.Body, &publisher)
				if err != nil {
					t.Errorf(err.Error())
					return
				}
				responseWriter.ResponseWriter(w, 200, "", map[string]interface{}{
					"id": "publisher.Id",
				}, &model.Response{})
			},
			ValidateFunc: func(err error) {
				if err != nil {
					t.Errorf("Want: %v, Got: %v", nil, err.Error())
				}
			},
		},
		{
			name:        "Failure:: Register Publisher",
			requestBody: map[string]string{"channel": ""},
			mockServerHandler: func(w http.ResponseWriter, r *http.Request) {
				log.Print("HIT")

				responseWriter.ResponseWriter(w, http.StatusBadRequest, constants.IncompleteData, nil, &model.Response{})
			},
			ValidateFunc: func(err error) {
				if err == nil {
					t.Errorf("Want: %v, Got: %v", "Key: 'Publisher.Channel' Error:Field validation for 'Channel' failed on the 'required' tag", nil)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			//sGrp.Add(1)
			_, url := testServer("/register/publisher", tt.mockServerHandler)
			//defer log.Fatal(router)
			calls := NewController(url)
			//reqBody := tt.requestBody
			key, err := calls.RegisterPub("channel")
			t.Log(err)
			t.Log(key)
			tt.ValidateFunc(err)
			//t.Log(key)

		})
	}
}
func Test_UpdateSubs(t *testing.T) {
	type tempStruct struct {
		msg     string
		key     string
		channel string
	}
	tests := []struct {
		name              string
		requestBody       tempStruct
		mockServerHandler func(w http.ResponseWriter, r *http.Request)
		ValidateFunc      func(error)
		expectedResponse  model.Response
	}{
		{
			name:        "Success:: PushMsg",
			requestBody: tempStruct{msg: "Hello World", channel: "c1"},
			mockServerHandler: func(w http.ResponseWriter, r *http.Request) {

				//defer sGrp.Done()
				log.Print("HIT")
				var update model.Updates
				err := parseRequest.ParseAndValidateRequest(r.Body, &update)
				if err != nil {
					t.Errorf(err.Error())
					return
				}
				responseWriter.ResponseWriter(w, 200, "", nil, &model.Response{})
			},
			ValidateFunc: func(err error) {
				if err != nil {
					t.Errorf("Want: %v, Got: %v", nil, err.Error())
				}
			},
		},
		{
			name:        "Failure::PushMsg",
			requestBody: tempStruct{msg: "Hello World", channel: ""},
			mockServerHandler: func(w http.ResponseWriter, r *http.Request) {

				//defer sGrp.Done()
				log.Print("HIT")

				responseWriter.ResponseWriter(w, 400, "", nil, &model.Response{})
			},
			ValidateFunc: func(err error) {
				if err == nil {
					t.Errorf("Want: %v, Got: %v", "non success status code received : 400", nil)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, url := testServer("/publish", tt.mockServerHandler)
			//defer log.Fatal(router)
			calls := NewController(url)
			reqBody := tt.requestBody
			err := calls.PushMsg(reqBody.msg, "", reqBody.channel)
			t.Log(err)

			tt.ValidateFunc(err)
		})
	}

}

/*
func testClient(c *TestServer, encrypted bool) {
	//expected := "dummy data"
	calls := NewController("http://localhost:9090")
	var privateKey *rsa.PrivateKey
	var err error

	if encrypted {
		privateKey, err = rsa.GenerateKey(rand.Reader, 2048)

		if err != nil {
			c.t.Log(err.Error())
		}

	}
	router := mux.NewRouter()
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		x := calls.ExtractMsg(r.Body, privateKey)
		msg := x()
		c.t.Errorf(msg)
	}).Methods(http.MethodPost)
	svr := httptest.NewServer(router)
	var url = svr.URL + "/ping"
	c.srv = svr
	c.t.Log(privateKey)
	var pubKey = ""
	if privateKey != nil {
		publicKey := privateKey.PublicKey
		pubKey = crypt.KeyAsPEMStr(&publicKey)
		c.t.Log(pubKey)
	}
}

type TestServer struct {
	srv *httptest.Server
	t   *testing.T
	i   controllerInterface.IController
	wg  *sync.WaitGroup
}*/

func Test_ReceiveMessage(t *testing.T) {
	type tempStruct struct {
		msg     string
		key     string
		channel string
	}
	tests := []struct {
		name              string
		requestBody       io.ReadCloser
		setupFunc         func() (string, *rsa.PrivateKey)
		mockServerHandler func(http.ResponseWriter, *http.Request)
		ValidateFunc      func(error)
		expectedResponse  model.Response
	}{
		{
			name:        "Success:: Update Subscribers",
			requestBody: io.NopCloser(strings.NewReader("Hello, world!")),
			setupFunc: func() (string, *rsa.PrivateKey) {
				privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
				if err != nil {
					t.Log(err.Error())
					return "", nil
				}
				publicKey := privateKey.PublicKey
				cipherMsg, err := crypt.RsaOaepEncrypt("Hello, world!", publicKey)
				if err != nil {
					t.Log(err.Error())
					return "", nil
				}
				return cipherMsg, privateKey
			},
			mockServerHandler: func(w http.ResponseWriter, r *http.Request) {

				//defer sGrp.Done()
				log.Print("HIT")
				var update model.Updates
				err := parseRequest.ParseAndValidateRequest(r.Body, &update)
				if err != nil {
					t.Errorf(err.Error())
					return
				}
				responseWriter.ResponseWriter(w, 200, "", nil, &model.Response{})
			},
			ValidateFunc: func(err error) {
				if err != nil {
					t.Errorf("Want: %v, Got: %v", nil, err.Error())
				}
			},
		},
		{
			name:        "Failure:: Update Subscribers",
			requestBody: io.NopCloser(strings.NewReader("Hello, world!")),
			setupFunc: func() (string, *rsa.PrivateKey) {
				privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
				if err != nil {
					t.Log(err.Error())
					return "", nil
				}
				publicKey := privateKey.PublicKey
				cipherMsg, err := crypt.RsaOaepEncrypt("Hello, world!", publicKey)
				if err != nil {
					t.Log(err.Error())
					return "", nil
				}
				privateKey.E = 100
				return cipherMsg, privateKey
			},
			mockServerHandler: func(w http.ResponseWriter, r *http.Request) {

				//defer sGrp.Done()
				log.Print("HIT")
				var update model.Updates
				err := parseRequest.ParseAndValidateRequest(r.Body, &update)
				if err != nil {
					t.Errorf(err.Error())
					return
				}
				responseWriter.ResponseWriter(w, 200, "", nil, &model.Response{})
			},
			ValidateFunc: func(err error) {
				if err != nil {
					t.Errorf("Want: %v, Got: %v", nil, err.Error())
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, url := testServer("/publish", tt.mockServerHandler)
			//defer log.Fatal(router)
			calls := NewController(url)
			reqBody := tt.requestBody
			err := calls.ExtractMsg(key)
			t.Log(err)
			getMsg := c
			tt.ValidateFunc(getMsg(readClosure))
		})
	}

}
