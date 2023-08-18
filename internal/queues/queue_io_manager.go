package queues

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"queue-manager/internal/queues/interfaces"
	"time"
)

type QueueIOManager struct {
	compressor interfaces.CompressorInterface
}

func NewQueueIOManager(compressor interfaces.CompressorInterface) interfaces.QueueIOManagerInterface {
	return &QueueIOManager{
		compressor: compressor,
	}
}
func (qio *QueueIOManager) SaveQueuesToFile(qm interfaces.QueueManagerInterface, fileName string) error {
	queues := qm.GetAllQueues()

	// Marshal the data to a byte slice
	jsonData, err := json.Marshal(queues)
	if err != nil {
		return err
	}
	data, err := qio.compressor.Compress(jsonData)

	if err != nil {
		return err
	}

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the compressed data to the binary file
	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (qio *QueueIOManager) LoadQueuesFromFile(qm interfaces.QueueManagerInterface, fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	bufferedReader := bufio.NewReader(file)

	var data []byte
	buffer := make([]byte, 1024)
	for {
		n, err := bufferedReader.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}

		data = append(data, buffer[:n]...)
	}

	decompressedData, err := qio.compressor.Decompress(data)
	if err != nil {
		return err
	}

	var queues map[string][]string
	err = json.Unmarshal(decompressedData, &queues)
	if err != nil {
		return err
	}

	for name, queue := range queues {
		qm.CreateQueue(name)
		for _, val := range queue {
			qm.GetQueue(name).Enqueue(val, time.Now())
		}
	}

	return nil
}
