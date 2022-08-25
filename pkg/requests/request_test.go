package requests

import (
	"crypto/rand"
	"crypto/rsa"
	"github.com/vatsal278/msgbroker/internal/model"
	"github.com/vatsal278/msgbroker/pkg/crypt"
	"io"
	"net/http"
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
		ValidateFunc     func(r *http.Response)
		expectedResponse model.Response
	}{
		{
			name:        "Success:: Register Subscriber",
			requestBody: tempStruct{method: "GET", callbackUrl: "http://localhost:8083/pong", publicKey: "", channel: "c1"},
			ValidateFunc: func(r *http.Response) {
				if r.Status != "201 Created" {
					t.Errorf("Want: %v, Got: %v", "201 Created", r.Status)
				}
			},
		},
		{
			name:        "Success :: Register Subscriber :: With Encryption",
			requestBody: tempStruct{method: "GET", callbackUrl: "http://localhost:8083/pong", publicKey: "LS0tLS1CRUdJTiBSU0EgUFVCTElDIEtFWS0tLS0tCk1JSUJDZ0tDQVFFQXZBWmZxM1lvVzdUTzBGYmJHMWxxRVBxNHQ4bGc5cTdla0NYMXJIVjVNNTdobmdyNlF1L3MKTnp0QXkzTmh1TG4xSm5PSVN5bzRXc29MMDRKWFI5WXI5UXVtZW1EdGVreWpOd2toQkFWM0xBN3BORjV3c2ZaSwpFbC9jY2U5aGZxRWtOcERtNUFFZklnRW5UZXdTMml5cGRCQm1pVmI5VzNzZFdUWHEwenNKY1pqb29obXZPNkN1CngyY01NOW1EeFQ4VXBYM2gweE1WNTBVd050TzRVbS9aWnFPeENqdFdhNE1STE16NTNMTG9lUm9UOE1tZEdlV1UKYTdHMitKU0c5K3V1MVJIVkYrelZGaEx2emtoM3dLTGdVdU1DcW0rL1U0Y3B3TDUxZU9TYVZNYUhjU1NiRXZCUgp0d0lZdHRHR3NDVC9mTEdyVXdjZm8xZ0xKaVNjU2taN1B3SURBUUFCCi0tLS0tRU5EIFJTQSBQVUJMSUMgS0VZLS0tLS0K", channel: "c1"},
			ValidateFunc: func(r *http.Response) {
				if r.Status != "201 Created" {
					t.Errorf("Want: %v, Got: %v", "201 Created", r.Status)
				}
			},
		},
		{
			name:        "Failure:: Register Subscriber :: Incorrect Method",
			requestBody: tempStruct{method: "POST", callbackUrl: "http://localhost:8083/pong", publicKey: ""},
			ValidateFunc: func(r *http.Response) {
				if r.Status != "400 Bad Request" {
					t.Errorf("Want: %v, Got: %v", "400 Bad Request", r.Status)
					return
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r, err := calls.RegisterSub(tt.requestBody.method, tt.requestBody.callbackUrl, tt.requestBody.publicKey, tt.requestBody.channel)
			if err != nil {
				t.Log(err.Error())
			}
			t.Log(r.Status)
			tt.ValidateFunc(r)
		})
	}

}
func Test_RegiterPub(t *testing.T) {
	calls := NewController("http://localhost:9090")
	tests := []struct {
		name             string
		requestBody      string
		ValidateFunc     func(*http.Response)
		expectedResponse model.Response
	}{
		{
			name:        "Success:: Register Publisher",
			requestBody: "c1",
			ValidateFunc: func(r *http.Response) {
				if r.Status != "201 Created" {
					t.Errorf("Want: %v, Got: %v", "400 Bad Request", r.Status)
					return
				}
			},
		},
		{
			name:        "Failure:: Register Publisher",
			requestBody: "",
			ValidateFunc: func(r *http.Response) {
				if r.Status != "400 Bad Request" {
					t.Errorf("Want: %v, Got: %v", "400 Bad Request", r.Status)
					return
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reqBody string
			reqBody = tt.requestBody
			key, r, err := calls.RegisterPub(reqBody)
			if err != nil {
				t.Errorf("Want: %v, Got: %v", nil, err.Error())
			}
			t.Log(key)
			tt.ValidateFunc(r)
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
		name             string
		requestBody      tempStruct
		setupFunc        func() string
		ValidateFunc     func(*http.Response)
		expectedResponse model.Response
	}{
		{
			name:        "Success:: Update Subscribers",
			requestBody: tempStruct{msg: "Hello World", key: "http://localhost:8083/pong", channel: "c1"},
			setupFunc: func() string {
				y := "c1"
				z, _, err := calls.RegisterPub(y)
				if err != nil {
					t.Log(err.Error())
				}
				return z
			},
			ValidateFunc: func(r *http.Response) {
				if r.Status != "200 OK" {
					t.Errorf("Want: %v, Got: %v", "200 OK", r.Status)
					return
				}
			},
		},
		{
			name:        "Failure:: Update Subscribers",
			requestBody: tempStruct{msg: "Hello World", key: "", channel: "c1"},
			setupFunc: func() string {
				y := "c1"
				z, _, err := calls.RegisterPub(y)
				if err != nil {
					t.Log(err.Error())
				}
				return z
			},
			ValidateFunc: func(r *http.Response) {
				if r.Status != "200 OK" {
					t.Errorf("Want: %v, Got: %v", "200 OK", r.Status)
					return
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			z := tt.setupFunc()
			r, err := calls.UpdateSubs("HelloWorld", z, "c1")
			if err != nil {
				t.Log(err.Error())
			}
			tt.ValidateFunc(r)
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
		x := calls.ReceiveMsg(r.Body, privateKey)
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
			getMsg := calls.ReceiveMsg(readClosure, key)
			tt.ValidateFunc(getMsg())
		})
	}

}
