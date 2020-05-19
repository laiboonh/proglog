package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

func getRootHandler(log *Log) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		switch method := request.Method; method {
		case http.MethodGet:
			var req ConsumeRequest
			err := json.NewDecoder(request.Body).Decode(&req)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusBadRequest)
				return
			}
			record, err := log.Read(req.Offset)
			if err == ErrOffsetNotFound {
				http.Error(writer, err.Error(), http.StatusNotFound)
				return
			}
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
			res := ConsumeResponse{Record: record}
			err = json.NewEncoder(writer).Encode(res)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}

		case http.MethodPost:
			var req ProduceRequest
			err := json.NewDecoder(request.Body).Decode(&req)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusBadRequest)
				return
			}
			offset, err := log.Append(req.Record)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
			response := ProduceResponse{Offset: offset}
			err = json.NewEncoder(writer).Encode(response)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}

		default:
			http.Error(writer, fmt.Sprintf("unexpected request method %s", method), http.StatusBadRequest)
		}
	}
}

func NewHttpServer(addr string) *http.Server {
	log := NewLog()
	mux := http.NewServeMux()
	mux.HandleFunc("/", getRootHandler(log))
	return &http.Server{
		Addr:    addr,
		Handler: mux,
	}
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

type Log struct {
	mu      sync.Mutex
	records []Record
}

func NewLog() *Log {
	return &Log{}
}

func (l *Log) Append(record Record) (uint64, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	record.Offset = uint64(len(l.records))
	l.records = append(l.records, record)
	return record.Offset, nil
}

func (l *Log) Read(offset uint64) (Record, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if offset >= uint64(len(l.records)) {
		return Record{}, ErrOffsetNotFound
	}
	return l.records[offset], nil
}

var ErrOffsetNotFound = fmt.Errorf("offset not found")

type Record struct {
	Value  []byte `json:"value"`
	Offset uint64 `json:"offset"`
}
