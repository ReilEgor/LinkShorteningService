package gRPC

import (
	"context" // Добавляем контекст
	"log/slog"
	"net"

	pb "github.com/ReilEgor/LinkShorteningService/internal/delivery/gRPC/gen"
	"github.com/ReilEgor/LinkShorteningService/internal/delivery/gRPC/handler"
	"github.com/ReilEgor/LinkShorteningService/internal/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func RunGRPCServer(ctx context.Context, port string, uc domain.LinkUsecase, logger *slog.Logger) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	linkHandler := handler.NewLinkGRPCHandler(uc, logger)
	pb.RegisterLinkServiceServer(s, linkHandler)
	reflection.Register(s)

	serverErrors := make(chan error, 1)

	go func() {
		logger.Info("gRPC server is starting", slog.String("port", port))
		if err := s.Serve(lis); err != nil {
			serverErrors <- err
		}
	}()

	select {
	case err := <-serverErrors:
		return err
	case <-ctx.Done():
		logger.Info("gRPC server is shutting down gracefully")
		s.GracefulStop()
		return nil
	}
}
