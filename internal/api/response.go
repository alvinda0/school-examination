package api

type Metadata struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type Response struct {
	Status   int         `json:"status"`
	Success  bool        `json:"success"`
	Message  string      `json:"message"`
	Data     interface{} `json:"data"`
	Metadata *Metadata   `json:"metadata,omitempty"`
}
