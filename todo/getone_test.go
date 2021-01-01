package todo

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_getOneHandler_getOne(t *testing.T) {
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
			fields: fields{func(string) (Item, error) {
				return Item{}, errors.New("error")
			}},
			wantErr: true,
			errMsg:  "failed to get one item: error",
		},
		{
			name: "should_write_status_ok_given_non_nil_item",
			fields: fields{func(string) (Item, error) {
				return Item{Key: "some-key"}, nil
			}},
			want:  http.StatusOK,
			want1: `{"key":"some-key","title":"","completed":false,"url":"http://example.com/some-key","order":0}`,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			h := &getOneHandler{
				getOneItem: test.fields.getOneItem,
			}
			w := httptest.NewRecorder()

			err := h.getOne(w, httptest.NewRequest("", "/", nil))

			makeAssertions(t, w, err, test.want, test.want1, test.wantErr, test.errMsg)
		})
	}
}
