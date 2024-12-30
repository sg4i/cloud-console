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

func runHTTPServer(grpcAddress string, httpAddress string, token string) error {
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
		grpc.WithUnaryInterceptor(clientAuthInterceptor(token)),
	}
	err := pb.RegisterConsoleServiceHandlerFromEndpoint(ctx, mux, grpcAddress, opts)
	if err != nil {
		return err
	}

	httpHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/healthz" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
			return
		}
		if !isValidHTTPToken(r, token) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		mux.ServeHTTP(w, r)
	})

	server := &http.Server{
		Addr:    httpAddress,
		Handler: httpHandler,
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
