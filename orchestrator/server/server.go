package server

import (
	"CUBUS-core/orchestrator"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
)

type Server struct {
	orchestrator.UnimplementedOrchestratorServer
	port        string
	db          *sql.DB
	cubeManager *CubeManager
}

func NewServer(port string) *Server {
	db, err := initDB()
	if err != nil {
		fmt.Println("Failed to initialize database: ", err)
		return nil
	}

	server := Server{
		port:        port,
		db:          db,
		cubeManager: NewCubeManager(db),
	}

	return &server
}

func (s *Server) CreateCube(_ context.Context, req *orchestrator.CreateCubeRequest) (*orchestrator.CreateCubeResponse, error) {
	println("Received request to create cube: ", req.GetConfig().Id)
	cubeConfig := req.GetConfig()
	err := saveCube(s.db, cubeConfig)
	if err != nil {
		fmt.Println("Failed to save cube: ", err)
		return nil, status.Error(codes.Internal, "Failed to save cube")
	}

	go s.cubeManager.StartCube(cubeConfig)

	return &orchestrator.CreateCubeResponse{Id: cubeConfig.Id}, nil
}

func (s *Server) GetAllCubes(_ context.Context, _ *orchestrator.GetAllCubesRequest) (*orchestrator.GetAllCubesResponse, error) {
	println("Received request to get all cubes")
	cubes, err := getAllCubes(s.db)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to get cubes")
	}

	return &orchestrator.GetAllCubesResponse{Cubes: cubes}, nil

}

func (s *Server) Start() {
	go func() {
		lis, err := net.Listen("tcp", s.port)
		if err != nil {
			fmt.Println("Failed to start server: ", err)
			return
		}

		grpcServer := grpc.NewServer()
		orchestrator.RegisterOrchestratorServer(grpcServer, s)

		fmt.Println("Server started on port", s.port)
		if err := grpcServer.Serve(lis); err != nil {
			fmt.Println("Failed to serve: ", err)
		}
	}()
	go s.startCubes()
}

func (s *Server) startCubes() {
	cubes, err := getAllCubes(s.db)
	if err != nil {
		fmt.Println("Failed to get cubes: ", err)
		return
	}

	for _, cube := range cubes {
		println("Starting cube: ", cube.Id)
		go s.cubeManager.StartCube(cube)
	}
}
