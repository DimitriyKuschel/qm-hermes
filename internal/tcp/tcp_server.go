package tcp

import (
	"bufio"
	"github.com/spf13/cast"
	"net"
	"queue-manager/internal/providers"
	"queue-manager/internal/structures"
	"sync"
)

const (
	ServerProtocol = "tcp"
)

type TcpServer struct {
	config *structures.Config
	logger providers.Logger
}

type ClientsMap struct {
	sync.RWMutex
	Connections map[net.Conn]bool
}

var Clients = &ClientsMap{
	Connections: make(map[net.Conn]bool),
}

func (t *TcpServer) Run() {
	t.logger.Infof(providers.TypeApp, "Starting server TCP server ...")

	listen, err := net.Listen(ServerProtocol, t.config.TcpServer.Host+":"+cast.ToString(t.config.TcpServer.Port))
	if err != nil {
		t.logger.Errorf(providers.TypeApp, "Error on listen: ", err.Error())
		return
	}

	defer listen.Close()

	t.logger.Infof(providers.TypeApp, "Listening TCP clients on "+t.config.TcpServer.Host+":"+cast.ToString(t.config.TcpServer.Port))

	for {
		conn, err := listen.Accept()
		if err != nil {
			t.logger.Errorf(providers.TypeApp, "Error on accept: ", err.Error())
			return
		}
		Clients.Lock()
		Clients.Connections[conn] = true
		Clients.Unlock()

		go t.handleRequest(conn)
	}
}

func (t *TcpServer) handleRequest(conn net.Conn) {
	for {
		_, err := bufio.NewReader(conn).ReadString('\n')

		if err != nil {
			t.logger.Errorf(providers.TypeApp, "Error reading:", err.Error())
			Clients.Lock()
			delete(Clients.Connections, conn)
			Clients.Unlock()
			conn.Close()
			break
		}

	}
}

func (t *TcpServer) Broadcast(message string) {
	Clients.Lock()
	for clientConn := range Clients.Connections {
		clientConn.Write([]byte(message))
	}
	Clients.Unlock()
}

func (t *TcpServer) GetClients() *ClientsMap {
	return Clients
}

func NewTcpServer(config *structures.Config, logger providers.Logger) *TcpServer {
	return &TcpServer{
		config: config,
		logger: logger,
	}
}
