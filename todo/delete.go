package todo

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tidwall/buntdb"
)

type deleteHandler struct {
	db database
}

// NewDeleteHandler function.
func NewDeleteHandler(db *buntdb.DB) func(w http.ResponseWriter, r *http.Request) {
	return deleteHandler{newDatabase(db)}.delete
}

func (h deleteHandler) delete(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]

	if err := h.db.delete(key); err != nil {
		handleInternalServerError(err, w)
	}
}
