package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func NewHTTPServer(addr string) *http.Server {
	srv := NewHTTPLogServer()
	r := mux.NewRouter()

	r.HandleFunc("/", srv.handleProduce).Methods("POST")
	r.HandleFunc("/", srv.handleConsume).Methods("GET")

	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}

type IHttpLogServer interface {
	Log() *Log
}

type httpLogServer struct {
	log *Log
}

func NewHTTPLogServer() *httpLogServer {
	return &httpLogServer{
		log: NewLog(),
	}
}

func (s *httpLogServer) Log() *Log {
	return s.log
}

type ProduceRequest struct {
	Record Record `json:"record"`
}

type ProduceResponse struct {
	Offset uint64 `json:"offset"`
}

type ConsumeRequest struct {
	Offset uint64 `json:"offset"`
}

type ConsumeResponse struct {
	Record Record `json:"record"`
}

func (s *httpLogServer) handleProduce(w http.ResponseWriter, r *http.Request) {
	var req ProduceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Println("error here")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	off, err := s.Log().Append(req.Record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Printf("offset: %d\n", off)
	fmt.Printf("record: %v\n", req.Record)

	res := ProduceResponse{Offset: off}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *httpLogServer) handleConsume(w http.ResponseWriter, r *http.Request) {
	var req ConsumeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	record, err := s.Log().Read(req.Offset)
	fmt.Printf("record: %v\n", record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	res := ConsumeResponse{Record: record}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
