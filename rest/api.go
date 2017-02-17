package rest

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

// Ls ...
type Ls struct {
	Childrens []string `json:"childrens"`
	Path      string   `json:"path"`
	State     string   `json:"state"`
	Error     string   `json:"error"`
}

// Get ...
type Get struct {
	Path  string `json:"path"`
	State string `json:"state"`
	Error string `json:"error"`
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
		w.Write([]byte("JSON parsing failure"))
		return
	}
	if childrens.Error != "" {
		w.WriteHeader(500)
		w.Write(data)
		return
	}

	w.WriteHeader(200)
	w.Write(data)
}

// ReadRequestBody ...
func ReadRequestBody(req *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return []byte(""), err
	}

	return body, err
}

// UP ...
func (zk ZooNode) UP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := vars["path"]

	ch := make(chan string)
	// Read request body as []byte
	content, err := ReadRequestBody(r)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	// Create node
	if r.Method == "PUT" {
		go func() { ch <- zk.CreateNode(path, content) }()
	} else if r.Method == "POST" {
		go func() { ch <- zk.UpdateNode(path, content) }()
	} else {
		e := strings.Join([]string{r.Method, "not alowed"}, " ")
		w.WriteHeader(500)
		w.Write([]byte(e))
		return
	}
	defer r.Body.Close()

	state := <-ch

	w.WriteHeader(200)
	w.Write([]byte(state))
}

// RM ...
func (zk ZooNode) RM(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		e := strings.Join([]string{r.Method, "not alowed"}, " ")
		w.WriteHeader(500)
		w.Write([]byte(e))
		return
	}
	vars := mux.Vars(r)
	path := vars["path"]

	go func() { zk.RMR(path) }()

	w.WriteHeader(200)
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
		w.Write([]byte("JSON parsing failure"))
		return
	}

	if childrens.Error != "" {
		w.WriteHeader(500)
		w.Write(data)
		return
	}

	w.WriteHeader(200)
	w.Write(data)
}

// Serve ...
func Serve(listen string, zk ZooNode) {
	r := mux.NewRouter()

	r.HandleFunc("/v1/ls{path:[a-z0-9-_/.:]+}", zk.LS)
	r.HandleFunc("/v1/get{path:[a-z0-9-_/.:]+}", zk.GET)
	r.HandleFunc("/v1/rmr{path:[a-z0-9-_/.:]+}", zk.RM)
	r.HandleFunc("/v1/up{path:[a-z0-9-_/.:]+}", zk.UP)

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
