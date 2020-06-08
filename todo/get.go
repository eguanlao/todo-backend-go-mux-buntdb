package todo

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/tidwall/buntdb"
)

type getHandler struct {
	db database
}

// NewGetAllHandler function.
func NewGetAllHandler(db *buntdb.DB) func(w http.ResponseWriter, r *http.Request) {
	return getHandler{newDatabase(db)}.getAll
}

func (h getHandler) getAll(w http.ResponseWriter, r *http.Request) {
	items, err := h.db.getAll()
	if err != nil {
		handleInternalServerError(err, w)
		return
	}

	for i := range items {
		items[i].URL = getURL(r, items[i].Key)
	}

	if err := json.NewEncoder(w).Encode(items); err != nil {
		handleInternalServerError(err, w)
	}
}

// NewGetOneHandler function.
func NewGetOneHandler(db *buntdb.DB) func(w http.ResponseWriter, r *http.Request) {
	return getHandler{newDatabase(db)}.getOne
}

func (h getHandler) getOne(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]

	item, err := h.db.getOne(key)
	if err != nil {
		handleInternalServerError(err, w)
		return
	}

	item.URL = getURL(r, item.Key)

	if err := json.NewEncoder(w).Encode(item); err != nil {
		handleInternalServerError(err, w)
	}
}

func getURL(r *http.Request, key string) string {
	scheme := "http://"
	if os.Getenv("PORT") != "" {
		scheme = "https://"
	}
	return scheme + r.Host + "/" + key
}
