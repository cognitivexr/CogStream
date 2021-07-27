module cognitivexr.at/cogstream

go 1.16

require (
	cognitivexr.at/cogstream/api v0.0.0
	github.com/go-redis/redis/v8 v8.11.0
	github.com/pion/rtcp v1.2.6
	github.com/pion/webrtc/v3 v3.0.25
	gocv.io/x/gocv v0.27.0
)

replace cognitivexr.at/cogstream/api v0.0.0 => ../api
