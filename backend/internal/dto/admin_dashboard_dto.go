package dto

type AdminDashboardResponse struct {
	Users      UserDashboardStats     `json:"users"`
	Tasks      TaskDashboardStats     `json:"tasks"`
	Categories CategoryDashboardStats `json:"categories"`
}

type UserDashboardStats struct {
	Total    int64 `json:"total"`
	Active   int64 `json:"active"`
	Inactive int64 `json:"inactive"`
}

type TaskDashboardStats struct {
	Total      int64             `json:"total"`
	ByStatus   TaskStatusStats   `json:"by_status"`
	ByPriority TaskPriorityStats `json:"by_priority"`
}

type TaskStatusStats struct {
	Todo       int64 `json:"todo"`
	InProgress int64 `json:"in_progress"`
	Done       int64 `json:"done"`
}

type TaskPriorityStats struct {
	Low    int64 `json:"low"`
	Medium int64 `json:"medium"`
	High   int64 `json:"high"`
}

type CategoryDashboardStats struct {
	Total int64 `json:"total"`
}
