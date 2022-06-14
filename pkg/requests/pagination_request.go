package requests

// PaginationRequest 分页请求
type PaginationRequest struct {
	Request
	Page int `json:"page" form:"page" valid:"page"`
	PerPage int `json:"per_page" form:"per_page" valid:"per_page"`
}