package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Server struct {
	port string
}

func NewServer(port string) *Server {
	return &Server{
		port: port,
	}
}

func (s *Server) createHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		fmt.Println("Received request on /create endpoint")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {

			}
		}(r.Body)

		var prettyJSON bytes.Buffer
		err = json.Indent(&prettyJSON, body, "", "  ")
		if err != nil {
			http.Error(w, "Failed to format JSON", http.StatusInternalServerError)
			return
		}

		fmt.Println("Request body: \n", prettyJSON.String())
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte("Request received"))
		if err != nil {
			return
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, err := w.Write([]byte("Only POST method is allowed"))
		if err != nil {
			return
		}
	}
}

func (s *Server) Start() {
	go func() {
		http.HandleFunc("/create", s.createHandler)
		err := http.ListenAndServe(s.port, nil)
		if err != nil {
			fmt.Println("Failed to start server: ", err)
		} else {
			fmt.Println("Server started on port", s.port)
		}
	}()
}
