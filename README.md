
<center><img src="https://raw.githubusercontent.com/DimitriyKuschel/qm-hermes/master/public/img/logo.svg" alt="Alt Text" width="300" height="200"></center>

### Queue Manager Hermes
![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)

Queue Manager is a versatile message queuing application designed to help you manage queues of messages efficiently. It provides both a RESTful API and a TCP server for seamless communication, making it easy to integrate with various applications and services.

#### Features

* **Queue Management**: Create, delete, and list queues with ease.
* **Message Handling**: Enqueue and dequeue messages from queues.
* **RESTful API**: Access queues and messages via well-defined API endpoints.
* **TCP Server**: Establish real-time communication with clients using a TCP connection.

#### Installation
* Install goreleaser from https://goreleaser.com/install/
* Clone the repository: git clone https://github.com/DimitriyKuschel/qm-hermes.git
* Navigate to the project directory: cd qm-hermes
* Build the project: `go get -u && go vet && git tag -f v1.0.0 ; goreleaser --rm-dist --skip-publish --skip-validate`
* Run the proper binary for your platform from: `dist` directory 
* The application will start with the default configuration. You can modify the configuration file according to your requirements.
* By Default web interface will be available at http://localhost:8090 and tcp server at 127.0.0.1:8081

#### Documentation

The Queue Manager Hermes API provides powerful features for managing queues of messages efficiently and seamlessly. With a combination of RESTful endpoints and a real-time TCP server, you can easily send, receive, and manage messages within your application.

##### Endpoints
The api schema will be available at http://localhost:8090/schema

Certainly! Here's the provided information formatted as a table in `README.md` format:

| Action                 | Description                                      | URL            | Method | Payload/Parameters                                         | Response                         |
|------------------------|--------------------------------------------------|----------------|--------|------------------------------------------------------------|----------------------------------|
| Create New Queue       | Create a new queue with a specified name.       | `/api/v1/queue/create`   | POST   | `{"name": "my-queue"}`                                     | Status 201 Created               |
| List Queues            | Retrieve a list of all available queues.        | `/api/v1/queue/list`  | GET    | -                                                          | `["queue1", "queue2", "queue3"]` |
| Enqueue Message        | Add a message to a specific queue.              | `/api/v1/message/create` | POST   | `{"queue_name": "my-queue", "message": "Hello, World!"}` | Status 201 Created               |
| Dequeue Message        | Retrieve and remove the next message from a queue.| `/api/v1/message/get` | GET    | `queue_name`: The name of the queue to dequeue from       | `"Hello, World!"`                |
| List Messages in Queue| Retrieve all messages from a specific queue.    | `/api/v1/message/list` | GET    | `queue_name`: The name of the queue to retrieve messages from| `["Message 1", "Message 2", "Message 3"]` | 


##### Real-time Communication

Queue Manager also provides a real-time TCP server that allows you to establish persistent connections for immediate updates when the queue changes. The TCP server broadcasts messages whenever a queue is updated.

#### Contributing

Contributions are welcome! If you find any issues or want to enhance the project, feel free to open a pull request or submit an issue.

#### ToDo
1. [ ] Add load tests and benchmarks
2. [ ] Add Unit tests
3. [ ] Add web-server SSL support
4. [ ] Add web-server basic auth support for end-points
5. [ ] User group and roles for the Dashboard
6. [ ] Master-Master replication support


#### License

This project is licensed under the Apache License 2.0. See the LICENSE file for details.

#### Contact

For questions or feedback, you can reach us at info@greenline-software.com.