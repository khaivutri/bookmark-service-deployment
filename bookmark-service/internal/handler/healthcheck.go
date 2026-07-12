package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/khaivutri/bookmark-service/internal/service"
	"github.com/khaivutri/bookmark-service/pkg/response"
	"github.com/rs/zerolog/log"
)

// HealthCheck defines the interface for health check handler.
type HealthCheck interface {
	HealthCheck(ctx *gin.Context)
}

type healthCheck struct {
	IsHealthy service.HealthCheck
}

// NewHealthCheck creates and returns a new HealthCheck handler.
func NewHealthCheck(isHealthy service.HealthCheck) HealthCheck {
	return &healthCheck{IsHealthy: isHealthy}
}

// HealthCheck handles the health check HTTP request.
//@Summary Health Check
//@Tags Health Check
//@Accept json
//@Produce json
//@Success 200 {object} model.HealthReport
//@Failure 503 {object} model.HealthReport
//@Router /health-check [get]
func (hc *healthCheck) HealthCheck(ctx *gin.Context) {
	report, err := hc.IsHealthy.Check(ctx)

	if report == nil {
		log.Error().Err(err).Str("from", "handler.healthCheck.HealthCheck").Msg("failed to check health")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.InternalServerErrResponse)
		return
	}

	status := http.StatusOK
	if report.Message != "OK" {
		status = http.StatusServiceUnavailable

		log.Error().
			Err(err).
			Str("from", "handler.healthCheck.HealthCheck").
			Str("report_message", report.Message).
			Msg("health check degraded")
	}
	ctx.JSON(status, report)
}
