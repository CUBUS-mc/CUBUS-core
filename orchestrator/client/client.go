package client

import (
	"CUBUS-core/shared/types"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) CreateNewCube(serverUrl string, config types.CubeConfig) error {
	configAsJson, err := config.ToJson()
	if err != nil {
		return err
	}
	resp, err := http.Post(serverUrl+"/create", "application/json", bytes.NewBuffer(configAsJson))
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Failed to close response body: ", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("server error: %s", string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("Response: ", string(body))
	return nil
}

func (c *Client) GetAllCubes(serverUrl string) ([]types.CubeConfig, error) {
	resp, err := http.Get(serverUrl + "/cubes")
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Failed to close response body: ", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("server error: %s", string(body))
	}

	var cubes []types.CubeConfig
	err = json.NewDecoder(resp.Body).Decode(&cubes)
	if err != nil {
		return nil, err
	}

	return cubes, nil
}
