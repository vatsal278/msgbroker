package responseWriter

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/vatsal278/msgbroker/mocks"
)

func TestResponseWriter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		w      http.ResponseWriter
		status int
		msg    string
		data   interface{}
	}
	tests := []struct {
		name     string
		input    args
		mocks    func(*mocks.MockResponse)
		validate func(http.ResponseWriter, error)
	}{
		{
			name: "SUCCESS:: ResponseWriter",
			input: args{
				w:      httptest.NewRecorder(),
				status: http.StatusOK,
				msg:    "OK",
				data:   "Hello",
			},
			mocks: func(controller *mocks.MockResponse) {
				controller.EXPECT().Update(http.StatusOK, "OK", "Hello").MaxTimes(1)
			},
			validate: func(w http.ResponseWriter, err error) {
				if err != nil {
					t.Errorf("Error::Want: %v, Got: %v", nil, err.Error())
				}
				v, ok := w.(*httptest.ResponseRecorder)
				if !ok {
					t.Errorf("ResponseWriter assert failed")
				}
				if v.Code != http.StatusOK {
					t.Errorf("Want: %v, Got: %v", http.StatusOK, v.Code)
				}
				if v.Header().Get("Content-Type") != "application/json" {
					t.Errorf("Content Type Want: %v, Got: %v", "application/json", v.Header().Get("Content-Type"))
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// registering the mocks
			MockResponse := mocks.NewMockResponse(ctrl)
			tt.mocks(MockResponse)
			// function call
			err := ResponseWriter(tt.input.w, tt.input.status, tt.input.msg, tt.input.data, MockResponse)
			// validate
			tt.validate(tt.input.w, err)
		})
	}
}
