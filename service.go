// Package service uses the registry service
package service

import (
	"context"
	"fmt"
	"time"

	pb "github.com/unistack-org/micro-registry-service/v3/proto"
	"github.com/unistack-org/micro/v3/client"
	"github.com/unistack-org/micro/v3/errors"
	"github.com/unistack-org/micro/v3/registry"
)

type serviceRegistry struct {
	opts registry.Options
	// name of the registry
	service string
	// address
	address []string
	// client to call registry
	client pb.RegistryService
}

func (s *serviceRegistry) callOpts() []client.CallOption {
	var opts []client.CallOption

	// set registry address
	if len(s.address) > 0 {
		opts = append(opts, client.WithAddress(s.address...))
	}

	// set timeout
	if s.opts.Timeout > time.Duration(0) {
		opts = append(opts, client.WithRequestTimeout(s.opts.Timeout))
	}

	return opts
}

func (s *serviceRegistry) Init(opts ...registry.Option) error {
	for _, o := range opts {
		o(&s.opts)
	}

	if len(s.opts.Addrs) > 0 {
		s.address = s.opts.Addrs
	}

	var cli client.Client
	if s.opts.Context != nil {
		if v, ok := s.opts.Context.Value(clientKey{}).(string); ok && v != "" {
			s.service = v
		}
		if v, ok := s.opts.Context.Value(clientKey{}).(client.Client); ok && v != nil {
			cli = v
		}
	}

	if cli == nil {
		return fmt.Errorf("missing Client option")
	}

	if s.service == "" {
		return fmt.Errorf("missing Service option")
	}

	s.client = pb.NewRegistryService(s.service, cli)

	return nil
}

func (s *serviceRegistry) Options() registry.Options {
	return s.opts
}

func (s *serviceRegistry) Connect(ctx context.Context) error {
	return nil
}

func (s *serviceRegistry) Disconnect(ctx context.Context) error {
	return nil
}

func (s *serviceRegistry) Register(ctx context.Context, srv *registry.Service, opts ...registry.RegisterOption) error {
	options := registry.NewRegisterOptions(opts...)

	// encode srv into protobuf and pack TTL and domain into it
	pbSrv := ToProto(srv)
	pbSrv.Options.Ttl = int64(options.TTL.Seconds())
	pbSrv.Options.Domain = options.Domain

	// register the service
	_, err := s.client.Register(ctx, pbSrv, s.callOpts()...)
	return err
}

func (s *serviceRegistry) Deregister(ctx context.Context, srv *registry.Service, opts ...registry.DeregisterOption) error {
	options := registry.NewDeregisterOptions(opts...)

	// encode srv into protobuf and pack domain into it
	pbSrv := ToProto(srv)
	pbSrv.Options.Domain = options.Domain

	// deregister the service
	_, err := s.client.Deregister(ctx, pbSrv, s.callOpts()...)
	return err
}

func (s *serviceRegistry) GetService(ctx context.Context, name string, opts ...registry.GetOption) ([]*registry.Service, error) {
	options := registry.NewGetOptions(opts...)

	rsp, err := s.client.GetService(ctx, &pb.GetRequest{
		Service: name, Options: &pb.Options{Domain: options.Domain},
	}, s.callOpts()...)

	if verr, ok := err.(*errors.Error); ok && verr.Code == 404 {
		return nil, registry.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	services := make([]*registry.Service, 0, len(rsp.Services))
	for _, service := range rsp.Services {
		services = append(services, ToService(service))
	}
	return services, nil
}

func (s *serviceRegistry) ListServices(ctx context.Context, opts ...registry.ListOption) ([]*registry.Service, error) {
	options := registry.NewListOptions(opts...)

	req := &pb.ListRequest{Options: &pb.Options{Domain: options.Domain}}
	rsp, err := s.client.ListServices(ctx, req, s.callOpts()...)
	if err != nil {
		return nil, err
	}

	services := make([]*registry.Service, 0, len(rsp.Services))
	for _, service := range rsp.Services {
		services = append(services, ToService(service))
	}

	return services, nil
}

func (s *serviceRegistry) Watch(ctx context.Context, opts ...registry.WatchOption) (registry.Watcher, error) {
	options := registry.NewWatchOptions(opts...)

	stream, err := s.client.Watch(ctx, &pb.WatchRequest{
		Service: options.Service, Options: &pb.Options{Domain: options.Domain},
	}, s.callOpts()...)

	if err != nil {
		return nil, err
	}

	return newWatcher(stream), nil
}

func (s *serviceRegistry) String() string {
	return "service"
}

// NewRegistry returns a new registry service client
func NewRegistry(opts ...registry.Option) registry.Registry {
	options := registry.NewOptions(opts...)

	addrs := options.Addrs
	if len(addrs) == 0 {
		addrs = []string{"127.0.0.1:8000"}
	}

	return &serviceRegistry{
		opts:    options,
		address: addrs,
	}
}
