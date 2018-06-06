package services

import (
	"time"
	"github.com/AlexanderDaw/apiagg/models"
	"context"
)

type(
	VariableResponseTestService interface {
		VariableResponse(context.Context, models.TriggerRequest) *models.TriggerResponse
	}

	VariableResponseTestServiceInstance struct {
		ResponseTimeDelayMs int
	}
)

func InitializeVariableResponseTestService(Delay int) (*VariableResponseTestServiceInstance, error){
	VRT := VariableResponseTestServiceInstance{
		ResponseTimeDelayMs: Delay,
	}
	return &VRT, nil
}

func (service *VariableResponseTestServiceInstance) VariableResponse(ctx context.Context, request models.TriggerRequest) *models.TriggerResponse{

	sleepTime := service.ResponseTimeDelayMs
	if request.ResponseTime != 0 {
		sleepTime = request.ResponseTime
	}

	time.Sleep(time.Duration(sleepTime)*time.Millisecond)
	tr := models.TriggerResponse{Response:"Complete"}
	return &tr
}