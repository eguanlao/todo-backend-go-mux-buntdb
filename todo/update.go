package todo

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/imdario/mergo"
	"github.com/tidwall/buntdb"
)

type updateHandler struct {
	getOneItem func(string) (Item, error)
	saveItem   func(Item) (Item, error)
}

// NewUpdateHandler function.
func NewUpdateHandler(db *buntdb.DB) func(w http.ResponseWriter, r *http.Request) error {
	return updateHandler{newDatabase(db).getOne, newDatabase(db).save}.update
}

func (u updateHandler) update(w http.ResponseWriter, r *http.Request) error {
	key := mux.Vars(r)["key"]
	requestBody, _ := ioutil.ReadAll(r.Body)

	update := Item{}
	if err := json.Unmarshal(requestBody, &update); err != nil {
		return &Error{err, "failed to unmarshal request body", http.StatusBadRequest}
	}

	item, err := u.getOneItem(key)
	if err != nil {
		return &Error{err, "failed to get one item", http.StatusInternalServerError}
	}

	if err = mergo.Merge(&item, update, mergo.WithOverride); err != nil {
		return &Error{err, "failed to merge structs", http.StatusInternalServerError}
	}

	saved, err := u.saveItem(item)
	if err != nil {
		return &Error{err, "failed to save item", http.StatusInternalServerError}
	}

	saved.URL = getURL(r, item.Key)

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(saved); err != nil {
		return &Error{err, "failed to encode response", http.StatusInternalServerError}
	}

	return nil
}
