gen:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
  		hello/hello.proto
s:
	go run server/main.go

c: 
	go run client/main.go