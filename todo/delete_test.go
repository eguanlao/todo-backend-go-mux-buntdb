package todo

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_deleteHandler_delete(t *testing.T) {
	t.Parallel()

	type fields struct {
		deleteItem func(string) error
	}

	tests := []struct {
		name    string
		fields  fields
		want    int
		want1   string
		wantErr bool
		errMsg  string
	}{
		{
			name: "should_return_error_given_delete_item_failure",
			fields: fields{deleteItem: func(string) error {
				return errors.New("error")
			}},
			wantErr: true,
			errMsg:  "failed to delete item: error",
		},
		{
			name: "should_write_status_no_content",
			fields: fields{deleteItem: func(string) error {
				return nil
			}},
			want:  http.StatusNoContent,
			want1: "",
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			h := deleteHandler{
				deleteItem: test.fields.deleteItem,
			}
			w := httptest.NewRecorder()

			err := h.delete(w, httptest.NewRequest("", "/", nil))

			if err != nil {
				assert.True(t, test.wantErr)
				assert.EqualError(t, err, test.errMsg)
			} else {
				assert.False(t, test.wantErr)

				resp := w.Result()
				body, _ := ioutil.ReadAll(resp.Body)
				defer resp.Body.Close()

				assert.Equal(t, test.want, resp.StatusCode)
				assert.Equal(t, test.want1, string(body))
			}
		})
	}
}
