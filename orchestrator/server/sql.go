package server

import (
	"CUBUS-core/shared/types"
	"database/sql"
	"log"
)

func initDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./cube.db")
	if err != nil {
		return nil, err
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS cubes (
  "id" TEXT PRIMARY KEY NOT NULL UNIQUE,
  "cube_type" TEXT NOT NULL,
  "cube_name" TEXT NOT NULL,
  "public_key" TEXT NOT NULL
 );`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func saveCube(db *sql.DB, data types.CubeConfig) error {
	insertCubeSQL := `INSERT INTO cubes (id, cube_type, cube_name, public_key) VALUES (?, ?, ?, ?)`
	statement, err := db.Prepare(insertCubeSQL)
	if err != nil {
		return err
	}
	defer func(statement *sql.Stmt) {
		err := statement.Close()
		if err != nil {
			log.Println("Failed to close statement of the sql db: ", err)
		}
	}(statement)

	var publicKey interface{}
	if data.PublicKey == nil {
		publicKey = ""
	} else {
		publicKey = data.PublicKey
	}
	_, err = statement.Exec(data.Id, data.CubeType.Value, data.CubeName, publicKey)
	if err != nil {
		return err
	}

	return nil
}

func getAllCubes(db *sql.DB) ([]types.CubeConfig, error) {
	rows, err := db.Query("SELECT id, cube_type, cube_name, public_key FROM cubes")
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println("Failed to close rows: ", err)
		}
	}(rows)

	var cubes []types.CubeConfig
	for rows.Next() {
		var cube types.CubeConfig
		var publicKey string
		err = rows.Scan(&cube.Id, &cube.CubeType.Value, &cube.CubeName, &publicKey)
		if err != nil {
			return nil, err
		}
		cube.PublicKey = publicKey
		cubes = append(cubes, cube)
	}

	return cubes, nil
}
