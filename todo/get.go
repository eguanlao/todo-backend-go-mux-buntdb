package todo

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/tidwall/buntdb"
)

type getHandler struct {
	getAllItems func() ([]Item, error)
	getOneItem  func(string) (Item, error)
}

// NewGetAllHandler function.
func NewGetAllHandler(db *buntdb.DB) func(w http.ResponseWriter, r *http.Request) error {
	return getHandler{getAllItems: newDatabase(db).getAll}.getAll
}

func (g getHandler) getAll(w http.ResponseWriter, r *http.Request) error {
	items, err := g.getAllItems()
	if err != nil {
		return &Error{err, "failed to get all items", http.StatusInternalServerError}
	}

	for i := range items {
		items[i].URL = getURL(r, items[i].Key)
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(items); err != nil {
		return &Error{err, "failed to encode response", http.StatusInternalServerError}
	}

	return nil
}

// NewGetOneHandler function.
func NewGetOneHandler(db *buntdb.DB) func(w http.ResponseWriter, r *http.Request) error {
	return getHandler{getOneItem: newDatabase(db).getOne}.getOne
}

func (g getHandler) getOne(w http.ResponseWriter, r *http.Request) error {
	key := mux.Vars(r)["key"]

	item, err := g.getOneItem(key)
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

func getURL(r *http.Request, key string) string {
	scheme := "http://"
	if os.Getenv("PORT") != "" {
		scheme = "https://"
	}
	return scheme + r.Host + "/" + key
}
