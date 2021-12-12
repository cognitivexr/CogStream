package cluster

import (
	"cognitivexr.at/cogstream/api/messages"
)

type NodeId string

type NodeInfo struct {
	NodeId           NodeId                    `json:"nodeId"`
	WebsocketPort    int                       `json:"websocketPort"`
	AvailableEngines messages.AvailableEngines `json:"availableEngines"`
}

const Heartbeat string = "ping"
