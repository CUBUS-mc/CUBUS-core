package client

import (
	"CUBUS-core/orchestrator"
	"CUBUS-core/shared/types"
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	client orchestrator.OrchestratorClient
	conn   *grpc.ClientConn
}

func NewClient(serverAddr string) (*Client, error) {
	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("Failed to connect to the server: ", err)
		return nil, err
	}
	return &Client{
		client: orchestrator.NewOrchestratorClient(conn),
		conn:   conn,
	}, nil
}

func (c *Client) Close() {
	err := c.conn.Close()
	if err != nil {
		fmt.Println("Failed to close connection: ", err)
	}
}

func (c *Client) CreateCube(ctx context.Context, config types.CubeConfig) (*orchestrator.CreateCubeResponse, error) {
	publicKeyStr, err := convertPublicKeyToString(&config.PublicKey)
	if err != nil {
		return nil, err
	}
	req := &orchestrator.CreateCubeRequest{Config: &orchestrator.CubeConfig{
		Id:        config.Id,
		Name:      config.CubeName,
		PublicKey: publicKeyStr,
		QueueServer: &orchestrator.QueueServerConfig{
			Url:      config.QueueServer.Url,
			Username: config.QueueServer.Username,
			Password: config.QueueServer.Password,
			Db:       int32(config.QueueServer.DB),
		},
	}}
	return c.client.CreateCube(ctx, req)
}

func (c *Client) GetAllCubes(ctx context.Context) ([]types.CubeConfig, error) {
	resp, err := c.client.GetAllCubes(ctx, &orchestrator.GetAllCubesRequest{})
	if err != nil {
		return nil, err
	}
	cubes := make([]types.CubeConfig, 0)
	for _, cube := range resp.Cubes {
		publicKey, err := convertStringToPublicKey(cube.PublicKey)
		if err != nil {
			return nil, err
		}
		cubes = append(cubes, types.CubeConfig{
			Id:        cube.Id,
			CubeName:  cube.Name,
			PublicKey: *publicKey,
			QueueServer: types.QueueServerConfig{
				Url:      cube.QueueServer.Url,
				Username: cube.QueueServer.Username,
				Password: cube.QueueServer.Password,
				DB:       int(cube.QueueServer.Db),
			},
		})
	}
	return cubes, nil
}

func convertPublicKeyToString(publicKey *rsa.PublicKey) (string, error) {
	if publicKey == nil || publicKey.N == nil || publicKey.E == 0 {
		return "", nil
	}
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(pubKeyBytes), nil
}

func convertStringToPublicKey(publicKeyStr string) (*rsa.PublicKey, error) {
	if publicKeyStr == "" {
		return &rsa.PublicKey{}, nil
	}
	pubKeyBytes, err := base64.StdEncoding.DecodeString(publicKeyStr)
	if err != nil {
		return nil, err
	}
	pubKey, err := x509.ParsePKIXPublicKey(pubKeyBytes)
	if err != nil {
		return nil, err
	}
	return pubKey.(*rsa.PublicKey), nil
}
