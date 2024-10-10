package tasks

import (
	"CUBUS-core/shared/types"
	"context"
	"encoding/json"
	"github.com/hibiken/asynq"
)

const (
	TypePing = "ping"
)

type PingPayload struct {
	Message string
}

func NewPingTask(message string) (*asynq.Task, error) {
	payload, err := json.Marshal(PingPayload{Message: message})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypePing, payload), nil
}

func HandlePingTask(_ context.Context, t *asynq.Task) error {
	var payload PingPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}
	_, err := t.ResultWriter().Write([]byte(payload.Message))
	if err != nil {
		return err
	}
	return nil
}

var Tasks = []types.Task{
	{
		Type:    TypePing,
		Handler: HandlePingTask,
	},
}
