package service

import (
	pb "github.com/unistack-org/micro-register-service/v3/proto"
	"github.com/unistack-org/micro/v3/register"
)

func values(v []*register.Value) []*pb.Value {
	if len(v) == 0 {
		return []*pb.Value{}
	}

	vs := make([]*pb.Value, 0, len(v))
	for _, vi := range v {
		vs = append(vs, &pb.Value{
			Name:   vi.Name,
			Type:   vi.Type,
			Values: values(vi.Values),
		})
	}
	return vs
}

func toValues(v []*pb.Value) []*register.Value {
	if len(v) == 0 {
		return []*register.Value{}
	}

	vs := make([]*register.Value, 0, len(v))
	for _, vi := range v {
		vs = append(vs, &register.Value{
			Name:   vi.Name,
			Type:   vi.Type,
			Values: toValues(vi.Values),
		})
	}
	return vs
}

func ToProto(s *register.Service) *pb.Service {
	endpoints := make([]*pb.Endpoint, 0, len(s.Endpoints))
	for _, ep := range s.Endpoints {
		var request, response *pb.Value

		if ep.Request != nil {
			request = &pb.Value{
				Name:   ep.Request.Name,
				Type:   ep.Request.Type,
				Values: values(ep.Request.Values),
			}
		}

		if ep.Response != nil {
			response = &pb.Value{
				Name:   ep.Response.Name,
				Type:   ep.Response.Type,
				Values: values(ep.Response.Values),
			}
		}

		endpoints = append(endpoints, &pb.Endpoint{
			Name:     ep.Name,
			Request:  request,
			Response: response,
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
		var request, response *register.Value

		if ep.Request != nil {
			request = &register.Value{
				Name:   ep.Request.Name,
				Type:   ep.Request.Type,
				Values: toValues(ep.Request.Values),
			}
		}

		if ep.Response != nil {
			response = &register.Value{
				Name:   ep.Response.Name,
				Type:   ep.Response.Type,
				Values: toValues(ep.Response.Values),
			}
		}

		endpoints = append(endpoints, &register.Endpoint{
			Name:     ep.Name,
			Request:  request,
			Response: response,
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
