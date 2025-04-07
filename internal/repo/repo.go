package repo

import (
	"context"
	"fmt"
	"github.com/HappyFreeman/rest-in-memory-cache/internal/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

// Слой репозитория, здесь должны быть все методы, связанные с базой данных

// SQL-запрос на вставку задачи
const (
	insertTaskQuery = `INSERT INTO tasks (title, description, user_id) VALUES ($1, $2, $3) RETURNING id;`
	getTaskQuery    = `SELECT title, description FROM tasks WHERE id = $1 and user_id = $2;`
	deleteTaskQuery = `DELETE FROM tasks WHERE id = $1 and user_id = $2;`
	updateTaskQuery = `UPDATE tasks SET title = $1, description = $2 WHERE id = $3 and user_id = $4;`
	getTasksQuery   = `SELECT title, description FROM tasks WHERE user_id = $1;`
)

type repository struct {
	pool *pgxpool.Pool
}

// Repository - интерфейс с методом создания задачи
// mockgen -source=C:/MyProjects/rest-in-memory-cache/internal/repo/repo.go -destination=C:/MyProjects/rest-in-memory-cache/internal/repo/mocks/repository.go -package=mocks
type Repository interface {
	CreateTask(ctx context.Context, task Task, userId int) (int, error)
	GetTaskById(ctx context.Context, id int, userId int) (Task, error)
	DeleteTaskById(ctx context.Context, id int, userId int) error
	UpdateTaskById(ctx context.Context, id int, task Task, userId int) error
	GetTasks(ctx context.Context, userId int) ([]Task, error)
}

// NewRepository - создание нового экземпляра репозитория с подключением к PostgreSQL
func NewRepository(ctx context.Context, cfg config.PostgreSQL) (Repository, error) {
	// Формируем строку подключения
	connString := fmt.Sprintf(
		`user=%s password=%s host=%s port=%d dbname=%s sslmode=%s 
        pool_max_conns=%d pool_max_conn_lifetime=%s pool_max_conn_idle_time=%s`,
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
		cfg.PoolMaxConns,
		cfg.PoolMaxConnLifetime.String(),
		cfg.PoolMaxConnIdleTime.String(),
	)

	// Парсим конфигурацию подключения
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse PostgreSQL config")
	}

	// Оптимизация выполнения запросов (кеширование запросов)
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeCacheDescribe

	// Создаём пул соединений с базой данных
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create PostgreSQL connection pool")
	}

	return &repository{pool}, nil
}

// CreateTask - вставка новой задачи в таблицу tasks
func (r *repository) CreateTask(ctx context.Context, task Task, userId int) (int, error) {
	var id int
	err := r.pool.QueryRow(ctx, insertTaskQuery, task.Title, task.Description, userId).Scan(&id)
	if err != nil {
		return 0, errors.Wrap(err, "failed to insert task")
	}
	return id, nil
}

// GetTaskById - получение задачи по ее id
func (r *repository) GetTaskById(ctx context.Context, id int, userId int) (Task, error) {
	var task Task

	err := r.pool.QueryRow(ctx, getTaskQuery, id, userId).Scan(&task.Title, &task.Description)

	if err != nil {
		return Task{}, errors.Wrap(err, "failed to get task")
	}

	return task, nil
}

// DeleteTaskById - удаление задачи по ее id
func (r *repository) DeleteTaskById(ctx context.Context, id int, userId int) error {
	_, err := r.pool.Exec(ctx, deleteTaskQuery, id, userId)
	if err != nil {
		return errors.Wrap(err, "failed to delete task")
	}
	return nil
}

// UpdateTaskById - обновление задачи по ее id
func (r *repository) UpdateTaskById(ctx context.Context, id int, task Task, userId int) error {
	_, err := r.pool.Exec(ctx, updateTaskQuery, task.Title, task.Description, id, userId)
	if err != nil {
		return errors.Wrap(err, "failed to update task")
	}
	return nil
}

// GetTasks - получение всех задач
// TODO: Реализовать пагинацию
func (r *repository) GetTasks(ctx context.Context, userId int) ([]Task, error) {
	var tasks []Task

	rows, err := r.pool.Query(ctx, getTasksQuery, userId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get tasks")
	}
	defer rows.Close()

	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.Title, &task.Description); err != nil {
			return nil, errors.Wrap(err, "failed to scan task")
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "failed to get tasks")
	}

	return tasks, nil
}
