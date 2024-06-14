package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/vmanilo/eth-parser/parser"
)

type server struct {
	parser parser.Parser
}

func NewServer(parser parser.Parser) *server {
	return &server{
		parser: parser,
	}
}

func (s *server) Serve(ctx context.Context, port string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/current-block", s.GetCurrentBlock)
	mux.HandleFunc("POST /api/subscribe", s.Subscribe)
	mux.HandleFunc("GET /api/transactions/{address}", s.GetTransactions)

	srv := http.Server{
		Addr:    port,
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		srv.Shutdown(ctx)
	}()

	return srv.ListenAndServe()
}

type GetCurrentBlockResponse struct {
	BlockNumber int `json:"blockNumber"`
}

func (s *server) GetCurrentBlock(w http.ResponseWriter, _ *http.Request) {
	var resp GetCurrentBlockResponse

	resp.BlockNumber = s.parser.GetCurrentBlock()

	sendJSON(w, resp, nil)
}

type SubscribeRequest struct {
	Address string `json:"address"`
}

type SubscribeResponse struct {
	Ok bool `json:"ok"`
}

func (s *server) Subscribe(w http.ResponseWriter, r *http.Request) {
	var req SubscribeRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var resp SubscribeResponse
	resp.Ok = s.parser.Subscribe(req.Address)

	sendJSON(w, resp, func(w http.ResponseWriter) {
		if resp.Ok {
			w.WriteHeader(http.StatusCreated)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	})
}

type GetTransactionsResponse struct {
	Transactions []parser.Transaction `json:"transactions"`
}

func (s *server) GetTransactions(w http.ResponseWriter, r *http.Request) {
	address := r.PathValue("address")

	var resp GetTransactionsResponse
	resp.Transactions = s.parser.GetTransactions(address)

	sendJSON(w, resp, nil)
}

func sendJSON(w http.ResponseWriter, resp any, addHeader func(w http.ResponseWriter)) {
	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if addHeader != nil {
		addHeader(w)
	}

	w.Write(data)
}
