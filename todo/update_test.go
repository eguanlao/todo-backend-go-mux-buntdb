package todo

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_updateHandler_update(t *testing.T) {
	t.Parallel()

	type fields struct {
		getOneItem func(string) (Item, error)
		saveItem   func(Item) (Item, error)
	}

	type args struct {
		body io.Reader
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		want1   string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "should_return_error_given_unmarshal_failure",
			wantErr: true,
			errMsg:  "failed to unmarshal request body: unexpected end of JSON input",
		},
		{
			name: "should_return_error_given_get_one_item_failure",
			fields: fields{getOneItem: func(string) (Item, error) {
				return Item{}, errors.New("error")
			}},
			args:    args{strings.NewReader(`{}`)},
			wantErr: true,
			errMsg:  "failed to get one item: error",
		},
		{
			name: "should_return_error_given_save_item_failure",
			fields: fields{
				getOneItem: func(string) (Item, error) {
					return Item{}, nil
				},
				saveItem: func(Item) (Item, error) {
					return Item{}, errors.New("error")
				},
			},
			args:    args{strings.NewReader(`{}`)},
			wantErr: true,
			errMsg:  "failed to save item: error",
		},
		{
			name: "should_write_status_ok",

			fields: fields{
				getOneItem: func(string) (Item, error) {
					return Item{Key: "some-key"}, nil
				},
				saveItem: func(Item) (Item, error) {
					return Item{Key: "some-key"}, nil
				},
			},
			args:  args{strings.NewReader(`{}`)},
			want:  http.StatusOK,
			want1: `{"key":"some-key","title":"","completed":false,"url":"http://example.com/some-key","order":0}`,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			h := &updateHandler{
				getOneItem: test.fields.getOneItem,
				saveItem:   test.fields.saveItem,
			}
			w := httptest.NewRecorder()
			req := httptest.NewRequest("", "/", test.args.body)

			err := h.update(w, req)

			if err != nil {
				assert.True(t, test.wantErr)
				assert.EqualError(t, err, test.errMsg)
			} else {
				assert.False(t, test.wantErr)

				resp := w.Result()
				body, _ := ioutil.ReadAll(resp.Body)
				defer resp.Body.Close()

				assert.Equal(t, test.want, resp.StatusCode)
				assert.Equal(t, test.want1, strings.TrimSpace(string(body)))
			}
		})
	}
}
