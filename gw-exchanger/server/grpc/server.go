package grpc

import (
	"gw-exchanger/gw-exchanger/proto/exchange"
	"gw-exchanger/internal/pkg/logger"
	"gw-exchanger/service"
	"net"

	"google.golang.org/grpc"
)

func RunGRPCServer(svc *service.ExchangeService, port string) error {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	exchange.RegisterExchangeServiceServer(server, svc)

	logger.L.Info("gRPC server listening on %s", port)
	return server.Serve(listener)
}
