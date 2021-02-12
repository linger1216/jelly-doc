package server

import (
	"github.com/linger1216/go-utils/config"
	"github.com/linger1216/go-utils/log"
	"net"
	"net/http"
	"net/http/pprof"
	// 3d Party
	"google.golang.org/grpc"

	// This Service
	"github.com/linger1216/jelly-doc/src/server/api-service/handlers"
	"github.com/linger1216/jelly-doc/src/server/api-service/svc"
	pb "github.com/linger1216/jelly-doc/src/server/pb"
)

// Config contains the required fields for running a server
type Config struct {
	HTTPAddr  string
	DebugAddr string
	GRPCAddr  string
}

func NewEndpoints(service pb.ApiServer) svc.Endpoints {
	// Business domain.

	// Wrap Service with middlewares. See handlers/middlewares.go
	service = handlers.WrapService(service)

	// Endpoint domain.
	var (
		createEndpoint = svc.MakeCreateEndpoint(service)
		getEndpoint    = svc.MakeGetEndpoint(service)
		listEndpoint   = svc.MakeListEndpoint(service)
		updateEndpoint = svc.MakeUpdateEndpoint(service)
		deleteEndpoint = svc.MakeDeleteEndpoint(service)
	)

	endpoints := svc.Endpoints{
		CreateEndpoint: createEndpoint,
		GetEndpoint:    getEndpoint,
		ListEndpoint:   listEndpoint,
		UpdateEndpoint: updateEndpoint,
		DeleteEndpoint: deleteEndpoint,
	}

	// Wrap selected Endpoints with middlewares. See handlers/middlewares.go
	endpoints = handlers.WrapEndpoints(endpoints)

	return endpoints
}

// Run starts a new http server, gRPC server, and a debug server with the
// passed config and logger
func Run(reader config.Reader) {

	// logger
	logger := log.NewLog()

	service := handlers.NewService(logger, reader)
	endpoints := NewEndpoints(service)

	// Mechanical domain.
	interrupt := make(chan error)

	// Interrupt handler.
	go handlers.InterruptHandler(interrupt)

	// Debug listener.
	go func() {
		addr := reader.GetString("server", "debugAddr")
		logger.Debugf("debug addr:%s", addr)

		m := http.NewServeMux()
		m.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
		m.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
		m.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
		m.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
		m.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))

		interrupt <- http.ListenAndServe(addr, m)
	}()

	// HTTP transport.
	go func() {
		addr := reader.GetString("server", "httpAddr")
		logger.Debugf("http addr:%s", addr)
		h := svc.MakeHTTPHandler(endpoints)
		interrupt <- http.ListenAndServe(addr, h)
	}()

	// gRPC transport.
	go func() {
		addr := reader.GetString("server", "grpcAddr")
		logger.Debugf("grpc addr:%s", addr)
		ln, err := net.Listen("tcp", addr)
		if err != nil {
			interrupt <- err
			return
		}

		srv := svc.MakeGRPCServer(endpoints)
		s := grpc.NewServer()
		pb.RegisterApiServer(s, srv)

		interrupt <- s.Serve(ln)
	}()

	// Run!
	logger.Debugf("exit", <-interrupt)
}
