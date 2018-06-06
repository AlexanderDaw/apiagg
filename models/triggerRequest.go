package models

/**
Request in object for variable response time handleing.
 */
type TriggerRequest struct {
	ResponseTime int `json:"responseTime"`
}