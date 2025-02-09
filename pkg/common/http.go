package common

type HTTPError struct {
	Error string `json:"error"`
	Type  string `json:"type"`
} // @name HTTPError

type Paging struct {
	Page  int `json:"page,omitempty" form:"page"`
	Limit int `json:"limit,omitempty" form:"limit"`
	// Total *int `json:"total,omitempty"`
} // @name Paging

func (p *Paging) Process() {
	if p.Page <= 0 {
		p.Page = 1
	}

	if p.Limit <= 0 {
		p.Limit = 10
	}
}

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Paging  Paging      `json:"paging,omitempty"`
} // @name Response

func SuccessResponse(data interface{}) Response {
	return Response{
		Success: true,
		Data:    data,
	}
}

func SuccessResponseWithPaging(data interface{}, paging Paging) Response {
	return Response{
		Success: true,
		Data:    data,
		Paging:  paging,
	}
}
