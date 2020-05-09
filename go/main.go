package main

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

type response struct {
	ID   string
	Name string
	Time int64
}

var client *http.Client
var url string

func init() {
	tr := &http.Transport{
		MaxIdleConns:        4000,
		MaxIdleConnsPerHost: 4000,
	}
	client = &http.Client{Transport: tr}

	url = "http://" + os.Getenv("HOST") + ":5002/data"
}

func main() {

	router := httprouter.New()
	router.GET("/", hello)
	router.GET("/test", test)

	addr := ":5001"
	fmt.Println("listening on " + addr)
	log.Fatal(http.ListenAndServe(addr, router))
}

func test(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	rsp, err := client.Get(url)
	if err != nil {
		serverError(w, err.Error())
		return
	}

	defer rsp.Body.Close()

	// deserialize
	obj := response{}
	err = json.NewDecoder(rsp.Body).Decode(&obj)
	if err != nil {
		serverError(w, err.Error())
		return
	}

	// serialize
	jsonStr, err := json.Marshal(&obj)
	if err != nil {
		serverError(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonStr); err != nil {
		serverError(w, err.Error())
		return
	}
}

func hello(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func serverError(w http.ResponseWriter, msg string) {
	http.Error(w, msg, http.StatusInternalServerError)
}
