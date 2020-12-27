package todo

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"github.com/tidwall/buntdb"
)

type updateHandler struct {
	db *database
}

// NewUpdateHandler function.
func NewUpdateHandler(db *buntdb.DB) func(w http.ResponseWriter, r *http.Request) error {
	return updateHandler{newDatabase(db)}.update
}

func (h updateHandler) update(w http.ResponseWriter, r *http.Request) error {
	key := mux.Vars(r)["key"]
	requestBytes, _ := ioutil.ReadAll(r.Body)

	update := Item{}
	if err := json.Unmarshal(requestBytes, &update); err != nil {
		return errors.Wrap(err, "failed to unmarshal request body")
	}

	item, err := h.db.getOne(key)
	if err != nil {
		return errors.Wrap(err, "failed to get item")
	}

	if err = mergo.Merge(&item, update, mergo.WithOverride); err != nil {
		handleInternalServerError(err, w)
		return errors.Wrap(err, "failed to merge structs")
	}

	saved, err := h.db.save(item)
	if err != nil {
		return errors.Wrap(err, "failed to save item")
	}

	saved.URL = getURL(r, item.Key)

	if err := json.NewEncoder(w).Encode(saved); err != nil {
		return errors.Wrap(err, "failed to encode response")
	}

	return nil
}
