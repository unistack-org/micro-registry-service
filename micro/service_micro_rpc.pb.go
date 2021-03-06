// Code generated by protoc-gen-micro
// source: service.proto
package service

import (
	context "context"
	proto "github.com/unistack-org/micro-register-service/v3/proto"
	api "github.com/unistack-org/micro/v3/api"
	client "github.com/unistack-org/micro/v3/client"
	server "github.com/unistack-org/micro/v3/server"
)

type registerClient struct {
	c    client.Client
	name string
}

func NewRegisterClient(name string, c client.Client) RegisterClient {
	return &registerClient{c: c, name: name}
}

func (c *registerClient) LookupService(ctx context.Context, req *proto.LookupRequest, opts ...client.CallOption) (*proto.LookupResponse, error) {
	rsp := &proto.LookupResponse{}
	err := c.c.Call(ctx, c.c.NewRequest(c.name, "Register.LookupService", req), rsp, opts...)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

func (c *registerClient) Register(ctx context.Context, req *proto.Service, opts ...client.CallOption) (*proto.EmptyResponse, error) {
	rsp := &proto.EmptyResponse{}
	err := c.c.Call(ctx, c.c.NewRequest(c.name, "Register.Register", req), rsp, opts...)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

func (c *registerClient) Deregister(ctx context.Context, req *proto.Service, opts ...client.CallOption) (*proto.EmptyResponse, error) {
	rsp := &proto.EmptyResponse{}
	err := c.c.Call(ctx, c.c.NewRequest(c.name, "Register.Deregister", req), rsp, opts...)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

func (c *registerClient) ListServices(ctx context.Context, req *proto.ListRequest, opts ...client.CallOption) (*proto.ListResponse, error) {
	rsp := &proto.ListResponse{}
	err := c.c.Call(ctx, c.c.NewRequest(c.name, "Register.ListServices", req), rsp, opts...)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

func (c *registerClient) Watch(ctx context.Context, req *proto.WatchRequest, opts ...client.CallOption) (Register_WatchClient, error) {
	stream, err := c.c.Stream(ctx, c.c.NewRequest(c.name, "Register.Watch", &proto.WatchRequest{}), opts...)
	if err != nil {
		return nil, err
	}
	if err := stream.Send(req); err != nil {
		return nil, err
	}
	return &registerClientWatch{stream}, nil
}

type registerClientWatch struct {
	stream client.Stream
}

func (s *registerClientWatch) Close() error {
	return s.stream.Close()
}

func (s *registerClientWatch) Context() context.Context {
	return s.stream.Context()
}

func (s *registerClientWatch) SendMsg(msg interface{}) error {
	return s.stream.Send(msg)
}

func (s *registerClientWatch) RecvMsg(msg interface{}) error {
	return s.stream.Recv(msg)
}

func (s *registerClientWatch) Recv() (*proto.Result, error) {
	msg := &proto.Result{}
	if err := s.stream.Recv(msg); err != nil {
		return nil, err
	}
	return msg, nil
}

type registerServer struct {
	RegisterServer
}

func (h *registerServer) LookupService(ctx context.Context, req *proto.LookupRequest, rsp *proto.LookupResponse) error {
	return h.RegisterServer.LookupService(ctx, req, rsp)
}

func (h *registerServer) Register(ctx context.Context, req *proto.Service, rsp *proto.EmptyResponse) error {
	return h.RegisterServer.Register(ctx, req, rsp)
}

func (h *registerServer) Deregister(ctx context.Context, req *proto.Service, rsp *proto.EmptyResponse) error {
	return h.RegisterServer.Deregister(ctx, req, rsp)
}

func (h *registerServer) ListServices(ctx context.Context, req *proto.ListRequest, rsp *proto.ListResponse) error {
	return h.RegisterServer.ListServices(ctx, req, rsp)
}

func (h *registerServer) Watch(ctx context.Context, stream server.Stream) error {
	msg := &proto.WatchRequest{}
	if err := stream.Recv(msg); err != nil {
		return err
	}
	return h.RegisterServer.Watch(ctx, msg, &registerWatchStream{stream})
}

type registerWatchStream struct {
	stream server.Stream
}

func (s *registerWatchStream) Close() error {
	return s.stream.Close()
}

func (s *registerWatchStream) Context() context.Context {
	return s.stream.Context()
}

func (s *registerWatchStream) SendMsg(msg interface{}) error {
	return s.stream.Send(msg)
}

func (s *registerWatchStream) RecvMsg(msg interface{}) error {
	return s.stream.Recv(msg)
}

func (s *registerWatchStream) Send(msg *proto.Result) error {
	return s.stream.Send(msg)
}

func RegisterRegisterServer(s server.Server, sh RegisterServer, opts ...server.HandlerOption) error {
	type register interface {
		LookupService(ctx context.Context, req *proto.LookupRequest, rsp *proto.LookupResponse) error
		Register(ctx context.Context, req *proto.Service, rsp *proto.EmptyResponse) error
		Deregister(ctx context.Context, req *proto.Service, rsp *proto.EmptyResponse) error
		ListServices(ctx context.Context, req *proto.ListRequest, rsp *proto.ListResponse) error
		Watch(ctx context.Context, stream server.Stream) error
	}
	type Register struct {
		register
	}
	h := &registerServer{sh}
	var nopts []server.HandlerOption
	for _, endpoint := range NewRegisterEndpoints() {
		nopts = append(nopts, api.WithEndpoint(endpoint))
	}
	return s.Handle(s.NewHandler(&Register{h}, append(nopts, opts...)...))
}
