module example

go 1.19

require mycache v0.0.0

require (
	github.com/golang/protobuf v1.5.2 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)

replace mycache => ./mycache
