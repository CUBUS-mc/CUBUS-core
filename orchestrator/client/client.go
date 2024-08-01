package client

import (
	"CUBUS-core/shared/types"
	"bytes"
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("Response: ", string(body))
	return nil
}
