build:
	protoc --proto_path=$GOPATH/src:. --micro_out=. --go_out=. proto/vessel/vessel.proto
	docker build -t vessel-service .

run:
	docker run -p 50053:50052 -e MICRO_SERVER_ADDRESS=:50052 vessel-service
