package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nnqq/eth-parser/pkg/eth"
	"github.com/nnqq/eth-parser/pkg/logger"
	"github.com/nnqq/eth-parser/pkg/parser"
	"net"
	"net/http"
	"strconv"
)

type Server struct {
	logger logger.Printer
	parser parser.Parser
	srv    *http.Server
}

func NewServer(logger logger.Printer, parser parser.Parser, host string, port int) *Server {
	s := &Server{
		logger: logger,
		parser: parser,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/current-block", s.currentBlock)
	mux.HandleFunc("/subscribe", s.subscribe)
	mux.HandleFunc("/transactions", s.transactions)
	mux.HandleFunc("/healthz", s.healthz)

	addr := net.JoinHostPort(host, strconv.Itoa(port))
	s.srv = &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	s.logger.Printf("http server init on %s", addr)

	return s
}

func (s *Server) Run() error {
	err := s.srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("s.srv.ListenAndServe: %w", err)
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) {
	s.logger.Printf("http graceful shutdown...")
	err := s.srv.Shutdown(ctx)
	if err != nil {
		s.logger.Printf("s.srv.Shutdown: %v", err)
	}
	s.logger.Printf("exit")
}

func (s *Server) currentBlock(w http.ResponseWriter, r *http.Request) {
	type res struct {
		Block int `json:"block"`
	}
	s.handle(w, r, "currentBlock", http.MethodGet, res{
		Block: s.parser.GetCurrentBlock(),
	})
}

func (s *Server) healthz(w http.ResponseWriter, r *http.Request) {
	type res struct {
		Success bool `json:"success"`
	}
	s.handle(w, r, "healthz", http.MethodGet, res{
		Success: true,
	})
}

func (s *Server) subscribe(w http.ResponseWriter, r *http.Request) {
	type req struct {
		Address string `json:"address"`
	}
	type res struct {
		Success bool `json:"success"`
	}

	defer func() {
		_ = r.Body.Close()
	}()

	var reqBody req
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		s.logger.Printf("subscribe: json.NewDecoder: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	s.handle(w, r, "subscribe", http.MethodPost, res{
		Success: s.parser.Subscribe(reqBody.Address),
	})
}

func (s *Server) transactions(w http.ResponseWriter, r *http.Request) {
	type res struct {
		Transactions []eth.Transaction `json:"transactions"`
	}

	addr := r.URL.Query().Get("address")
	if addr == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	s.handle(w, r, "transactions", http.MethodGet, res{
		Transactions: s.parser.GetTransactions(addr),
	})
}

func (s *Server) handle(w http.ResponseWriter, r *http.Request, funcName, httpMethod string, res interface{}) {
	setJSON(w)

	if r.Method != httpMethod {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	b, err := json.Marshal(res)
	if err != nil {
		s.logger.Printf(funcName+": json.Marshal: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(b)
	if err != nil {
		s.logger.Printf(funcName+": w.Write: %v", err)
	}
}

func setJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}
