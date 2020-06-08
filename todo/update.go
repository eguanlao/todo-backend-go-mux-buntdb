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
	db database
}

// NewUpdateHandler function.
func NewUpdateHandler(db *buntdb.DB) func(w http.ResponseWriter, r *http.Request) {
	return updateHandler{newDatabase(db)}.update
}

func (h updateHandler) update(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	requestBody, _ := ioutil.ReadAll(r.Body)

	item, err := h.db.getOne(key)
	if err != nil {
		handleInternalServerError(err, w)
		return
	}

	update := Item{}
	if err = json.Unmarshal([]byte(string(requestBody)), &update); err != nil {
		handleBadRequest(err, w)
		return
	}

	if err = mergo.Merge(&item, update, mergo.WithOverride); err != nil {
		handleInternalServerError(err, w)
		return
	}

	saved, err := h.db.save(item)
	if err != nil {
		handleInternalServerError(err, w)
		return
	}

	saved.URL = getURL(r, item.Key)

	if err := json.NewEncoder(w).Encode(saved); err != nil {
		handleInternalServerError(err, w)
	}
}
