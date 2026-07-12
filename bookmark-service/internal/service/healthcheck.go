package service

import (
	"context"
	"fmt"
	"time"

	"github.com/khaivutri/bookmark-service/internal/model"
	"github.com/khaivutri/bookmark-service/internal/repository"
)

// HealthCheck defines the interface for health check service.
type HealthCheck interface {
	Check(ctx context.Context) (*model.HealthReport, error)
}
type redisPinger interface {
	Ping(ctx context.Context) error
}

type healthCheck struct {
	serviceName 	string
	instanceID  	string
	
	redisAdapter 	redisPinger
}

// NewHealthCheck creates and returns a new HealthCheck service instance.
func NewHealthCheck(serviceName, instanceID string, redisAdapter redisPinger) HealthCheck {
	return &healthCheck{	
							serviceName: 	serviceName,
							instanceID:		instanceID,
							redisAdapter: 	redisAdapter,
						}
}

// Check returns the health status report for the service.
func (hc *healthCheck) Check(ctx context.Context) (*model.HealthReport, error) {
	dpc := make(map[string]string)

	var msg string 
	var firstErr error

	pingCtx, cancel := context.WithTimeout(ctx, 5 * time.Second)
	defer cancel()

	if pingErr := hc.redisAdapter.Ping(pingCtx); pingErr != nil {
		dpc["redis"] = "DOWN"
		msg = "DEGRADED"
		firstErr = fmt.Errorf("%w: redis: %v", repository.ErrDependencyDown, pingErr)
	} else {
		dpc["redis"] = "UP"
		msg = "OK"
	}

	return &model.HealthReport{	Message: msg, 
								ServiceName: hc.serviceName, 
								InstanceID: hc.instanceID, 
								Dependencies: dpc,
							}, firstErr

}