package models

/**
A request object
 */
type AggRequest struct {
	//Requests that this service should make.
	Requests []RequestInstance 	`json:"requests"`
	//The number of ms to wait before giving up on the aggregation.
	Timeout int 					`json:"timeout"`
	//The number of goroutines to run
	Concurrency int 				`json:"concurrency"`
}