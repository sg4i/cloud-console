package server

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func AuthInterceptor(token string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "无元数据")
		}

		authToken, ok := md["authorization"]
		if !ok || len(token) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "缺少授权令牌")
		}

		tokenStr := authToken[0]
		if len(tokenStr) <= 7 {
			return nil, status.Errorf(codes.Unauthenticated, "无效的令牌格式")
		}

		// 移除 "Bearer " 前缀（如果存在）
		if tokenStr[:7] == "Bearer " {
			tokenStr = tokenStr[7:]
		} else {
			return nil, status.Errorf(codes.Unauthenticated, "无效的令牌格式")
		}

		if !isValidToken(tokenStr, token) {
			return nil, status.Errorf(codes.Unauthenticated, "无效的令牌")
		}

		return handler(ctx, req)
	}
}

func isValidToken(value, token string) bool {
	return value == token
}

func clientAuthInterceptor(token string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
