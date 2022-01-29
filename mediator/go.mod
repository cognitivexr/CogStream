module cognitivexr.at/cogstream/mediator

go 1.16

require (
	cognitivexr.at/cogstream/api v0.0.0
	github.com/containerd/containerd v1.5.9 // indirect
	github.com/docker/docker v20.10.12+incompatible // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/gorilla/websocket v1.4.2
	google.golang.org/grpc v1.44.0 // indirect
)

replace cognitivexr.at/cogstream/api v0.0.0 => ../api
