package todo

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/tidwall/buntdb"
)

type deleteHandler struct {
	db *database
}

// NewDeleteHandler function.
func NewDeleteHandler(db *buntdb.DB) func(w http.ResponseWriter, r *http.Request) error {
	return deleteHandler{newDatabase(db)}.delete
}

func (h deleteHandler) delete(w http.ResponseWriter, r *http.Request) error {
	key := mux.Vars(r)["key"]

	if err := h.db.delete(key); err != nil {
		return errors.Wrap(err, "failed to delete item")
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}
