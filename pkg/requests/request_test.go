package requests

import (
	"github.com/vatsal278/msgbroker/internal/handler"
	"github.com/vatsal278/msgbroker/internal/model"
	"net/http/httptest"
	"testing"
)

func Test_RegiterSub(t *testing.T) {
	type tempStruct struct {
		method      string
		callbackUrl string
		publicKey   string
		channel     string
	}
	NewController("http://localhost:9090")
	tests := []struct {
		name             string
		requestBody      tempStruct
		ValidateFunc     func(*httptest.ResponseRecorder, controllerInterface.IController, interface{})
		expectedResponse model.Response
	}{
		{
			name:        "Success:: Register Publisher",
			requestBody: tempStruct{method: "GET", callbackUrl: "http://localhost:8083/pong", publicKey: "", channel: "c1"},
			ValidateFunc: func(w *httptest.ResponseRecorder, i controllerInterface.IController, reqbody interface{}) {

			},
		},
		{
			name:        "Success :: Register Publisher :: With Encryption",
			requestBody: tempStruct{method: "GET", callbackUrl: "http://localhost:8083/pong", publicKey: "LS0tLS1CRUdJTiBSU0EgUFVCTElDIEtFWS0tLS0tCk1JSUJDZ0tDQVFFQXZBWmZxM1lvVzdUTzBGYmJHMWxxRVBxNHQ4bGc5cTdla0NYMXJIVjVNNTdobmdyNlF1L3MKTnp0QXkzTmh1TG4xSm5PSVN5bzRXc29MMDRKWFI5WXI5UXVtZW1EdGVreWpOd2toQkFWM0xBN3BORjV3c2ZaSwpFbC9jY2U5aGZxRWtOcERtNUFFZklnRW5UZXdTMml5cGRCQm1pVmI5VzNzZFdUWHEwenNKY1pqb29obXZPNkN1CngyY01NOW1EeFQ4VXBYM2gweE1WNTBVd050TzRVbS9aWnFPeENqdFdhNE1STE16NTNMTG9lUm9UOE1tZEdlV1UKYTdHMitKU0c5K3V1MVJIVkYrelZGaEx2emtoM3dLTGdVdU1DcW0rL1U0Y3B3TDUxZU9TYVZNYUhjU1NiRXZCUgp0d0lZdHRHR3NDVC9mTEdyVXdjZm8xZ0xKaVNjU2taN1B3SURBUUFCCi0tLS0tRU5EIFJTQSBQVUJMSUMgS0VZLS0tLS0K", channel: "c1"},
			ValidateFunc: func(w *httptest.ResponseRecorder, i controllerInterface.IController, reqbody interface{}) {

			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ApiCalls.RegisterSub(tt.requestBody.method, tt.requestBody.callbackUrl, tt.requestBody.publicKey, tt.requestBody.channel)
			if err != nil {
				t.Log(err.Error())
			}
		})
	}

}
func Test_RegiterPub(t *testing.T) {

	key, err := ApiCalls.RegisterPub("c1")
	if err != nil {
		t.Log(err.Error())
	}

	t.Log(key)
}
func Test_UpdateSubs(t *testing.T) {
	//x := Publisher{Channel: "c1"}
	y := "c1"
	z, err := ApiCalls.RegisterPub(y)

	err = ApiCalls.UpdateSubs("HelloWorld", z, "c1")
	if err != nil {
		t.Log(err.Error())
	}

}
