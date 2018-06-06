package models

type RequestInstance struct {
	//The Url we are going to request
	Url string					`json:"requestUrl"`
	//Any headers needed for pass along.
	Headers []RequestHeader		`json:"headers"`
	//Method of the request [POST/GET/PUT/DELETE]
	Method	string				`json:"method"`
	//Body, an optional body
	Body	string				`json:"body"`
}

/**
Header object for http requests
 */
type RequestHeader struct {
	Key string		`json:"key"`
	Value string	`json:"value"`
}