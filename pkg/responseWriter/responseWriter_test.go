package responseWriter_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vatsal278/msgbroker/internal/model"
	"github.com/vatsal278/msgbroker/pkg/responseWriter"
)

func TestResponseWriter(t *testing.T) {
	t.Run("TestValidate", func(t *testing.T) {
		w := httptest.NewRecorder()
		//r := httptest.
		_ = responseWriter.ResponseWriter(w, http.StatusOK, "Successfully Registered as publisher to the channel", nil, &model.Response{})
		if w.Code != http.StatusOK {
			t.Errorf("Want: %v, Got: %v", http.StatusOK, w.Code)
		}
	},
	)
}
