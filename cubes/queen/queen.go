package queen

import (
	"CUBUS-core/cubes"
	"CUBUS-core/shared/tasks"
	"CUBUS-core/shared/types"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/hibiken/asynq"
	"log"
	"os"
	"time"
)

type Queen struct {
	cubes.BaseStruct
	config      types.QueenConfig
	privateKey  crypto.PrivateKey
	asynqClient *asynq.Client
	logger      *log.Logger
}

func New(config types.QueenConfig) *Queen {
	logger := log.New(os.Stdout, "Cube (type: "+config.CubeType.Value+", id: "+config.Id+"):  ", log.LstdFlags)
	q := Queen{config: config, logger: logger}
	q.InitMessageChannel()
	if q.config.PublicKey == nil || q.config.PublicKey == "" {
		q.generateKeyPair()
	} else {
		q.loadPrivateKey()
	}
	go q.start()
	return &q
}

func (q *Queen) generateKeyPair() {
	q.logger.Println("Generating new key pair")
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		q.logger.Fatalln("Failed to generate key pair: ", err)
	}
	q.logger.Println("Key pair generated")
	q.privateKey = privateKey
	q.config.PublicKey = &privateKey.PublicKey

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})
	err = os.WriteFile(q.config.Id+".pem", privateKeyPEM, 0600)
	if err != nil {
		q.logger.Fatalln("Failed to write private key to file: ", err)
	}
	go q.SendMessage(types.Message{
		MessageType: "UPDATE PUBLIC KEY",
		Message:     q.config.PublicKey,
	})
	q.logger.Println("Key pair saved to file")
}

func (q *Queen) loadPrivateKey() {
	q.logger.Println("Loading private key")
	privateKeyPEM, err := os.ReadFile(q.config.Id + ".pem")
	if err != nil {
		q.logger.Fatalln("Failed to read private key from file: ", err)
	}

	block, _ := pem.Decode(privateKeyPEM)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		q.logger.Fatalln("Failed to decode private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		q.logger.Fatalln("Failed to parse private key: ", err)
	}
	q.privateKey = privateKey
	q.logger.Println("Private key loaded")
}

func (q *Queen) startServer() {
	q.logger.Println("Starting asynq server")
	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     q.config.RedisAddress,
			Password: q.config.RedisPassword,
			DB:       q.config.RedisDB,
		},
		asynq.Config{
			Concurrency: 1<<31 - 1,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	mux := asynq.NewServeMux()

	q.logger.Println("Registering tasks")
	for _, task := range q.config.Tasks {
		mux.HandleFunc(task.Type, task.Handler)
		q.logger.Println("Registered task: ", task.Type)
	}

	go func() {
		err := srv.Run(mux)
		if err != nil {

		}
	}()
}

func (q *Queen) startClient() {
	q.logger.Println("Starting asynq client")
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: q.config.RedisAddress, Password: q.config.RedisPassword, DB: q.config.RedisDB})
	defer func(client *asynq.Client) {
		err := client.Close()
		if err != nil {
			q.logger.Println("Failed to close asynq client: ", err)
		}
	}(client)
	q.asynqClient = client
	q.logger.Println("Asynq client started")
}

func (q *Queen) start() {
	q.startServer()
	time.Sleep(1 * time.Second)
	q.startClient()
}

func (q *Queen) Enqueue(task *asynq.Task, options ...asynq.Option) (*asynq.TaskInfo, error) {
	for q.asynqClient == nil {
	}
	return q.asynqClient.Enqueue(task, options...)
}

func (q *Queen) GetConfig() types.CubeConfig {
	return q.config.CubeConfig
}

func (q *Queen) Ping() {
	task, _ := tasks.NewPingTask(q.config.Id)
	info, err := q.Enqueue(task)
	if err != nil {
		q.logger.Println("Failed to enqueue ping task: ", err)
		return
	}
	q.logger.Println("Ping task enqueued: ", info.ID)
}
