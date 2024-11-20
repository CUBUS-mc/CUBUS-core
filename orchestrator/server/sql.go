package server

import (
	"CUBUS-core/orchestrator"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/base64"
	"errors"
	"log"
)

func initDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./cube.db")
	if err != nil {
		return nil, err
	}

	createTableSQL := ` CREATE TABLE IF NOT EXISTS queue_servers (
	  "id" INTEGER PRIMARY KEY AUTOINCREMENT,
	  "url" TEXT NOT NULL,
	  "username" TEXT NOT NULL,
	  "password" TEXT NOT NULL,
	  "db" INTEGER NOT NULL
   );
	`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}

	createTableSQL = `CREATE TABLE IF NOT EXISTS cubes (
	  "id" TEXT PRIMARY KEY NOT NULL UNIQUE,
	  "cube_name" TEXT NOT NULL,
	  "public_key" TEXT,
	 	  "queue_server_id" INTEGER,
		  FOREIGN KEY (queue_server_id) REFERENCES queue_servers(id)
	 );`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func saveCube(db *sql.DB, data *orchestrator.CubeConfig) error {
	var queueServerID int

	checkQueueServerSQL := `SELECT id FROM queue_servers WHERE url = ? AND username = ? AND password = ? AND db = ?`
	err := db.QueryRow(checkQueueServerSQL, data.QueueServer.Url, data.QueueServer.Username, data.QueueServer.Password, data.QueueServer.Db).Scan(&queueServerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			insertQueueServerSQL := `INSERT INTO queue_servers (url, username, password, db) VALUES (?, ?, ?, ?)`
			result, err := db.Exec(insertQueueServerSQL, data.QueueServer.Url, data.QueueServer.Username, data.QueueServer.Password, data.QueueServer.Db)
			if err != nil {
				return err
			}
			queueServerID64, err := result.LastInsertId()
			if err != nil {
				return err
			}
			queueServerID = int(queueServerID64)
		} else {
			return err
		}
	}

	insertCubeSQL := `INSERT INTO cubes (id, cube_name, public_key, queue_server_id) VALUES (?, ?, ?, ?)`
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

	_, err = statement.Exec(data.Id, data.Name, data.PublicKey, queueServerID)
	if err != nil {
		return err
	}

	return nil
}

func getAllCubes(db *sql.DB) ([]*orchestrator.CubeConfig, error) {
	rows, err := db.Query("SELECT id, cube_name, public_key, queue_server_id FROM cubes")
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println("Failed to close rows: ", err)
		}
	}(rows)

	var cubes []*orchestrator.CubeConfig
	for rows.Next() {
		var cube orchestrator.CubeConfig
		var serverID int
		err = rows.Scan(&cube.Id, &cube.Name, &cube.PublicKey, &serverID)
		if err != nil {
			return nil, err
		}

		var queueServer orchestrator.QueueServerConfig
		err = db.QueryRow("SELECT url, username, password, db FROM queue_servers WHERE id = ?", serverID).Scan(&queueServer.Url, &queueServer.Username, &queueServer.Password, &queueServer.Db)
		if err != nil {
			return nil, err
		}
		cube.QueueServer = &queueServer
		cubes = append(cubes, &cube)
	}

	return cubes, nil
}

func updatePublicKey(db *sql.DB, id string, publicKey *rsa.PublicKey) error {
	updateSQL := `UPDATE cubes SET public_key = ? WHERE id = ?`
	statement, err := db.Prepare(updateSQL)
	if err != nil {
		log.Printf("Failed to prepare statement: %v", err)
		return err
	}
	defer func(statement *sql.Stmt) {
		err := statement.Close()
		if err != nil {
			log.Println("Failed to close statement of the sql db: ", err)
		}
	}(statement)

	publicKeyStr, err := convertPublicKeyToString(publicKey)
	if err != nil {
		log.Printf("Failed to convert public key to string: %v", err)
		return err
	}

	println("Updating public key to: ", publicKeyStr, " for cube with id: ", id)
	_, err = statement.Exec(publicKeyStr, id)
	if err != nil {
		log.Printf("Failed to execute statement: %v", err)
		return err
	}

	return nil
}

func convertPublicKeyToString(publicKey *rsa.PublicKey) (string, error) {
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(pubKeyBytes), nil
}
