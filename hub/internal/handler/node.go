package handler

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kanaya/jobboard-hub/internal/database/repo"
	"github.com/kanaya/jobboard-hub/internal/middleware"
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
	ID           int64     `json:"id"`
	NodeName     string    `json:"node_name"`
	CurrentJobID *int64    `json:"current_job_id"`
	CreatedAt    time.Time `json:"created_at"`
	NodeToken    string    `json:"node_token,omitempty"`
}

func nodeToResponse(node repo.Node) nodeResponse {
	var createdAt time.Time
	if node.CreatedAt.Valid {
		createdAt = node.CreatedAt.Time
	}
	return nodeResponse{
		ID:           node.ID,
		NodeName:     node.NodeName,
		CurrentJobID: node.CurrentJobID,
		CreatedAt:    createdAt,
	}
}

func (h *NodeHandler) List(c *gin.Context) {
	clusterID := c.GetString(middleware.ClusterIDContextKey)
	nodes, err := h.queries.ListNodesByCluster(c.Request.Context(), clusterID)
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

type createNodeResponse struct {
	nodeResponse
	NodeToken string `json:"node_token"`
}

func (h *NodeHandler) Create(c *gin.Context) {
	var req createNodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	clusterID := c.GetString(middleware.ClusterIDContextKey)
	nodeToken, err := generateNodeToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate node token"})
		return
	}

	node, err := h.queries.CreateNode(c.Request.Context(), repo.CreateNodeParams{
		ClusterID:     clusterID,
		NodeName:      req.NodeName,
		NodeTokenHash: hashNodeToken(nodeToken),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create node"})
		return
	}

	c.JSON(http.StatusOK, createNodeResponse{
		nodeResponse: nodeToResponse(node),
		NodeToken:    nodeToken,
	})
}

func (h *NodeHandler) Delete(c *gin.Context) {
	clusterID := c.GetString(middleware.ClusterIDContextKey)
	nodeID, err := strconv.ParseInt(c.Param("node_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid node id"})
		return
	}

	rows, err := h.queries.DeleteNodeByCluster(c.Request.Context(), repo.DeleteNodeByClusterParams{
		ID:        nodeID,
		ClusterID: clusterID,
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

func generateNodeToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func hashNodeToken(nodeToken string) string {
	sum := sha256.Sum256([]byte(nodeToken))
	return hex.EncodeToString(sum[:])
}
