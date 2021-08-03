module cognitivexr.at/cogstream/engines

go 1.16

require (
	cognitivexr.at/cogstream v0.0.0
	cognitivexr.at/cogstream/api v0.0.0
	gocv.io/x/gocv v0.27.0
)

replace (
	cognitivexr.at/cogstream v0.0.0 => ../../cogstream-go
	cognitivexr.at/cogstream/api v0.0.0 => ../../api
)
