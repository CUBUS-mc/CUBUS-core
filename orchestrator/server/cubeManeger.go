package server

import (
	"CUBUS-core/cubes"
	"CUBUS-core/orchestrator"
	"crypto/rsa"
	"database/sql"
	"log"
)

type CubeManager struct {
	db *sql.DB
}

func NewCubeManager(db *sql.DB) *CubeManager {
	cm := CubeManager{db: db}
	return &cm
}

func (cm *CubeManager) StartCube(cube *orchestrator.CubeConfig) {
	log.Println("Starting cube: ", cube.Id, " with the name ", cube.Name)
	/*
		q := queen.New(types.QueenConfig{
			CubeConfig:    *cube,
			RedisAddress:  "localhost:6379",
			RedisPassword: "",
			RedisDB:       0,
			Tasks:         tasks.Tasks,
		})
		go cm.Listen(q)
		q.Ping()
	*/
}

func (cm *CubeManager) Listen(cube cubes.Base) {
	println("Listening for messages from cube: ", cube.GetConfig().Id)
	for message := range cube.GetMessageChannel() {
		println("Received message: ", message.MessageType)
		switch message.MessageType {
		case "UPDATE PUBLIC KEY":
			err := updatePublicKey(cm.db, cube.GetConfig().Id, message.Message.(*rsa.PublicKey))
			if err != nil {
				println("Failed to update public key: ", err)
				continue
			}
		}
	}
}
