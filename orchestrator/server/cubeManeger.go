package server

import (
	"CUBUS-core/cubes"
	"CUBUS-core/cubes/queen"
	"CUBUS-core/shared/tasks"
	"CUBUS-core/shared/types"
	"database/sql"
)

type CubeManager struct {
	cubes []cubes.Base
	db    *sql.DB
}

func NewCubeManager(db *sql.DB) *CubeManager {
	cm := CubeManager{db: db}
	cm.cubes = make([]cubes.Base, 0)
	go cm.Listen()
	return &cm
}

func (cm *CubeManager) StartCube(cube *types.CubeConfig) {
	switch cube.CubeType {
	case types.CubeTypes.GenericWorker:
	case types.CubeTypes.Queen:
		cm.cubes = append(cm.cubes, queen.New(types.QueenConfig{
			CubeConfig:    *cube,
			RedisAddress:  "localhost:6379",
			RedisPassword: "",
			RedisDB:       0,
			Tasks:         tasks.Tasks,
		}))
	default:
	}
}

// FIXME: Fix this function so it updates the public key in the database when a cube sends a message to update its public key
func (cm *CubeManager) Listen() {
	for {
		for _, cube := range cm.cubes {
			go func(cube cubes.Base) {
				for message := range cube.GetMessageChannel() {
					println("Received message: ", message.MessageType)
					switch message.MessageType {
					case "UPDATE PUBLIC KEY":
						err := updatePublicKey(cm.db, cube.GetConfig().Id, message.Message.(types.CubeConfig).PublicKey)
						if err != nil {
							println("Failed to update public key: ", err)
							return
						}
					}
				}
			}(cube)
		}
	}
}
