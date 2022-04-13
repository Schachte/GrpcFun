gen:
	protoc */**/*.proto --go_out=plugins=grpc:.

clean:
	rm -r */**/*.pb.go

_client_greeting:
	go run greet/greet_client/client.go

_server_greeting:
	go run greet/greet_server/server.go

_client_calculator:
	go run calculator/client/client.go

_server_calculator:
	go run calculator/server/server.go