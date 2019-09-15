package main

import (
	"context"
	"fmt"
	"os"

	pb "microservices/go-microservice-example/consignment-service/proto/consignment"

	vesselProto "github.com/trongtb88/go-microservice-example/vessel-service/proto/vessel"

	micro "github.com/micro/go-micro"

	log "github.com/sirupsen/logrus"
)

const (
	port = ":50051"
)

var logger = log.New()

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

}

type repository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
	GetAll() []*pb.Consignment
}

// Repository - Dummy repository, this simulates the use of a datastore
// of some kind. We'll replace this with a real implementation later on.
type Repository struct {
	consignments []*pb.Consignment
}

func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	updated := append(repo.consignments, consignment)
	repo.consignments = updated
	return consignment, nil
}

func (repo *Repository) GetAll() []*pb.Consignment {
	return repo.consignments
}

type service struct {
	repo          repository
	logger        *log.Logger
	vesselService vesselProto.VesselService
}

func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {
	var vesselResponse vesselProto.Response
	err := s.vesselService.FindAvailable(context.Background(), &vesselProto.Specification{
		MaxWeight: req.Weight,
		Capacity:  int32(len(req.Containers)),
	}, &vesselResponse)
	logger.WithField("vesselResponse.Vessel.Name", vesselResponse.Vessel.Name).Info("Found vessel:")
	if err != nil {
		logger.Error(err)
		return err
	}
	req.VesselId = vesselResponse.Vessel.GetId()
	// Save our consignment
	consignment, err := s.repo.Create(req)
	if err != nil {
		return err
	}

	// Return matching the `Response` message we created in our
	// protobuf definition.
	res.Created = true
	res.Consignment = consignment
	return nil
}

func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	consignments := s.repo.GetAll()
	res.Consignments = consignments
	return nil
}

func main() {

	repo := &Repository{}

	srv := micro.NewService(
		micro.Name("consignment"),
	)

	srv.Init()

	vesselService := vesselProto.NewVesselService("vessel", srv.Client())

	// Register handler
	pb.RegisterShippingServiceHandler(srv.Server(), &service{repo, logger, vesselService})

	// Run the server
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
