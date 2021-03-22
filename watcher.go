package service

import (
	pbmicro "github.com/unistack-org/micro-register-service/v3/micro"
	"github.com/unistack-org/micro/v3/register"
)

type serviceWatcher struct {
	stream pbmicro.Register_WatchClient
	closed chan bool
}

func (s *serviceWatcher) Next() (*register.Result, error) {
	// check if closed
	select {
	case <-s.closed:
		return nil, register.ErrWatcherStopped
	default:
	}

	r, err := s.stream.Recv()
	if err != nil {
		return nil, err
	}

	return &register.Result{
		Action:  r.Action,
		Service: ToService(r.Service),
	}, nil
}

func (s *serviceWatcher) Stop() {
	select {
	case <-s.closed:
		return
	default:
		close(s.closed)
		s.stream.Close()
	}
}

func newWatcher(stream pbmicro.Register_WatchClient) register.Watcher {
	return &serviceWatcher{
		stream: stream,
		closed: make(chan bool),
	}
}
