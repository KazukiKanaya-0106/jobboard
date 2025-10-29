package handler

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kanaya/jobboard-hub/internal/database/repo"
	"github.com/kanaya/jobboard-hub/internal/middleware"
	"golang.org/x/crypto/bcrypt"
)

type NodeHandler struct {
	queries repo.Querier
}

func NewNodeHandler(queries repo.Querier) *NodeHandler {
	return &NodeHandler{
		queries: queries,
	}
}

type nodeResponse struct {
	Id           int64     `json:"id"`
	NodeName     string    `json:"node_name"`
	CurrentJobId *int64    `json:"current_job_id"`
	CreatedAt    time.Time `json:"created_at"`
}

func nodeToResponse(node repo.Node) nodeResponse {
	var createdAt time.Time
	if node.CreatedAt.Valid {
		createdAt = node.CreatedAt.Time
	}
	return nodeResponse{
		Id:           node.ID,
		NodeName:     node.NodeName,
		CurrentJobId: node.CurrentJobID,
		CreatedAt:    createdAt,
	}
}

func (h *NodeHandler) List(c *gin.Context) {
	clusterId := c.GetString(middleware.ClusterIdContextKey)
	nodes, err := h.queries.ListNodesByCluster(c.Request.Context(), clusterId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list nodes"})
		return
	}
	resp := make([]nodeResponse, 0, len(nodes))
	for _, node := range nodes {
		resp = append(resp, nodeToResponse(node))
	}
	c.JSON(http.StatusOK, resp)
}

type createNodeRequest struct {
	NodeName string `json:"node_name" binding:"required"`
}

func (h *NodeHandler) Create(c *gin.Context) {
	var req createNodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	clusterId := c.GetString(middleware.ClusterIdContextKey)
	webhookSecret, err := generateWebhookSecret()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate webhook secret"})
		return
	}

	secretHash, err := bcrypt.GenerateFromPassword([]byte(webhookSecret), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash webhook secret"})
		return
	}

	node, err := h.queries.CreateNode(c.Request.Context(), repo.CreateNodeParams{
		ClusterID:         clusterId,
		NodeName:          req.NodeName,
		WebhookSecretHash: string(secretHash),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create node"})
		return
	}

	c.JSON(http.StatusOK, nodeToResponse(node))
}

func (h *NodeHandler) Delete(c *gin.Context) {
	clusterId := c.GetString(middleware.ClusterIdContextKey)
	nodeId, err := strconv.ParseInt(c.Param("node_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid node id"})
		return
	}

	rows, err := h.queries.DeleteNodeByCluster(c.Request.Context(), repo.DeleteNodeByClusterParams{
		ID:        nodeId,
		ClusterID: clusterId,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete node"})
		return
	}

	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "node not found"})
		return
	}

	c.Status(http.StatusNoContent)
}

func generateWebhookSecret() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}
