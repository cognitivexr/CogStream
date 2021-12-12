package cluster

import (
	"cognitivexr.at/cogstream/api/messages"
	"cognitivexr.at/cogstream/mediator/pkg/log"
	"testing"
	"time"
)

func TestCluster(t *testing.T) {
	leader := NewLeader("0.0.0.0:4345")

	time.AfterFunc(1*time.Second, func() {
		info := &NodeInfo{
			NodeId:        "hello",
			WebsocketPort: 9501,
			AvailableEngines: messages.AvailableEngines{
				Engines: make([]*messages.EngineSpec, 10),
			},
		}

		log.Info("creating new cluster connection")
		c := NewClusterConnection("127.0.0.1:4345", info)

		time.AfterFunc(6*time.Second, func() {
			log.Info("closing client")
			c.Close()
		})

		err := c.Run()
		if err != nil {
			log.Error("error with cluster connection: %s", err)
		}
	})

	time.AfterFunc(10*time.Second, func() {
		log.Info("shutting down leader")
		leader.Shutdown()
	})

	go func() {
		time.Sleep(2 * time.Second)

		for {
			if !leader.IsRunning() {
				return
			}
			for _, connection := range leader.ListWorkers() {
				log.Info("available: connection: %s", connection)
			}
			time.Sleep(1 * time.Second)
		}

	}()

	log.Info("starting leader")
	leader.Run()
	log.Info("exiting")
}
