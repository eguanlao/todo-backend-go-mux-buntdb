package todo

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"github.com/tidwall/buntdb"
)

type saveHandler struct {
	generateKey func() string
	saveItem    func(Item) (Item, error)
}

// NewSaveHandler function.
func NewSaveHandler(db *buntdb.DB) func(w http.ResponseWriter, r *http.Request) error {
	h := &saveHandler{func() string { return uuid.New().String() }, newDatabase(db).save}
	return h.save
}

func (h *saveHandler) save(w http.ResponseWriter, r *http.Request) error {
	key := h.generateKey()
	requestBody, _ := ioutil.ReadAll(r.Body)

	item := Item{Key: key}
	if err := json.Unmarshal(requestBody, &item); err != nil {
		return &Error{err, "failed to unmarshal request body", http.StatusBadRequest}
	}

	saved, err := h.saveItem(item)
	if err != nil {
		return &Error{err, "failed to save item", http.StatusInternalServerError}
	}

	saved.URL = getURL(r, item.Key)

	w.Header().Add("Location", saved.URL)
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(saved); err != nil {
		return &Error{err, "failed to encode response", http.StatusInternalServerError}
	}

	return nil
}
