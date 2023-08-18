package structures

// swagger:model DetailedQueuesDashboardResponse
type DetailedQueuesDashboardResponse struct {
	QueueName     string   `json:"queue_name"`
	QueueSize     int      `json:"queue_size"`
	QueueMessages []string `json:"queue_messages"`
}
