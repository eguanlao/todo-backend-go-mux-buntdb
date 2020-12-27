package todo

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/tidwall/buntdb"
)

type getHandler struct {
	db *database
}

// NewGetAllHandler function.
func NewGetAllHandler(db *buntdb.DB) func(w http.ResponseWriter, r *http.Request) error {
	return getHandler{newDatabase(db)}.getAll
}

func (h getHandler) getAll(w http.ResponseWriter, r *http.Request) error {
	items, err := h.db.getAll()
	if err != nil {
		return errors.Wrap(err, "failed to get all items")
	}

	for i := range items {
		items[i].URL = getURL(r, items[i].Key)
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(items); err != nil {
		return errors.Wrap(err, "failed to encode items")
	}

	return nil
}

// NewGetOneHandler function.
func NewGetOneHandler(db *buntdb.DB) func(w http.ResponseWriter, r *http.Request) error {
	return getHandler{newDatabase(db)}.getOne
}

func (h getHandler) getOne(w http.ResponseWriter, r *http.Request) error {
	key := mux.Vars(r)["key"]

	item, err := h.db.getOne(key)
	if err != nil {
		return errors.Wrap(err, "failed to get item")
	}

	item.URL = getURL(r, item.Key)

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(item); err != nil {
		return errors.Wrap(err, "failed to encode response")
	}

	return nil
}

func getURL(r *http.Request, key string) string {
	scheme := "http://"
	if os.Getenv("PORT") != "" {
		scheme = "https://"
	}
	return scheme + r.Host + "/" + key
}
