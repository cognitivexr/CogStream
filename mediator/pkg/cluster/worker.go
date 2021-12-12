package cluster

import (
	"bufio"
	"cognitivexr.at/cogstream/mediator/pkg/log"
	"encoding/json"
	"net"
	"time"
)

type ClusterConnection struct {
	address string
	info    *NodeInfo
	conn    net.Conn
}

func NewClusterConnection(address string, info *NodeInfo) *ClusterConnection {
	return &ClusterConnection{
		address: address,
		info:    info,
	}
}

func (c *ClusterConnection) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *ClusterConnection) Run() error {
	info, err := json.Marshal(c.info)
	if err != nil {
		return err
	}

	log.Debug("connecting to cluster leader %s", c.address)
	conn, err := net.Dial("tcp", c.address)
	if err != nil {
		return err
	}
	c.conn = conn
	defer conn.Close()

	w := bufio.NewWriter(conn)
	log.Debug("writing info %s to leader", info)
	_, err = w.Write(info)
	if err != nil {
		return err
	}
	err = w.Flush()
	if err != nil {
		return err
	}

	for {
		log.Debug("sending heartbeat to %s", conn.RemoteAddr())
		_, err := w.WriteString(Heartbeat + "\n")
		if err != nil {
			return err
		}
		err = w.Flush()
		if err != nil {
			return err
		}
		time.Sleep(2 * time.Second)
	}
}
