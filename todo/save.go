package todo

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/tidwall/buntdb"
)

type saveHandler struct {
	db *database
}

// NewSaveHandler function.
func NewSaveHandler(db *buntdb.DB) func(w http.ResponseWriter, r *http.Request) error {
	return saveHandler{newDatabase(db)}.save
}

func (h saveHandler) save(w http.ResponseWriter, r *http.Request) error {
	key := uuid.New().String()
	requestBytes, _ := ioutil.ReadAll(r.Body)

	item := Item{
		Key: key,
	}
	if err := json.Unmarshal(requestBytes, &item); err != nil {
		return errors.Wrap(err, "failed to unmarshal request bytes")
	}

	saved, err := h.db.save(item)
	if err != nil {
		return errors.Wrap(err, "failed to save item")
	}

	saved.URL = getURL(r, item.Key)

	w.Header().Add("Location", saved.URL)
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(saved); err != nil {
		return errors.Wrap(err, "failed to encode response")
	}

	return nil
}
