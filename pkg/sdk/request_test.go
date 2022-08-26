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

func Test_RegiterSub(t *testing.T) {
	type tempStruct struct {
		method      string
		callbackUrl string
		publicKey   string
		channel     string
	}
	calls := NewController("http://localhost:9090")
	tests := []struct {
		name             string
		requestBody      tempStruct
		ValidateFunc     func(err error)
		expectedResponse model.Response
	}{
		{
			name:        "Success:: Register Subscriber",
			requestBody: tempStruct{method: "GET", callbackUrl: "http://localhost:8083/pong", publicKey: "", channel: "c1"},
			ValidateFunc: func(err error) {
				if err != nil {
					t.Errorf("Want: %v, Got: %v", nil, err.Error())
				}
			},
		},
		{
			name:        "Success :: Register Subscriber :: With Encryption",
			requestBody: tempStruct{method: "GET", callbackUrl: "http://localhost:8083/pong", publicKey: "LS0tLS1CRUdJTiBSU0EgUFVCTElDIEtFWS0tLS0tCk1JSUJDZ0tDQVFFQXZBWmZxM1lvVzdUTzBGYmJHMWxxRVBxNHQ4bGc5cTdla0NYMXJIVjVNNTdobmdyNlF1L3MKTnp0QXkzTmh1TG4xSm5PSVN5bzRXc29MMDRKWFI5WXI5UXVtZW1EdGVreWpOd2toQkFWM0xBN3BORjV3c2ZaSwpFbC9jY2U5aGZxRWtOcERtNUFFZklnRW5UZXdTMml5cGRCQm1pVmI5VzNzZFdUWHEwenNKY1pqb29obXZPNkN1CngyY01NOW1EeFQ4VXBYM2gweE1WNTBVd050TzRVbS9aWnFPeENqdFdhNE1STE16NTNMTG9lUm9UOE1tZEdlV1UKYTdHMitKU0c5K3V1MVJIVkYrelZGaEx2emtoM3dLTGdVdU1DcW0rL1U0Y3B3TDUxZU9TYVZNYUhjU1NiRXZCUgp0d0lZdHRHR3NDVC9mTEdyVXdjZm8xZ0xKaVNjU2taN1B3SURBUUFCCi0tLS0tRU5EIFJTQSBQVUJMSUMgS0VZLS0tLS0K", channel: "c1"},
			ValidateFunc: func(err error) {
				if err != nil {
					t.Errorf("Want: %v, Got: %v", nil, err.Error())
				}
			},
		},
		{
			name:        "Failure:: Register Subscriber :: Incorrect Method",
			requestBody: tempStruct{method: "POST", callbackUrl: "http://localhost:8083/pong", publicKey: ""},
			ValidateFunc: func(err error) {
				if err != nil {
					t.Errorf("Want: %v, Got: %v", "400 Bad Request", nil)
					return
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := calls.RegisterSub(tt.requestBody.method, tt.requestBody.callbackUrl, tt.requestBody.publicKey, tt.requestBody.channel)
			tt.ValidateFunc(err)
		})
	}

}
func testServer(url string, f func(w http.ResponseWriter, r *http.Request)) (*mux.Router, string) {
	router := mux.NewRouter()
	router.HandleFunc(url, f).Methods(http.MethodPost)
	svr := httptest.NewServer(router)
	//url = svr.URL + url

	return router, svr.URL

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
			reqBody := tt.requestBody
			key, err := calls.RegisterPub(reqBody["channel"])
			t.Log(err)
			t.Log(key)
			tt.ValidateFunc(err)
			//t.Log(key)

		})
	}
}
func Test_UpdateSubs(t *testing.T) {
	calls := NewController("http://localhost:9090")
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
			name:        "Success:: Update Subscribers",
			requestBody: tempStruct{msg: "Hello World", channel: "c1"},
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
			name:        "Failure:: Update Subscribers",
			mockServerHandler: func(w http.ResponseWriter, r *http.Request)
			ValidateFunc: func(err error) {
				if err == nil {
					t.Errorf("Want: %v, Got: %v", "error", nil)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, url := testServer("/register/publisher", tt.mockServerHandler)
			//defer log.Fatal(router)
			calls := NewController(url)
			reqBody := tt.requestBody
			key, err := calls.PushMsg()

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
	calls := NewController("http://localhost:9090")
	type tempStruct struct {
		msg     string
		key     string
		channel string
	}
	tests := []struct {
		name             string
		requestBody      io.ReadCloser
		setupFunc        func() (string, *rsa.PrivateKey)
		ValidateFunc     func(string, error)
		expectedResponse model.Response
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
			ValidateFunc: func(msg string, err error) {
				if msg != "Hello, world!" {
					t.Errorf("Want: %v, Got: %v", "Hello, world!", msg)
					return
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
			ValidateFunc: func(msg string, err error) {
				if msg != "" {
					t.Errorf("Want: %v, Got: %v", "Hello, world!", msg)
					return
				}
				//testErr:=errors.New("crypto/rsa: decryption error")
				if err.Error() != "crypto/rsa: decryption error" {
					t.Errorf("Want: %v, Got: %v", "crypto/rsa: decryption error", msg)
					return
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x, key := tt.setupFunc()
			readClosure := io.NopCloser(strings.NewReader(x))
			getMsg := calls.ExtractMsg(key)
			tt.ValidateFunc(getMsg(readClosure))
		})
	}

}
