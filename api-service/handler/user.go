package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	proto "testTwoServices/proto/sonyyy04.user.v1"
)

type Handler struct {
	client proto.UserServiceClient
	logger *zap.SugaredLogger
}

func NewHandler(client proto.UserServiceClient, logger *zap.SugaredLogger) *Handler {
	return &Handler{
		client: client,
		logger: logger,
	}
}

func (h *Handler) CreateUser(c *gin.Context) {
	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := c.BindJSON(&req); err != nil {
		h.logger.Errorw("Invalid input", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	ctx := c.Request.Context()
	resp, err := h.client.CreateUser(ctx, &proto.CreateUserRequest{
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		h.logger.Errorw("gRPC CreateUser failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}
	h.logger.Infow("User created", "id", resp.Id)
	c.JSON(http.StatusOK, gin.H{"id": resp.Id})
}

func (h *Handler) GetUser(c *gin.Context) {
	id := c.Param("id")
	ctx := c.Request.Context()
	resp, err := h.client.GetUser(ctx, &proto.GetUserRequest{
		Id: id,
	})
	if err != nil {
		h.logger.Errorw("gRPC GetUser failed", "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "failed to get user"})
		return
	}
	h.logger.Infow("User got", "id", resp.Id, "name", resp.Name, "email", resp.Email)
	c.JSON(http.StatusOK, &gin.H{
		"id":    resp.Id,
		"name":  resp.Name,
		"email": resp.Email,
	})
}
