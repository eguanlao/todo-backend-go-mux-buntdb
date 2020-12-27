package main

import (
	"log"
	"net/http"
	"os"

	"github.com/eguanlao/todo-backend-go-mux-buntdb/todo"
	"github.com/gorilla/mux"
	"github.com/tidwall/buntdb"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("env var PORT is empty")
	}

	db, err := buntdb.Open(":memory:")
	if err != nil {
		log.Fatal(err)
	}
	if err = db.CreateIndex("order", "*", buntdb.IndexJSON("order")); err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()

	r.Handle("/", appHandler(todo.NewGetAllHandler(db))).Methods("GET", "OPTIONS")
	r.Handle("/", appHandler(todo.NewSaveHandler(db))).Methods("POST")
	r.Handle("/", appHandler(todo.NewDeleteHandler(db))).Methods("DELETE")
	r.Handle("/{key}", appHandler(todo.NewGetOneHandler(db))).Methods("GET", "OPTIONS")
	r.Handle("/{key}", appHandler(todo.NewUpdateHandler(db))).Methods("PATCH")
	r.Handle("/{key}", appHandler(todo.NewDeleteHandler(db))).Methods("DELETE")

	r.Use(mux.CORSMethodMiddleware(r))
	r.Use(responseHeaderSetter)

	log.Fatal(http.ListenAndServe(":"+port, r))
	defer db.Close()
}

type appHandler func(http.ResponseWriter, *http.Request) error

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func responseHeaderSetter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
