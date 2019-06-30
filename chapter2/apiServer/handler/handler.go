package handler

import (
	"../heartbeat"
	"../locate"
	"fmt"
	"go-implement-your-object-storage/src/lib/objectstream"
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

//func put(w http.ResponseWriter, r *http.Request) {
//	f, e := os.Create(PATH + strings.Split(r.URL.EscapedPath(), "/")[2])
//	if e != nil {
//		log.Println(e)
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//	defer f.Close()
//	io.Copy(f, r.Body)
//}

func put(w http.ResponseWriter, r *http.Request) {
	object := strings.Split(r.URL.EscapedPath(), "/")[2]
	c, e := storeObject(r.Body, object)
	if e != nil {
		log.Println(e)
	}
	w.WriteHeader(c)
}

func storeObject(r io.Reader, object string) (int, error) {
	stream, e := putStream(object)
	if e != nil {
		return http.StatusServiceUnavailable, e
	}

	io.Copy(stream, r)
	e = stream.Close()
	if e != nil {
		return http.StatusInternalServerError, e
	}
	return http.StatusOK, nil
}

func putStream(object string) (*objectstream.PutStream, error) {
	server := heartbeat.ChooseRandomDataServer()
	if server == "" {
		return nil, fmt.Errorf("cannot find any dataServer")
	}

	return objectstream.NewPutStream(server, object), nil
}

//func get(w http.ResponseWriter, r *http.Request) {
//	f, e := os.Open(PATH + strings.Split(r.URL.EscapedPath(), "/")[2])
//	if e != nil {
//		log.Println(e)
//		w.WriteHeader(http.StatusNotFound)
//		return
//	}
//	defer f.Close()
//	io.Copy(w, f)
//}

func get(w http.ResponseWriter, r *http.Request) {
	object := strings.Split(r.URL.EscapedPath(), "/")[2]
	stream, e := getStream(object)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	io.Copy(w, stream)
}

func getStream(object string) (io.Reader, error) {
	server := locate.Locate(object)
	if server == "" {
		return nil, fmt.Errorf("object %s locate fail", object)
	}
	return objectstream.NewGetStream(server, object)
}
