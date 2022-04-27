package v1

type EndpointStat struct {
	Url   string  `json:"url"`
	Views int32   `json:"views"`
	Score float32 `json:"relevanceScore"`
}

type Stats struct {
	Data []EndpointStat `json:"data"`
}

type APIResponse struct {
	Stats `json:",inline"`
	Count int `json:"count"`
}

type APIErrorResponse struct {
	Message   string `json:"message"`
	ErrorCode int    `json:"error"`
}
