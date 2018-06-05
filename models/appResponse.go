package model

type AggResponse struct {
	Responses []RequestResponse	`json:"responses"`
	Metrics	[]ResponseMetrics	`json:"metrics"`
}