package controllers

import (
	"encoding/json"
	"net/http"
	"queue-manager/internal/providers"
	"queue-manager/internal/queues/interfaces"
	"queue-manager/internal/structures"
	"queue-manager/internal/tcp"
	"sort"
	"time"
)

//swagger:model CreateQueuePayload
type CreateQueuePayload struct {
	Name string `json:"name"`
}

//swagger:model CreateMessagePayload
type CreateMessagePayload struct {
	QueueName string      `json:"queue_name"`
	Message   interface{} `json:"message"`
}

type ApiController struct {
	qm        interfaces.QueueManagerInterface
	TcpServer *tcp.TcpServer
	logger    providers.Logger
}

func (ac *ApiController) List(w http.ResponseWriter, r *http.Request) {
	queuesList := ac.qm.ListQueues()
	gson, e := json.Marshal(queuesList)
	if e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(200)
	w.Write(gson)

}

func (ac *ApiController) DetailedQueuesList(w http.ResponseWriter, r *http.Request) {
	queuesList := ac.qm.ListQueues()
	sort.Slice(queuesList, func(i, j int) bool {
		return queuesList[i] < queuesList[j]
	})
	queues := make([]structures.DetailedQueuesDashboardResponse, len(queuesList))
	for k, queue := range queuesList {
		messages := ac.qm.GetQueue(queue).ReadAll()
		queues[k] = structures.DetailedQueuesDashboardResponse{
			QueueName:     queue,
			QueueSize:     len(messages),
			QueueMessages: messages,
		}
	}

	gson, e := json.Marshal(queues)
	if e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(200)
	w.Write(gson)
}

func (ac *ApiController) ListMessagesInQueue(w http.ResponseWriter, r *http.Request) {
	queue := ac.qm.GetQueue(r.URL.Query().Get("queue_name"))
	gson, e := json.Marshal(queue.ReadAll())
	if e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(200)
	w.Write(gson)

}

func (ac *ApiController) CreateQueue(w http.ResponseWriter, r *http.Request) {
	var payload CreateQueuePayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ac.qm.CreateQueue(payload.Name)
	w.WriteHeader(201)
}

func (ac *ApiController) CreateMessage(w http.ResponseWriter, r *http.Request) {
	var payload CreateMessagePayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	msg, err := json.Marshal(payload.Message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ac.qm.GetQueue(payload.QueueName).Enqueue(string(msg), time.Now())
	w.WriteHeader(201)

	ac.TcpServer.Broadcast("queue_name:" + payload.QueueName + " updated:1")
}

func (ac *ApiController) GetMessage(w http.ResponseWriter, r *http.Request) {
	queue := ac.qm.GetQueue(r.URL.Query().Get("queue_name"))

	gson, e := json.Marshal(queue.Dequeue())
	if e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(200)

	w.Write(gson)
}

func (ac *ApiController) Queues(w http.ResponseWriter, r *http.Request) {
	queue := ac.qm.GetQueue(r.URL.Query().Get("queue_name"))

	gson, e := json.Marshal(queue.Dequeue())
	if e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(200)
	w.Write(gson)
}

func (ac *ApiController) TCPClients(w http.ResponseWriter, r *http.Request) {
	gson, e := json.Marshal(map[string]int{"count": len(ac.TcpServer.GetClients().Connections)})
	if e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(200)
	w.Write(gson)
}

func (ac *ApiController) MessagesCount(w http.ResponseWriter, r *http.Request) {
	count := 0
	queue := ac.qm.ListQueues()

	for _, v := range queue {
		count += ac.qm.GetQueue(v).Size()
	}
	gson, e := json.Marshal(map[string]int{"count": count})
	if e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(200)
	w.Write(gson)
}

func (ac *ApiController) DeleteQueue(w http.ResponseWriter, r *http.Request) {
	queue := r.URL.Query().Get("queue_name")
	ac.qm.DeleteQueue(queue)
	w.WriteHeader(200)
}

func NewApiController(qm interfaces.QueueManagerInterface, logger providers.Logger) *ApiController {
	return &ApiController{
		qm:     qm,
		logger: logger,
	}
}
