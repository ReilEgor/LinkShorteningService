package handler

import (
	"context"
	"errors"
	"log/slog"

	pb "github.com/ReilEgor/LinkShorteningService/internal/delivery/gRPC/gen"
	"github.com/ReilEgor/LinkShorteningService/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LinkGRPCHandler struct {
	pb.UnimplementedLinkServiceServer
	uc     domain.LinkUsecase
	logger *slog.Logger
}

func NewLinkGRPCHandler(uc domain.LinkUsecase, logger *slog.Logger) *LinkGRPCHandler {
	return &LinkGRPCHandler{
		uc:     uc,
		logger: logger,
	}
}

func (h *LinkGRPCHandler) AddLink(ctx context.Context, req *pb.AddLinkRequest) (*pb.LinkResponse, error) {
	link, err := h.uc.AddLink(ctx, req.GetLongURL())
	if err != nil {
		h.logger.Error("failed to add link", slog.Any("error", err))
		return nil, status.Errorf(codes.Internal, "failed to add link: %v", err)
	}

	return &pb.LinkResponse{
		Id:       link.ID,
		LongURL:  link.LongURL,
		ShortURL: link.ShortURL,
	}, nil
}

func (h *LinkGRPCHandler) GetLink(ctx context.Context, req *pb.GetLinkRequest) (*pb.LinkResponse, error) {
	link, err := h.uc.GetLink(ctx, req.GetShortURL())
	if err != nil {
		if errors.Is(err, domain.ErrLinkNotFound) {
			return nil, status.Error(codes.NotFound, "link not found")
		}
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return &pb.LinkResponse{
		Id:       link.ID,
		LongURL:  link.LongURL,
		ShortURL: link.ShortURL,
	}, nil
}
