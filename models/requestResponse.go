package model

type RequestResponse struct {
	Url	string				`json:"url"`
	ResponseHeaders []ResponseHeader 	`json:"responseHeaders"`
	ResponseBody	string 			`json:"responseBody"`
	StatusCode	int			`json:"statusCode"`
}

type ResponseHeader struct {
	Key   string        `json:"key"`
	Value []string        `json:"value"`
}

/**
Object used to track the speed of the requests down stream in the stack.
 */
type ResponseMetrics struct {
	RequestUrl string	`json:"requestUrl"`
	ContentLength int64	`json:"contentLength"`
	ResponseTimeMs int	`json:"responseTimeMs"`
	StatusCode int	`json:"StatusCode"`
	TransferEncoding []string `json:transferEncoding`
}