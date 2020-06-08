package todo

import (
	"log"
	"net/http"
)

func handleBadRequest(err error, w http.ResponseWriter) {
	msg := "unexpected bad request"
	log.Println(msg, err)
	w.WriteHeader(http.StatusBadRequest)
	if _, err := w.Write([]byte(msg)); err != nil {
		log.Println(err)
	}
}

func handleInternalServerError(err error, w http.ResponseWriter) {
	msg := "unexpected internal server error"
	log.Println(msg, err)
	w.WriteHeader(http.StatusInternalServerError)
	if _, err := w.Write([]byte(msg)); err != nil {
		log.Println(err)
	}
}
