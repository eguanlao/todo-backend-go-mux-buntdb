package todo

import (
	"log"
	"net/http"
)

func handleInternalServerError(err error, w http.ResponseWriter) {
	msg := "unexpected internal server error"
	log.Println(msg, err)
	w.WriteHeader(http.StatusInternalServerError)
	if _, err := w.Write([]byte(msg)); err != nil {
		log.Println(err)
	}
}
