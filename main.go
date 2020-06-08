package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/eguanlao/todo-backend-go-mux-buntdb/todo"
	"github.com/gorilla/mux"
	"github.com/tidwall/buntdb"
)

func main() {
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		log.Println(pair[0], "=", pair[1])
	}

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

	r.HandleFunc("/", todo.NewGetAllHandler(db)).Methods("GET", "OPTIONS")
	r.HandleFunc("/", todo.NewSaveHandler(db)).Methods("POST")
	r.HandleFunc("/", todo.NewDeleteHandler(db)).Methods("DELETE")
	r.HandleFunc("/{key}", todo.NewGetOneHandler(db)).Methods("GET", "OPTIONS")
	r.HandleFunc("/{key}", todo.NewUpdateHandler(db)).Methods("PATCH")
	r.HandleFunc("/{key}", todo.NewDeleteHandler(db)).Methods("DELETE")

	r.Use(mux.CORSMethodMiddleware(r))
	r.Use(responseHeaderSetter)

	log.Fatal(http.ListenAndServe(":"+port, r))
	defer db.Close()
}

func responseHeaderSetter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
