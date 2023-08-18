package queues

import (
	"queue-manager/internal/queues/interfaces"
	"sync"
)

type QueueManager struct {
	queues map[string]interfaces.QueueInterface
	mu     sync.Mutex
}

func (qm *QueueManager) GetQueue(name string) interfaces.QueueInterface {
	qm.mu.Lock()
	defer qm.mu.Unlock()

	if qm.queues == nil {
		qm.queues = make(map[string]interfaces.QueueInterface)
	}

	if _, ok := qm.queues[name]; !ok {
		qm.queues[name] = &FIFOQueue{}
	}

	return qm.queues[name]
}

func (qm *QueueManager) DeleteQueue(name string) {
	qm.mu.Lock()
	defer qm.mu.Unlock()

	delete(qm.queues, name)
}

func (qm *QueueManager) ListQueues() []string {
	qm.mu.Lock()
	defer qm.mu.Unlock()

	var queues []string
	for name := range qm.queues {
		queues = append(queues, name)
	}

	return queues
}

func (qm *QueueManager) ReadAll(name string) []string {
	qm.mu.Lock()
	defer qm.mu.Unlock()

	if _, ok := qm.queues[name]; !ok {
		return []string{}
	}

	return qm.queues[name].ReadAll()
}

func (qm *QueueManager) CreateQueue(name string) {
	qm.mu.Lock()
	defer qm.mu.Unlock()

	if _, ok := qm.queues[name]; !ok {
		qm.queues[name] = &FIFOQueue{}
	}
}

func (qm *QueueManager) GetAllQueues() map[string][]string {
	qm.mu.Lock()
	defer qm.mu.Unlock()

	queuesCopy := make(map[string][]string)
	for key, queue := range qm.queues {
		queuesCopy[key] = queue.ReadAll()
	}
	return queuesCopy
}

func (qm *QueueManager) SetQueue(name string, queue interfaces.QueueInterface) {
	qm.mu.Lock()
	defer qm.mu.Unlock()

	qm.queues[name] = queue
}

func NewQueueManager() interfaces.QueueManagerInterface {
	qm := &QueueManager{}
	qm.queues = make(map[string]interfaces.QueueInterface)
	return qm
}
