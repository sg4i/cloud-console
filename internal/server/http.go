package server

import (
	"context"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/sg4i/cloud-console/internal/logger"
	pb "github.com/sg4i/cloud-console/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func RunHTTPServer(grpcAddress string, httpAddress string, token string) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux(
		runtime.WithMetadata(func(ctx context.Context, req *http.Request) metadata.MD {
			md := make(map[string]string)
			if auth := req.Header.Get("Authorization"); auth != "" {
				md["authorization"] = auth
			}
			return metadata.New(md)
		}),
	)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	if token != "" {
		opts = append(opts, grpc.WithUnaryInterceptor(clientAuthInterceptor(token)))
	}

	err := pb.RegisterConsoleServiceHandlerFromEndpoint(ctx, mux, grpcAddress, opts)
	if err != nil {
		logger.Log.Errorf("注册 gRPC 服务处理程序失败: %v", err)
		return err
	}

	httpHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Log.Infof("处理请求: %s %s", r.Method, r.URL.Path)

		if r.URL.Path == "/healthz" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
			return
		}

		if token != "" {
			if !isValidHTTPToken(r, token) {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		}

		mux.ServeHTTP(w, r)
	})

	corsHandler := CORSHandler(httpHandler)
	server := &http.Server{
		Addr:    httpAddress,
		Handler: corsHandler,
	}

	logger.Log.Infof("Starting HTTP server on %s", httpAddress)
	return server.ListenAndServe()
}

func isValidHTTPToken(r *http.Request, token string) bool {
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		return false
	}
	authToken := strings.TrimPrefix(auth, "Bearer ")
	return isValidToken(authToken, token)
}
