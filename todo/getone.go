package todo

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tidwall/buntdb"
)

type getOneHandler struct {
	getOneItem func(string) (Item, error)
}

// NewGetOneHandler function.
func NewGetOneHandler(db *buntdb.DB) func(w http.ResponseWriter, r *http.Request) error {
	h := &getOneHandler{newDatabase(db).getOne}
	return h.getOne
}

func (h *getOneHandler) getOne(w http.ResponseWriter, r *http.Request) error {
	key := mux.Vars(r)["key"]

	item, err := h.getOneItem(key)
	if err != nil {
		return &Error{err, "failed to get one item", http.StatusInternalServerError}
	}

	item.URL = getURL(r, item.Key)

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(item); err != nil {
		return &Error{err, "failed to encode response", http.StatusInternalServerError}
	}

	return nil
}
