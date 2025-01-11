package biz_entity

type IStreamGenerateQueue interface {
	PushErr(err error)
	Push(chunk IQueueEvent)
	Final(chunk IQueueEvent)
	Fork() IStreamGenerateQueue
	Close()
	Listen()
	GetQueues() (chan *MessageQueueMessage, chan *MessageQueueMessage, chan *MessageQueueMessage)
	CloseOutErr()
	CloseOutNormalExit()
	Debug()
}
