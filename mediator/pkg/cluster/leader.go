package cluster

import (
	"bufio"
	"cognitivexr.at/cogstream/mediator/pkg/log"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
)

type WorkerConnection struct {
	conn    net.Conn
	info    *NodeInfo
	reader  *bufio.Reader
	remote  net.Addr
	isAlive bool
}

func (w *WorkerConnection) String() string {
	addrString := "<nil>"
	address, err := w.WebsocketAddress()
	if err == nil {
		addrString = "ws://" + address
	}
	return fmt.Sprintf("WorkerConnection(%s, isAlive=%t)", addrString, w.isAlive)
}

func (w *WorkerConnection) NodeInfo() *NodeInfo {
	return w.info
}

func (w *WorkerConnection) WebsocketAddress() (string, error) {
	if !w.isAlive {
		return "", errors.New("worker is not alive")
	}
	if w.remote == nil {
		return "", errors.New("no address set")
	}
	s := w.remote.String()
	if s == "<nil>" {
		return s, errors.New("no address set")
	}

	parts := strings.Split(s, ":")

	return fmt.Sprintf("%s:%d", parts[0], w.info.RpcPort), nil
}

func (w *WorkerConnection) IsAlive() bool {
	if w.conn == nil {
		return false
	}
	if w.remote == nil {
		return false
	}
	return w.isAlive
}

type Leader struct {
	address string
	workers map[NodeId]*WorkerConnection
	mutex   sync.Mutex
	socket  net.Listener
	running bool
}

func (l *Leader) ListWorkers() []*WorkerConnection {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	nodes := make([]*WorkerConnection, 0)

	for _, connection := range l.workers {
		if connection.IsAlive() {
			nodes = append(nodes, connection)
		}
	}

	return nodes
}

func NewLeader(address string) *Leader {
	return &Leader{
		address: address,
		workers: make(map[NodeId]*WorkerConnection),
		running: false,
	}
}

func (l *Leader) IsRunning() bool {
	return l.running
}

func (l *Leader) Shutdown() (err error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if l.socket != nil {
		err = l.socket.Close()
	}

	for _, worker := range l.workers {
		if worker.conn == nil {
			continue
		}
		worker.conn.Close()
	}

	return
}

func (l *Leader) Run() error {
	socket, err := net.Listen("tcp", l.address)
	if err != nil {
		return err
	}
	l.socket = socket
	defer socket.Close()

	for {
		l.running = true
		conn, err := socket.Accept()
		if err != nil {
			log.Error("error accepting connection: %s", err)
			l.running = false
			return err
		}
		go l.workerConnection(conn)
	}
}

func (l *Leader) workerConnection(conn net.Conn) {
	reader := bufio.NewReader(conn)

	decoder := json.NewDecoder(reader)
	var info NodeInfo
	err := decoder.Decode(&info)
	if err != nil {
		log.Error("error initializing worker connection %s: %s", conn.LocalAddr(), err)
		conn.Close()
		return
	}

	l.mutex.Lock()

	if oldConn, exists := l.workers[info.NodeId]; exists {
		log.Info("connection for node % exists, closing old one", oldConn)
		oldConn.isAlive = false
		oldConn.conn.Close()
	}

	wc := &WorkerConnection{
		conn:    conn,
		info:    &info,
		reader:  reader,
		remote:  conn.RemoteAddr(),
		isAlive: false,
	}
	l.workers[info.NodeId] = wc
	l.mutex.Unlock()

	// TODO: start timeout watchdog

	// start heartbeat read loop
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text != Heartbeat {
			err = errors.New("protocol error: did not receive heartbeat")
			break
		}
		wc.isAlive = true
		log.Debug("received heartbeat from %s", wc.remote)
	}
	wc.isAlive = false

	if err == nil {
		err = scanner.Err()
	}

	if err != nil {
		log.Info("connection for node %s terminated: %s", wc.remote, err)
	}

	// remove connection
	l.mutex.Lock()
	delete(l.workers, info.NodeId)
	l.mutex.Unlock()
}
