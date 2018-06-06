/**
Service responsible for taking one or more REST requests and running them
concurrently, then streaming their responses back on a common channel.
 */
package services

import (
	"github.com/AlexanderDaw/apiagg/models"
	"net/http"
	"log"
	"io/ioutil"
	"context"
	"sync"
	"bytes"
	"time"
	"strconv"
	"net/url"
)

/**
aggregation of call from the backend.
 */
type(

	AggregationService interface {
		FanOutRequests(context.Context, models.AggRequest) *models.AggResponse
	}

	AggregationServiceInstance struct {
		//Primary settings here.
		threadCount int
	}
)

/**
Spin up the Aggregation service
 */
func InitializeAggregationService(aggThreadCount int) (*AggregationServiceInstance, error) {

	asi := AggregationServiceInstance{
		threadCount: aggThreadCount,
	}

	return &asi, nil
}

/**
Given a request fan it out.
 */
func (service *AggregationServiceInstance) FanOutRequests(cxt context.Context, aggReq models.AggRequest) *models.AggResponse {

	totalRequests := 0
	var requesters []chan models.RequestInstance
	var responders []<-chan models.RequestsPipelineInstance
	//All responders will funnel to this chan.
	responseChannel := make(chan models.RequestsPipelineInstance, 100)

	threadConcurrency := service.threadCount

	if aggReq.Concurrency != 0 {
		threadConcurrency = aggReq.Concurrency
	}

	log.Println("Querying with "+strconv.Itoa(threadConcurrency)+" threads")

	//Spin up the request pool.
	for i:=0; i< threadConcurrency; i++{
		newRequesterChan := make(chan models.RequestInstance, 10)
		newResponseChan := make(chan models.RequestsPipelineInstance, 10)
		requesters = append(requesters, newRequesterChan)
		responders = append(responders, newResponseChan)
		go service.RequestConcurrent(newRequesterChan, newResponseChan)
	}

	//Setup consumers for this service.
	go combine(responders, responseChannel)

	//task distribution
	for _, request := range aggReq.Requests{
		totalRequests+=1
		requesters[totalRequests%len(requesters)] <- request
	}

	//Clean up
	for _, chanInst := range requesters{
		close(chanInst)
	}

	//Setup time gate
	timeGateChan := make(chan bool)
	go service.RequestTimeLatch(aggReq.Timeout, timeGateChan)

	responseObj := models.AggResponse{}
	//loop on top level response.
completeResponseCollection:
	for{
		select{
		case _, ok := <- timeGateChan:
			if !ok {
				log.Println("Some Requests Timed out, returning now. ")
				//Insert stubs for uncompleted work.
				var urlResponded bool
				for _, request := range aggReq.Requests{
					urlResponded = false
					for _, response := range responseObj.Metrics {
						if request.Url == response.RequestUrl{
							urlResponded = true
							break
						}

					}
					if urlResponded == false{
						responseObj.Metrics = append(responseObj.Metrics, models.ResponseMetrics{
							StatusCode: -1,
							RequestUrl: request.Url,
							ResponseTimeMs: aggReq.Timeout,
						} )
					}
				}

				break completeResponseCollection
			}
		case response, ok := <- responseChannel:
			if !ok {
				break completeResponseCollection
			}
			responseObj.Responses = append(responseObj.Responses, response.Response)
			responseObj.Metrics = append(responseObj.Metrics, response.ResponseMetrics)
		}
	}


	return &responseObj

}

func (service AggregationServiceInstance) RequestTimeLatch(timeout int, timeChan chan bool){
	time.AfterFunc(time.Duration(timeout)*time.Millisecond, func(){
		if timeChan != nil {
			close(timeChan)
		}
	})
}

func (service AggregationServiceInstance) RequestConcurrent(requestChan chan models.RequestInstance, responseChan chan models.RequestsPipelineInstance ) {

requestThreadComplete:
	for {
		select {
		case instance, ok := <- requestChan:
			if !ok {
				break requestThreadComplete
			}

			//Santize url.
			requestUrl, err := url.Parse(instance.Url)
			if err!= nil {
				log.Println(err.Error())
				PiplineResponse := models.RequestsPipelineInstance{
					ResponseMetrics:models.ResponseMetrics{
						RequestUrl:instance.Url,
						ResponseTimeMs: 0,
						StatusCode: -1,
					},
					Response:models.RequestResponse{
						Url:instance.Url,
						StatusCode: -1,
						ResponseBody: err.Error(),
					},
				}
				responseChan <- PiplineResponse
				continue
			}


			client := &http.Client{}
			log.Println("Requesting "+instance.Url)
			req, err := http.NewRequest(instance.Method, requestUrl.String(), bytes.NewBuffer([]byte(instance.Body)))

			if err != nil {
				log.Printf("Error creating HTTP Requst for the Entity Search for Query :: " + instance.Url)
			}

			for _, header := range (instance.Headers) {
				req.Header.Add(header.Key, header.Value)
			}

			//Start Timer
			startTime := time.Now().UnixNano() / int64(time.Millisecond)
			log.Println("STARTTIME == "+strconv.Itoa(int(startTime)))
			resp, err := client.Do(req)
			endTime := time.Now().UnixNano() / int64(time.Millisecond)
			log.Println("ENDTIME == "+strconv.Itoa(int(endTime)))
			responseTime := int((endTime - startTime))
			//log.Println(resp.Status)
			if err!= nil {
				log.Println(err.Error())
				PiplineResponse := models.RequestsPipelineInstance{
					ResponseMetrics:models.ResponseMetrics{
						RequestUrl:instance.Url,
						ResponseTimeMs: responseTime,
						StatusCode: -1,
					},
					Response:models.RequestResponse{
						Url:instance.Url,
						StatusCode: -1,
						ResponseBody: err.Error(),
					},
				}
				responseChan <- PiplineResponse
			}else {
				requestResponse, _ := ioutil.ReadAll(resp.Body)

				reqResp := models.RequestResponse{
					Url:          instance.Url,
					ResponseBody: string(requestResponse),
					StatusCode: resp.StatusCode,
				}

				/**
					Iterate over the response headers and append them to the aggregated request.
				*/
				for k, v := range resp.Header {
					respHeader := models.ResponseHeader{Key: k, Value: v }
					reqResp.ResponseHeaders = append(reqResp.ResponseHeaders, respHeader)
				}

				resp.Body.Close()

				if err != nil {
					log.Println("Error getting data from service  for query : ", instance.Url)
				}

				responseMetrics := models.ResponseMetrics{
					ResponseTimeMs:   responseTime,
					ContentLength:    resp.ContentLength,
					StatusCode:       resp.StatusCode,
					TransferEncoding: resp.TransferEncoding,
					RequestUrl: instance.Url,
				}

				pipelineResponse := models.RequestsPipelineInstance{
					Response:        reqResp,
					ResponseMetrics: responseMetrics,
				}

				responseChan <- pipelineResponse
			}
		}
	}
	close(responseChan)
}

//Given a bunch of input channels multiplex them togther.
func combine(inputs []<-chan models.RequestsPipelineInstance, output chan<- models.RequestsPipelineInstance) {
	var group sync.WaitGroup
	for i := range inputs {
		group.Add(1)
		//Closure has access to the group scope.
		go func(input <-chan models.RequestsPipelineInstance) {
			for val := range input {
				output <- val
			}
			group.Done()
		}(inputs[i])
	}
	go func() {
		group.Wait()
		close(output)
	}()

}