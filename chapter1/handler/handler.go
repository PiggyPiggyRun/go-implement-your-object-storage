package handler

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var PATH = os.Getenv("STORAGE_ROOT") + "/objects/"

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPut {
		put(w, r)
		return
	}
	if r.Method == http.MethodGet {
		get(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func put(w http.ResponseWriter, r *http.Request) {
	f, e := os.Create(PATH + strings.Split(r.URL.EscapedPath(), "/")[2])
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	io.Copy(f, r.Body)
}

func get(w http.ResponseWriter, r *http.Request) {
	f, e := os.Open(PATH + strings.Split(r.URL.EscapedPath(), "/")[2])
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer f.Close()
	io.Copy(w, f)
}
