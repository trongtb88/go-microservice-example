package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	pb "github.com/trongtb88/go-microservice-example/consignment-service/proto/consignment"

	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/metadata"
)

const (
	address         = "127.0.0.1:50051"
	defaultFilename = "consignment.json"
)

func parseFile(file string) (*pb.Consignment, error) {
	var consignment *pb.Consignment
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(data, &consignment)
	log.Println(consignment.GetDescription())
	return consignment, err
}

func main() {

	serviceConsignment := micro.NewService(micro.Name("consignment.cli"))
	serviceConsignment.Init()
	// Create new greeter client
	shipping := pb.NewShippingService("consignment", serviceConsignment.Client())

	// Contact the server and print out its response.
	file := defaultFilename
	var token string
	if len(os.Args) > 1 {
		file = os.Args[1]
		token = os.Args[2]
	}

	consignment, err := parseFile(file)

	if err != nil {
		log.Fatalf("Could not parse file: %v", err)
	}

	ctx := metadata.NewContext(context.Background(), map[string]string{
		"token": token,
	})

	r, err := shipping.CreateConsignment(ctx, consignment)
	if err != nil {
		log.Fatalf("Could not create: %v", err)
	}
	log.Printf("Created: %t", r.Created)

	getAll, err := shipping.GetConsignments(ctx, &pb.GetRequest{})
	if err != nil {
		log.Fatalf("Could not list consignments: %v", err)
	}
	for _, v := range getAll.Consignments {
		log.Println(v)
	}
}
