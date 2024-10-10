package cubes

import "CUBUS-core/shared/types"

type Base interface {
	InitMessageChannel()
	SendMessage(message types.Message)
	GetMessageChannel() chan types.Message
	GetConfig() types.CubeConfig
}

type BaseStruct struct {
	messageChannel chan types.Message
}

func (b *BaseStruct) InitMessageChannel() {
	b.messageChannel = make(chan types.Message)
}

func (b *BaseStruct) SendMessage(message types.Message) {
	b.messageChannel <- message
}

func (b *BaseStruct) GetMessageChannel() chan types.Message {
	return b.messageChannel
}

func (b *BaseStruct) GetConfig() types.CubeConfig {
	return types.CubeConfig{}
}
