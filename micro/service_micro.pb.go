// Code generated by protoc-gen-micro
// source: service.proto
package service

import (
	context "context"
	proto "github.com/unistack-org/micro-register-service/v3/proto"
	api "github.com/unistack-org/micro/v3/api"
	client "github.com/unistack-org/micro/v3/client"
)

func NewRegisterEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

type RegisterClient interface {
	LookupService(ctx context.Context, req *proto.LookupRequest, opts ...client.CallOption) (*proto.LookupResponse, error)
	Register(ctx context.Context, req *proto.Service, opts ...client.CallOption) (*proto.EmptyResponse, error)
	Deregister(ctx context.Context, req *proto.Service, opts ...client.CallOption) (*proto.EmptyResponse, error)
	ListServices(ctx context.Context, req *proto.ListRequest, opts ...client.CallOption) (*proto.ListResponse, error)
	Watch(ctx context.Context, req *proto.WatchRequest, opts ...client.CallOption) (Register_WatchClient, error)
}

type Register_WatchClient interface {
	Context() context.Context
	SendMsg(msg interface{}) error
	RecvMsg(msg interface{}) error
	Close() error
	Recv() (*proto.Result, error)
}

type RegisterServer interface {
	LookupService(ctx context.Context, req *proto.LookupRequest, rsp *proto.LookupResponse) error
	Register(ctx context.Context, req *proto.Service, rsp *proto.EmptyResponse) error
	Deregister(ctx context.Context, req *proto.Service, rsp *proto.EmptyResponse) error
	ListServices(ctx context.Context, req *proto.ListRequest, rsp *proto.ListResponse) error
	Watch(ctx context.Context, req *proto.WatchRequest, stream Register_WatchStream) error
}

type Register_WatchStream interface {
	Context() context.Context
	SendMsg(msg interface{}) error
	RecvMsg(msg interface{}) error
	Close() error
	Send(msg *proto.Result) error
}
