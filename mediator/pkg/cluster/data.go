package cluster

import (
	"cognitivexr.at/cogstream/api/engines"
)

type NodeId string

type NodeInfo struct {
	NodeId  NodeId                      `json:"nodeId"`
	RpcPort int                         `json:"rpcPort"`
	Engines []*engines.EngineDescriptor `json:"engines"`
}

const Heartbeat string = "ping"
