package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kanaya/jobboard-hub/internal/apierror"
	"github.com/kanaya/jobboard-hub/internal/database/repo"
	"github.com/kanaya/jobboard-hub/internal/middleware"
)

type JobHandler struct {
	queries repo.Querier
}

func NewJobHandler(queries repo.Querier) *JobHandler {
	return &JobHandler{
		queries: queries,
	}
}

type jobResponse struct {
	ID            int64      `json:"id"`
	ClusterID     string     `json:"cluster_id"`
	NodeID        int64      `json:"node_id"`
	Status        string     `json:"status"`
	StartedAt     *time.Time `json:"started_at,omitempty"`
	FinishedAt    *time.Time `json:"finished_at,omitempty"`
	DurationHours *float64   `json:"duration_hours,omitempty"`
	Tag           *string    `json:"tag,omitempty"`
	ErrorText     *string    `json:"error_text,omitempty"`
}

func (h *JobHandler) List(c *gin.Context) {
	clusterID := c.GetString(middleware.ClusterIDContextKey)

	jobs, err := h.queries.ListJobsByCluster(c.Request.Context(), clusterID)
	if err != nil {
		log.Printf("failed to list jobs: %v", err)
		apierror.Write(c, apierror.Internal)
		return
	}

	resp := make([]jobResponse, 0, len(jobs))
	for _, job := range jobs {
		resp = append(resp, jobToResponse(job))
	}

	c.JSON(200, resp)
}

func (h *JobHandler) Get(c *gin.Context) {
	clusterID := c.GetString(middleware.ClusterIDContextKey)
	jobID, err := strconv.ParseInt(c.Param("job_id"), 10, 64)
	if err != nil {
		apierror.Write(c, apierror.InvalidRequest)
		return
	}

	job, err := h.queries.GetJobByClusterAndJobID(c.Request.Context(), repo.GetJobByClusterAndJobIDParams{
		ClusterID: clusterID,
		ID:        jobID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			apierror.Write(c, apierror.JobNotFound)
			return
		}
		log.Printf("failed to load job: %v", err)
		apierror.Write(c, apierror.Internal)
		return
	}

	c.JSON(http.StatusOK, jobToResponse(job))
}

func (h *JobHandler) ListByNode(c *gin.Context) {
	clusterID := c.GetString(middleware.ClusterIDContextKey)
	nodeID, err := strconv.ParseInt(c.Param("node_id"), 10, 64)
	if err != nil {
		apierror.Write(c, apierror.InvalidRequest)
		return
	}

	nodes, err := h.queries.ListNodesByCluster(c.Request.Context(), clusterID)
	if err != nil {
		log.Printf("failed to load nodes: %v", err)
		apierror.Write(c, apierror.Internal)
		return
	}

	owned := false
	for _, node := range nodes {
		if node.ID == nodeID {
			owned = true
			break
		}
	}

	if !owned {
		apierror.Write(c, apierror.NodeNotFound)
		return
	}

	jobs, err := h.queries.ListJobsByNode(c.Request.Context(), nodeID)
	if err != nil {
		log.Printf("failed to list jobs: %v", err)
		apierror.Write(c, apierror.Internal)
		return
	}

	resp := make([]jobResponse, 0, len(jobs))
	for _, job := range jobs {
		resp = append(resp, jobToResponse(job))
	}

	c.JSON(http.StatusOK, resp)
}

func timestamptzPtr(ts pgtype.Timestamptz) *time.Time {
	if !ts.Valid {
		return nil
	}
	return &ts.Time
}

func jobToResponse(job repo.Job) jobResponse {
	return jobResponse{
		ID:            job.ID,
		ClusterID:     job.ClusterID,
		NodeID:        job.NodeID,
		Status:        job.Status,
		StartedAt:     timestamptzPtr(job.StartedAt),
		FinishedAt:    timestamptzPtr(job.FinishedAt),
		DurationHours: intervalToHours(job.DurationHours),
		Tag:           job.Tag,
		ErrorText:     job.ErrorText,
	}
}
