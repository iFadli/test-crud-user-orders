package template

type ResponseHTTP struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Error   error       `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Page    interface{} `json:"page,omitempty"`
}

type PagePagination struct {
	Limit int64 `json:"limit"`
	Page  int64 `json:"page"`
	Show  int   `json:"show"`
	Total int64 `json:"total,omitempty"`
}
