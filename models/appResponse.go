package models

type AggResponse struct {
	Responses []RequestResponse	`json:"responses"`
	Metrics	[]ResponseMetrics	`json:"metrics"`
}