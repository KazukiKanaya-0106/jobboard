package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kanaya/jobboard-hub/internal/database/repo"
	"github.com/kanaya/jobboard-hub/internal/middleware"
)

type ClusterHandler struct {
	queries repo.Querier
}

func NewClusterHandler(queries repo.Querier) *ClusterHandler {
	return &ClusterHandler{
		queries: queries,
	}
}

type clusterResponse struct {
	Id        string    `json:"cluster_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (h *ClusterHandler) Me(c *gin.Context) {
	clusterId := c.GetString(middleware.ClusterIdContextKey)
	cluster, err := h.queries.GetCluster(c.Request.Context(), clusterId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load cluster"})
		return
	}

	var createdAt time.Time
	if cluster.CreatedAt.Valid {
		createdAt = cluster.CreatedAt.Time
	}

	c.JSON(http.StatusOK, clusterResponse{
		Id:        cluster.ID,
		CreatedAt: createdAt,
	})
}
