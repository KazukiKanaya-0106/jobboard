package handler

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kanaya/jobboard-hub/internal/apierror"
	"github.com/kanaya/jobboard-hub/internal/database/repo"
)

type JobTriggerHandler struct {
	queries repo.Querier
}

func NewJobTriggerHandler(queries repo.Querier) *JobTriggerHandler {
	return &JobTriggerHandler{
		queries: queries,
	}
}

type startJobRequest struct {
	NodeToken string     `json:"node_token" binding:"required"`
	Tag       *string    `json:"tag"`
	StartedAt *time.Time `json:"started_at"`
}

type finishJobRequest struct {
	NodeToken     string     `json:"node_token" binding:"required"`
	Status        string     `json:"status" binding:"omitempty,oneof=completed failed"`
	FinishedAt    *time.Time `json:"finished_at"`
	DurationHours *float64   `json:"duration_hours"`
	ErrorText     *string    `json:"error_text"`
}

type JobTriggerResponse struct {
	Success bool `json:"success"`
}

func (h *JobTriggerHandler) StartJob(c *gin.Context) {
	var req startJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apierror.Write(c, apierror.InvalidRequest)
		return
	}

	node, ok := h.getNodeByNodeToken(c, req.NodeToken)
	if !ok {
		return
	}

	if node.CurrentJobID != nil {
		apierror.Write(c, apierror.JobAlreadyRunning)
		return
	}

	var started any
	if req.StartedAt != nil {
		started = timestamptz(*req.StartedAt)
	}

	job, err := h.queries.CreateJob(c.Request.Context(), repo.CreateJobParams{
		ClusterID: node.ClusterID,
		NodeID:    node.ID,
		Column3:   started,
		Column4:   nil,
		Tag:       req.Tag,
	})
	if err != nil {
		log.Printf("failed to create job: %v", err)
		apierror.Write(c, apierror.Internal)
		return
	}

	_, err = h.queries.UpdateNodeCurrentJob(c.Request.Context(), repo.UpdateNodeCurrentJobParams{
		ID:           node.ID,
		CurrentJobID: &job.ID,
	})
	if err != nil {
		log.Printf("failed to update node state: %v", err)
		apierror.Write(c, apierror.Internal)
		return
	}

	c.JSON(http.StatusCreated, JobTriggerResponse{Success: true})
}

func (h *JobTriggerHandler) FinishJob(c *gin.Context) {
	var req finishJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apierror.Write(c, apierror.InvalidRequest)
		return
	}

	node, ok := h.getNodeByNodeToken(c, req.NodeToken)
	if !ok {
		return
	}

	if node.CurrentJobID == nil {
		apierror.Write(c, apierror.JobNotRunning)
		return
	}

	finishedAt := time.Now()
	if req.FinishedAt != nil {
		finishedAt = req.FinishedAt.UTC()
	} else {
		finishedAt = finishedAt.UTC()
	}

	status := req.Status
	if status == "" {
		status = "completed"
	}

	var duration pgtype.Interval
	if req.DurationHours != nil {
		duration = intervalFromHours(*req.DurationHours)
	}

	_, err := h.queries.UpdateJob(c.Request.Context(), repo.UpdateJobParams{
		ID:            *node.CurrentJobID,
		StartedAt:     pgtype.Timestamptz{},
		FinishedAt:    timestamptz(finishedAt),
		Status:        status,
		DurationHours: duration,
		ErrorText:     req.ErrorText,
	})
	if err != nil {
		log.Printf("failed to update job: %v", err)
		apierror.Write(c, apierror.Internal)
		return
	}

	_, err = h.queries.UpdateNodeCurrentJob(c.Request.Context(), repo.UpdateNodeCurrentJobParams{
		ID:           node.ID,
		CurrentJobID: nil,
	})
	if err != nil {
		log.Printf("failed to reset node state: %v", err)
		apierror.Write(c, apierror.Internal)
		return
	}

	c.JSON(http.StatusOK, JobTriggerResponse{Success: true})
}

func (h *JobTriggerHandler) getNodeByNodeToken(c *gin.Context, secret string) (repo.Node, bool) {
	node, err := h.queries.GetNodeByNodeTokenHash(c.Request.Context(), hashNodeToken(secret))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			apierror.Write(c, apierror.NodeNotFound)
			return repo.Node{}, false
		}
		log.Printf("failed to load node: %v", err)
		apierror.Write(c, apierror.Internal)
		return repo.Node{}, false
	}
	return node, true
}

func timestamptz(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{
		Time:  t.UTC(),
		Valid: true,
	}
}

func intervalFromHours(hours float64) pgtype.Interval {
	us := int64(hours * float64(time.Hour/time.Microsecond))
	return pgtype.Interval{
		Microseconds: us,
		Valid:        true,
	}
}

func intervalToHours(iv pgtype.Interval) *float64 {
	if !iv.Valid {
		return nil
	}
	hours := float64(iv.Microseconds) / float64(time.Hour/time.Microsecond)
	return &hours
}
