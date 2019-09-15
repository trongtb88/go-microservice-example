package main

import (
	"context"
	"errors"
	"fmt"

	pb_vessel "microservices/go-microservice-example/vessel-service/proto/vessel"

	micro "github.com/micro/go-micro"
)

type Repository interface {
	FindAvailable(*pb_vessel.Specification) (*pb_vessel.Vessel, error)
}

type VesselRepository struct {
	vessels []pb_vessel.Vessel
}

func (repo *VesselRepository) FindAvailable(spec *pb_vessel.Specification) (*pb_vessel.Vessel, error) {
	for _, vessel := range repo.vessels {
		if spec.GetCapacity() <= vessel.GetCapacity() && spec.MaxWeight <= vessel.MaxWeight {
			return &vessel, nil
		}
	}
	return nil, errors.New("No vessel found by that spec")
}

type service struct {
	repo Repository
}

func (s *service) FindAvailable(ctx context.Context, req *pb_vessel.Specification, res *pb_vessel.Response) error {
	vessel, err := s.repo.FindAvailable(req)
	if err != nil {
		return err
	}
	res.Vessel = vessel
	return nil
}

func main() {
	vessels := []pb_vessel.Vessel{
		pb_vessel.Vessel{Id: "vessel001", Name: "Boaty McBoatface", MaxWeight: 200000, Capacity: 500},
	}
	repo := &VesselRepository{vessels}

	srv := micro.NewService(
		micro.Name("vessel"),
	)

	srv.Init()

	// Register our implementation with
	pb_vessel.RegisterVesselServiceHandler(srv.Server(), &service{repo})

	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
