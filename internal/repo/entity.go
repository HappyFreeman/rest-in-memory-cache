package repo

// Task - структура, соответствующая таблице tasks
type Task struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CreateTaskParams struct {
	UserID      int64
	Title       string
	Description string
}

type GetTaskByIdParams struct {
	TaskID int64
	UserID int64
}

type UpdateTaskParams struct {
	TaskID      int64
	UserID      int64
	Title       string
	Description string
}

type GetTasksParams struct {
	UserID int64
}

type DeleteTaskByIdParams struct {
	TaskID int64
	UserID int64
}
