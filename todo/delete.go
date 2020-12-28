package todo

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tidwall/buntdb"
)

type deleteHandler struct {
	deleteItem func(string) error
}

// NewDeleteHandler function.
func NewDeleteHandler(db *buntdb.DB) func(w http.ResponseWriter, r *http.Request) error {
	return deleteHandler{newDatabase(db).delete}.delete
}

func (d deleteHandler) delete(w http.ResponseWriter, r *http.Request) error {
	key := mux.Vars(r)["key"]

	if err := d.deleteItem(key); err != nil {
		return &Error{err, "failed to delete item", http.StatusInternalServerError}
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}
