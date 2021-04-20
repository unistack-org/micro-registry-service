package service

import (
	pb "github.com/unistack-org/micro-register-service/v3/proto"
	"github.com/unistack-org/micro/v3/register"
)

func ToProto(s *register.Service) *pb.Service {
	endpoints := make([]*pb.Endpoint, 0, len(s.Endpoints))
	for _, ep := range s.Endpoints {
		endpoints = append(endpoints, &pb.Endpoint{
			Name:     ep.Name,
			Request:  ep.Request,
			Response: ep.Response,
			Metadata: ep.Metadata,
		})
	}

	nodes := make([]*pb.Node, 0, len(s.Nodes))

	for _, node := range s.Nodes {
		nodes = append(nodes, &pb.Node{
			Id:       node.Id,
			Address:  node.Address,
			Metadata: node.Metadata,
		})
	}

	return &pb.Service{
		Name:      s.Name,
		Version:   s.Version,
		Metadata:  s.Metadata,
		Endpoints: endpoints,
		Nodes:     nodes,
		Options:   new(pb.Options),
	}
}

func ToService(s *pb.Service) *register.Service {
	endpoints := make([]*register.Endpoint, 0, len(s.Endpoints))
	for _, ep := range s.Endpoints {
		endpoints = append(endpoints, &register.Endpoint{
			Name:     ep.Name,
			Request:  ep.Request,
			Response: ep.Response,
			Metadata: ep.Metadata,
		})
	}

	nodes := make([]*register.Node, 0, len(s.Nodes))
	for _, node := range s.Nodes {
		nodes = append(nodes, &register.Node{
			Id:       node.Id,
			Address:  node.Address,
			Metadata: node.Metadata,
		})
	}

	return &register.Service{
		Name:      s.Name,
		Version:   s.Version,
		Metadata:  s.Metadata,
		Endpoints: endpoints,
		Nodes:     nodes,
	}
}
