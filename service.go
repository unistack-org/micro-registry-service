// Package service uses the register service
package service

import (
	"context"
	"fmt"
	"time"

	pb "github.com/unistack-org/micro-register-service/v3/proto"
	"github.com/unistack-org/micro/v3/client"
	"github.com/unistack-org/micro/v3/errors"
	"github.com/unistack-org/micro/v3/register"
)

type serviceRegister struct {
	opts register.Options
	// name of the register
	service string
	// address
	address []string
	// client to call register
	client pb.RegisterService
}

func (s *serviceRegister) callOpts() []client.CallOption {
	var opts []client.CallOption

	// set register address
	if len(s.address) > 0 {
		opts = append(opts, client.WithAddress(s.address...))
	}

	// set timeout
	if s.opts.Timeout > time.Duration(0) {
		opts = append(opts, client.WithRequestTimeout(s.opts.Timeout))
	}

	return opts
}

func (s *serviceRegister) Init(opts ...register.Option) error {
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

	s.client = pb.NewRegisterService(s.service, cli)

	return nil
}

func (s *serviceRegister) Options() register.Options {
	return s.opts
}

func (s *serviceRegister) Connect(ctx context.Context) error {
	return nil
}

func (s *serviceRegister) Disconnect(ctx context.Context) error {
	return nil
}

func (s *serviceRegister) Register(ctx context.Context, srv *register.Service, opts ...register.RegisterOption) error {
	options := register.NewRegisterOptions(opts...)

	// encode srv into protobuf and pack TTL and domain into it
	pbSrv := ToProto(srv)
	pbSrv.Options.Ttl = int64(options.TTL.Seconds())
	pbSrv.Options.Domain = options.Domain

	// register the service
	_, err := s.client.Register(ctx, pbSrv, s.callOpts()...)
	return err
}

func (s *serviceRegister) Deregister(ctx context.Context, srv *register.Service, opts ...register.DeregisterOption) error {
	options := register.NewDeregisterOptions(opts...)

	// encode srv into protobuf and pack domain into it
	pbSrv := ToProto(srv)
	pbSrv.Options.Domain = options.Domain

	// deregister the service
	_, err := s.client.Deregister(ctx, pbSrv, s.callOpts()...)
	return err
}

func (s *serviceRegister) LookupService(ctx context.Context, name string, opts ...register.LookupOption) ([]*register.Service, error) {
	options := register.NewLookupOptions(opts...)

	rsp, err := s.client.LookupService(ctx, &pb.LookupRequest{
		Service: name, Options: &pb.Options{Domain: options.Domain},
	}, s.callOpts()...)

	if verr, ok := err.(*errors.Error); ok && verr.Code == 404 {
		return nil, register.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	services := make([]*register.Service, 0, len(rsp.Services))
	for _, service := range rsp.Services {
		services = append(services, ToService(service))
	}
	return services, nil
}

func (s *serviceRegister) ListServices(ctx context.Context, opts ...register.ListOption) ([]*register.Service, error) {
	options := register.NewListOptions(opts...)

	req := &pb.ListRequest{Options: &pb.Options{Domain: options.Domain}}
	rsp, err := s.client.ListServices(ctx, req, s.callOpts()...)
	if err != nil {
		return nil, err
	}

	services := make([]*register.Service, 0, len(rsp.Services))
	for _, service := range rsp.Services {
		services = append(services, ToService(service))
	}

	return services, nil
}

func (s *serviceRegister) Watch(ctx context.Context, opts ...register.WatchOption) (register.Watcher, error) {
	options := register.NewWatchOptions(opts...)

	stream, err := s.client.Watch(ctx, &pb.WatchRequest{
		Service: options.Service, Options: &pb.Options{Domain: options.Domain},
	}, s.callOpts()...)

	if err != nil {
		return nil, err
	}

	return newWatcher(stream), nil
}

func (s *serviceRegister) String() string {
	return "service"
}

func (s *serviceRegister) Name() string {
	return s.opts.Name
}

// NewRegister returns a new register service client
func NewRegister(opts ...register.Option) register.Register {
	options := register.NewOptions(opts...)

	addrs := options.Addrs
	if len(addrs) == 0 {
		addrs = []string{"127.0.0.1:8000"}
	}

	return &serviceRegister{
		opts:    options,
		address: addrs,
	}
}
