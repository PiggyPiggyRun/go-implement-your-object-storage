package main

import (
	"./handler"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/objects/", handler.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("PORT"), nil))
}
