package rest

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

// Ls ...
type Ls struct {
	Childrens []string `json:"childrens"`
	Path      string   `json:"path"`
	State     string   `json:"state"`
	Error     error    `json:"error"`
}

// Get ...
type Get struct {
	Path  string `json:"path"`
	State string `json:"state"`
	Error error  `json:"error"`
	Data  []byte `json:"data"`
}

// LS ...
func (zk ZooNode) LS(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := vars["path"]

	ch := make(chan Ls)

	go func() { ch <- zk.GetChildrens(path) }()

	childrens := <-ch

	data, err := json.Marshal(childrens)
	if err != nil {
		w.WriteHeader(500)
	}
	w.WriteHeader(200)
	w.Write(data)
}

// GET ...
func (zk ZooNode) GET(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := vars["path"]

	ch := make(chan Get)

	go func() { ch <- zk.GetNode(path) }()

	childrens := <-ch

	data, err := json.Marshal(childrens)
	if err != nil {
		w.WriteHeader(500)
	}
	w.WriteHeader(200)
	w.Write(data)
}

// Serve ...
func Serve(listen string, zk ZooNode) {
	r := mux.NewRouter()

	r.HandleFunc("/v1/ls{path:[a-z0-9-_/.:]+}", zk.LS)
	r.HandleFunc("/v1/get{path:[a-z0-9-_/.:]+}", zk.GET)

	http.Handle("/", r)

	srv := http.Server{
		Handler:      r,
		Addr:         listen,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Print("Listening API on ", listen)
	log.Fatal(srv.ListenAndServe())
}
