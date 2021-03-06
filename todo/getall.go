package todo

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/tidwall/buntdb"
)

type getAllHandler struct {
	getAllItems func() ([]Item, error)
}

// NewGetAllHandler function.
func NewGetAllHandler(db *buntdb.DB) func(w http.ResponseWriter, r *http.Request) error {
	h := &getAllHandler{newDatabase(db).getAll}
	return h.getAll
}

func (h *getAllHandler) getAll(w http.ResponseWriter, r *http.Request) error {
	items, err := h.getAllItems()
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

func getURL(r *http.Request, key string) string {
	scheme := "http://"
	if os.Getenv("PORT") != "" {
		scheme = "https://"
	}
	return scheme + r.Host + "/" + key
}
