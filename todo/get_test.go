package todo

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getHandler_getAll(t *testing.T) {
	t.Parallel()

	type fields struct {
		getAllItems func() ([]Item, error)
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
			name: "should_return_error_given_get_all_items_failure",
			fields: fields{getAllItems: func() ([]Item, error) {
				return nil, errors.New("error")
			}},
			wantErr: true,
			errMsg:  "failed to get all items: error",
		},
		{
			name: "should_write_status_ok_given_empty_items",
			fields: fields{getAllItems: func() ([]Item, error) {
				return []Item{}, nil
			}},
			want:  http.StatusOK,
			want1: "[]",
		},
		{
			name: "should_write_status_ok_given_non_empty_items",
			fields: fields{getAllItems: func() ([]Item, error) {
				return []Item{{Key: "some-key"}}, nil
			}},
			want:  http.StatusOK,
			want1: `[{"key":"some-key","title":"","completed":false,"url":"http://example.com/some-key","order":0}]`,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			g := getHandler{
				getAllItems: test.fields.getAllItems,
			}
			w := httptest.NewRecorder()

			err := g.getAll(w, httptest.NewRequest("", "/", nil))

			makeAssertions(t, w, err, test.want, test.want1, test.wantErr, test.errMsg)
		})
	}
}

func Test_getHandler_getOne(t *testing.T) {
	t.Parallel()

	type fields struct {
		getOneItem func(string) (Item, error)
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
			name: "should_return_error_given_get_one_item_failure",
			fields: fields{getOneItem: func(string) (Item, error) {
				return Item{}, errors.New("error")
			}},
			wantErr: true,
			errMsg:  "failed to get one item: error",
		},
		{
			name: "should_write_status_ok_given_non_nil_item",
			fields: fields{getOneItem: func(string) (Item, error) {
				return Item{Key: "some-key"}, nil
			}},
			want:  http.StatusOK,
			want1: `{"key":"some-key","title":"","completed":false,"url":"http://example.com/some-key","order":0}`,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			g := getHandler{
				getOneItem: test.fields.getOneItem,
			}
			w := httptest.NewRecorder()

			err := g.getOne(w, httptest.NewRequest("", "/", nil))

			makeAssertions(t, w, err, test.want, test.want1, test.wantErr, test.errMsg)
		})
	}
}

func makeAssertions(t *testing.T, w *httptest.ResponseRecorder, err error, want int, want1 string, wantErr bool, errMsg string) {
	if err != nil {
		assert.True(t, wantErr)
		assert.EqualError(t, err, errMsg)
	} else {
		assert.False(t, wantErr)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()

		assert.Equal(t, want, resp.StatusCode)
		assert.Equal(t, want1, strings.TrimSpace(string(body)))
	}
}
