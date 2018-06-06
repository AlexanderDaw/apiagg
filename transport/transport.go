package transport

import (
	"github.com/gorilla/mux"
	httptransport "github.com/go-kit/kit/transport/http"
	mainhttp "net/http"
	"github.com/AlexanderDaw/apiagg/services"
	"github.com/AlexanderDaw/apiagg/models"
	"encoding/json"
	"context"
	"log"
)

/**
Match the routes to the endpoints
 */
func RegisterRoutes(r *mux.Router, svc services.AggregationService, testSvc services.VariableResponseTestService) {

	e := MakeEndpoints(svc, testSvc)

	r.Handle("/1.0/aggregate/auto",
		httptransport.NewServer(
			e.GetAggregationEndpoint,
			decodeRequest,
			encodeJSONResponse))


	r.Handle("/1.0/varresp/simulate",
		httptransport.NewServer(
			e.GetTestServiceEndpoint,
			decodeTestRequest,
			encodeJSONResponse))

}


func decodeRequest(_ context.Context, req *mainhttp.Request) (interface{}, error) {

	request := models.AggRequest{}

	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		log.Println("unable to decode json")
		return nil, err
	}
	return &request, nil
}

/**
Decode the test request.
 */
func decodeTestRequest(_ context.Context, req *mainhttp.Request) (interface{}, error){
	request := models.TriggerRequest{}

	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		log.Println("unable to decode test json ",err.Error())
		return nil, err
	}

	return &request, nil
}

func encodeJSONResponse(_ context.Context, w mainhttp.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
