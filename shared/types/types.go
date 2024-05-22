package types

import "crypto"

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
	PublicKey crypto.PublicKey
}
