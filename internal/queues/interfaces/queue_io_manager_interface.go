package interfaces

type QueueIOManagerInterface interface {
	SaveQueuesToFile(qm QueueManagerInterface, fileName string) error
	LoadQueuesFromFile(qm QueueManagerInterface, fileName string) error
}
