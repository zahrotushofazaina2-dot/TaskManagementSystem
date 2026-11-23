package dto

type UserDashboardResponse struct {
	Total           int64 `json:"total_tasks"`
	TotalTodo       int64 `json:"todo"`
	TotalInProgress int64 `json:"in_progress"`
	TotalDone       int64 `json:"done"`
}
