package handler

import (
	"log/slog"
	"net/http"

	"github.com/ReilEgor/LinkShorteningService/internal/domain"
	"github.com/gin-gonic/gin"
)

type LinkHandler struct {
	uc     domain.LinkUsecase
	logger *slog.Logger
}

func NewLinkHandler(uc domain.LinkUsecase, logger *slog.Logger) *LinkHandler {
	return &LinkHandler{
		uc:     uc,
		logger: logger,
	}
}

func (h *LinkHandler) AddLink(c *gin.Context) {
	var req struct {
		LongURL string `json:"longURL" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	ctx := c.Request.Context()

	link, err := h.uc.AddLink(ctx, req.LongURL)
	if err != nil {
		h.logger.Error("failed to add link",
			slog.String("url", req.LongURL),
			slog.Any("error", err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	h.logger.InfoContext(ctx, "link added successfully", slog.String("short", link.ShortURL))
	c.JSON(http.StatusCreated, link)
}

func (h *LinkHandler) GetLink(c *gin.Context) {
	shortCode := c.Param("shortURL")
	ctx := c.Request.Context()
	origLink, err := h.uc.GetLink(ctx, shortCode)
	if err != nil {
		h.logger.Error("failed to get link",
			slog.String("code", shortCode),
			slog.Any("error", err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if origLink.LongURL == "" {
		h.logger.Warn(domain.ErrLinkNotFound.Error(), slog.String("code", shortCode))
		c.JSON(http.StatusNotFound, gin.H{"error": domain.ErrLinkNotFound.Error()})
		return
	}

	h.logger.Debug("redirecting", slog.String("code", shortCode), slog.String("to", origLink.LongURL))
	c.Redirect(http.StatusFound, origLink.LongURL)
}
