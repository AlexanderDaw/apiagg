package transport

import (
	"github.com/AlexanderDaw/apiagg/services"
	"github.com/AlexanderDaw/apiagg/models"
	"github.com/go-kit/kit/endpoint"
	"errors"
	"context"
)

type (
	Endpoints struct {
		GetAggregationEndpoint endpoint.Endpoint
		GetTestServiceEndpoint endpoint.Endpoint
	}
)


//Error for issues deserializing the requst for search.
var ErrTypeAssertQueryGetAggregationRequest = errors.New("type assertion failed on AggregationEndpoint")
var ErrTypeAssertTriggerRequest = errors.New("type assertion failed on trigger request")

/**
Setup endpoints
 */
func MakeEndpoints(svc services.AggregationService, testSvc services.VariableResponseTestService) Endpoints {

	return Endpoints{
		GetAggregationEndpoint: MakeAggregationEndpoint(svc),
		GetTestServiceEndpoint: MakeSlowEndpoint(testSvc),
	}

}

/**
Main endpoint for search aggregations.
 */
func MakeAggregationEndpoint(svc services.AggregationService) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (res interface{}, err error){

		aggR, ok := req.(*models.AggRequest)
		if !ok {
			return nil, ErrTypeAssertQueryGetAggregationRequest
		}

		aggResp := svc.FanOutRequests(ctx, *aggR)

		return aggResp, nil
	}
}

/**
Test endpoint for variable service response
 */
func MakeSlowEndpoint(svc services.VariableResponseTestService) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (res interface{}, err error){

		slowR, ok := req.(*models.TriggerRequest)
		if !ok {
			return nil, ErrTypeAssertTriggerRequest
		}
		varResp := svc.VariableResponse(ctx, *slowR)

		return varResp, nil
	}

}


