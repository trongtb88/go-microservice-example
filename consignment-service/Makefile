build:
	protoc --proto_path=$GOPATH/src:. --micro_out=. --go_out=. proto/consignment/consignment.proto

		GOOS=linux GOARCH=amd64 go build -o consignment-service && \
		  docker build -t consignment-service .

run:
	docker run -p 50052:50051 \
		-e MICRO_SERVER_ADDRESS=:50051 \
		consignment-service
