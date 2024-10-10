package types

import (
	"context"
	"crypto"
	"encoding/json"
	"github.com/hibiken/asynq"
)

type CubeType struct {
	Value string
}

type CubeTypesStruct struct {
	Queen         CubeType
	Security      CubeType
	Database      CubeType
	Api           CubeType
	CubusMod      CubeType
	DiscordBot    CubeType
	Web           CubeType
	Drone         CubeType
	GenericWorker CubeType
}

var CubeTypes = CubeTypesStruct{
	Queen:         CubeType{Value: "queen"},
	Security:      CubeType{Value: "security"},
	Database:      CubeType{Value: "database"},
	Api:           CubeType{Value: "api"},
	CubusMod:      CubeType{Value: "cubus-mod"},
	DiscordBot:    CubeType{Value: "discord-bot"},
	Web:           CubeType{Value: "web"},
	Drone:         CubeType{Value: "drone"},
	GenericWorker: CubeType{Value: "generic-worker"},
}

type CubeConfig struct {
	Id        string
	CubeType  CubeType
	CubeName  string
	PublicKey crypto.PublicKey
}

func (cc *CubeConfig) ToJson() ([]byte, error) {
	return json.Marshal(cc)
}

type Message struct {
	MessageType string
	Message     interface{}
}

type QueenConfig struct {
	CubeConfig
	RedisAddress  string
	RedisPassword string
	RedisDB       int
	Tasks         []Task
}

type Task struct {
	Type    string
	Handler func(ctx context.Context, t *asynq.Task) error
}

type CustomError struct {
	Message string
}

func (e *CustomError) Error() string {
	return e.Message
}
