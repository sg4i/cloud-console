package console

import (
	pb "github.com/sg4i/cloud-console/proto"

	"net"

	"github.com/sg4i/cloud-console/internal/logger"

	grpc_server "github.com/sg4i/cloud-console/internal/server"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedConsoleServiceServer
}

func NewServer() *server {
	return &server{}
}

func StartRPCServer(grpcAddress string, httpAddress string, authToken string) error {
	// 启动gRPC服务器
	lis, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		return err
	}
	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_server.AuthInterceptor(authToken)),
	)
	pb.RegisterConsoleServiceServer(s, NewServer())

	// 启动gRPC服务器
	go func() {
		logger.Log.Infof("Starting gRPC server on %s", grpcAddress)
		if err := s.Serve(lis); err != nil {
			logger.Log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	// 如果httpAddress不为空,则启动HTTP网关
	if httpAddress != "" {
		return grpc_server.RunHTTPServer(grpcAddress, httpAddress, authToken)
	}

	// 如httpAddress为空,则阻塞主线程,保持gRPC服务器运行
	select {}
}
