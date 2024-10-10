package server

import (
	"CUBUS-core/shared/types"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"net/http"
)

type Server struct {
	port string
	db   *sql.DB
	cm   *CubeManager
}

func NewServer(port string) *Server {
	db, err := initDB()
	if err != nil {
		fmt.Println("Failed to initialize database: ", err)
		return nil
	}

	server := Server{
		port: port,
		db:   db,
		cm:   NewCubeManager(db),
	}

	go server.startCubes()

	return &server
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
				fmt.Println("Failed to close request body: ", err)
			}
		}(r.Body)

		var cubeConfig types.CubeConfig
		err = json.Unmarshal(body, &cubeConfig)
		if err != nil {
			http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
			return
		}

		err = saveCube(s.db, cubeConfig)
		if err != nil {
			http.Error(w, "Failed to save cube", http.StatusInternalServerError)
			return
		}

		go s.cm.StartCube(&cubeConfig)

		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte("Request received and cube saved"))
		if err != nil {
			fmt.Println("Failed to write response: ", err)
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, err := w.Write([]byte("Only POST method is allowed"))
		if err != nil {
			fmt.Println("Failed to write response: ", err)
		}
	}
}

func (s *Server) getAllCubesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		fmt.Println("Received request on /cubes endpoint")

		cubes, err := getAllCubes(s.db)
		if err != nil {
			http.Error(w, "Failed to get cubes", http.StatusInternalServerError)
			return
		}

		response, err := json.Marshal(cubes)
		if err != nil {
			http.Error(w, "Failed to parse JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(response)
		if err != nil {
			fmt.Println("Failed to write response: ", err)
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, err := w.Write([]byte("Only GET method is allowed"))
		if err != nil {
			fmt.Println("Failed to write response: ", err)
		}
	}
}

func (s *Server) Start() {
	go func() {
		http.HandleFunc("/create", s.createHandler)
		http.HandleFunc("/cubes", s.getAllCubesHandler)
		err := http.ListenAndServe(s.port, nil)
		if err != nil {
			fmt.Println("Failed to start server: ", err)
		} else {
			fmt.Println("Server started on port", s.port)
		}
	}()
}

func (s *Server) startCubes() {
	cubes, err := getAllCubes(s.db)
	if err != nil {
		fmt.Println("Failed to get cubes: ", err)
		return
	}

	for _, cube := range cubes {
		go s.cm.StartCube(&cube)
	}
}
