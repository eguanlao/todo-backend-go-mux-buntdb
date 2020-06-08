package todo

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"github.com/tidwall/buntdb"
)

type saveHandler struct {
	db database
}

// NewSaveHandler function.
func NewSaveHandler(db *buntdb.DB) func(w http.ResponseWriter, r *http.Request) {
	return saveHandler{newDatabase(db)}.save
}

func (h saveHandler) save(w http.ResponseWriter, r *http.Request) {
	key := uuid.New().String()
	requestBytes, _ := ioutil.ReadAll(r.Body)

	item := Item{
		Key: key,
	}
	if err := json.Unmarshal(requestBytes, &item); err != nil {
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
