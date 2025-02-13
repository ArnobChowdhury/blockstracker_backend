package models

type CreateTaskRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
}

type UpdateTaskRequest struct {
	TaskID      int    `json:"task_id" binding:"required"`
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
}

type UpdateRepetitiveTaskRequest struct {
	TaskID    int    `json:"task_id" binding:"required"`
	Frequency string `json:"frequency" binding:"required"`
}
