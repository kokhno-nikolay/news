package server

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	desc "github.com/kokhno-nikolay/news/api/proto"
	"github.com/kokhno-nikolay/news/config"
	"github.com/kokhno-nikolay/news/internal/service"
)

type Server struct {
	desc.UnimplementedPostsServer
	postService service.PostService
}

func NewServer(postService service.PostService) *Server {
	return &Server{
		postService: postService,
	}
}

func (s *Server) StartGrpcServer(cfg *config.Config) error {
	grpcServer := grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
	)

	reflection.Register(grpcServer)

	desc.RegisterPostsServer(grpcServer, &Server{})

	list, err := net.Listen("tcp", cfg.GrpcAddress)
	if err != nil {
		return err
	}

	log.Printf("gRPC server listening at %v\n", cfg.GrpcAddress)

	return grpcServer.Serve(list)
}

func (s *Server) StartHttpServer(ctx context.Context, cfg *config.Config) error {
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err := desc.RegisterPostsHandlerFromEndpoint(ctx, mux, cfg.GrpcAddress, opts)
	if err != nil {
		return err
	}

	log.Printf("http server listening at %v\n", cfg.HttpAddress)

	return http.ListenAndServe(cfg.HttpAddress, mux)
}
