package parser_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/vatsal278/msgbroker/internal/model"
	"github.com/vatsal278/msgbroker/pkg/parser"
)

func TestParser(t *testing.T) {
	tests := []struct {
		name              string
		requestBody      interface{}
		setupVariable
		expectedResponse interface{}
	}{
		{
			name:"SUCCESS:: Parser",
			requestBody: &model.Publisher{
				Name:    "publisher1",
				Channel: "c4",
			},
			expectedResponse: &model.Publisher{
				Name:    "publisher1",
				Channel: "c4",
			},
		},
		{
			name: "FAILURE:: Parser",
		}
	}
	for _, tt := range tests {
		t.Run("test", func(t *testing.T) {
			//errorCase = tt.ErrorCase
			var publisher model.Publisher
			//w := httptest.NewRecorder()
			//t.Log(w.Code)

			var testPub = &model.Publisher{
				Name:    "publisher1",
				Channel: "c4",
			}
			jsonValue, _ := json.Marshal(testPub)
			r := httptest.NewRequest("POST", "/register/publisher", bytes.NewBuffer(jsonValue))

			parser.Parse(r.Body, &publisher)
			if !reflect.DeepEqual(&publisher, testPub) {
				t.Errorf("Want: %v, Got: %v", testPub, &publisher)
			}

		})

}
