package console

import (
	"context"

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

func (s *server) GenerateRoleLoginURL(ctx context.Context, req *pb.GenerateRoleLoginURLRequest) (*pb.GenerateRoleLoginURLResponse, error) {
	// 创建 Console 实例
	console, err := New(&Options{
		Mode:     "grpc",
		Provider: Provider(req.Provider),
	})
	if err != nil {
		return nil, err
	}

	// 处理可选字段
	var token, roleArn, destination, loginUrl string
	if req.Token != nil {
		token = *req.Token
	}
	if req.RoleArn != nil {
		roleArn = *req.RoleArn
	}
	if req.Desiontion != nil {
		destination = *req.Desiontion
	}
	if req.LoginUrl != nil {
		loginUrl = *req.LoginUrl
	}

	// 将请求参数转换为 LoginOptions
	loginOpts := NewLoginOptions(
		req.SecretId,
		req.SecretKey,
		token,
		roleArn,
		destination,
		loginUrl,
	)

	// 调用 GetLoginURL 获取登录 URL
	url, err := console.GetLoginURL(loginOpts)
	if err != nil {
		return nil, err
	}

	return &pb.GenerateRoleLoginURLResponse{Url: url}, nil
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
