package interfaces

type QueueManagerInterface interface {
	GetQueue(name string) QueueInterface
	DeleteQueue(name string)
	ListQueues() []string
	ReadAll(name string) []string
	CreateQueue(name string)
	GetAllQueues() map[string][]string
	SetQueue(name string, queue QueueInterface)
}
