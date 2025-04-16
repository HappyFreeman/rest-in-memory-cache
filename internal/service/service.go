package service

import (
	"github.com/HappyFreeman/rest-in-memory-cache/internal/dto"
	"github.com/HappyFreeman/rest-in-memory-cache/internal/repo"
	"github.com/HappyFreeman/rest-in-memory-cache/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// Слой бизнес-логики. Тут должна быть основная логика сервиса

// Service - интерфейс для бизнес-логики
type Service interface {
	CreateTask(ctx *fiber.Ctx) error
	GetTask(ctx *fiber.Ctx) error
	GetTasks(ctx *fiber.Ctx) error
	DeleteTask(ctx *fiber.Ctx) error
	UpdateTask(ctx *fiber.Ctx) error
}

type service struct {
	repo repo.Repository
	log  *zap.SugaredLogger
}

// NewService - конструктор сервиса
func NewService(repo repo.Repository, logger *zap.SugaredLogger) Service {
	return &service{
		repo: repo,
		log:  logger,
	}
}

// CreateTask - обработчик запроса на создание задачи
func (s *service) CreateTask(ctx *fiber.Ctx) error {
	var req TaskRequest

	// Десериализация JSON-запроса json.Unmarshal(ctx.Body(), &req)
	if err := ctx.BodyParser(&req); err != nil {
		s.log.Error("Invalid request body", zap.Error(err))
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "Invalid request body")
	}

	// Валидация входных данных
	if vErr := validator.Validate(ctx.Context(), req); vErr != nil {
		return dto.BadResponseError(ctx, dto.FieldIncorrect, vErr.Error())
	}

	userId, ok := ctx.Locals("userId").(int64)

	if !ok {
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "Invalid token")
	}

	taskID, err := s.repo.CreateTask(ctx.Context(), repo.CreateTaskParams{
		UserID:      userId,
		Title:       req.Title,
		Description: req.Description,
	})

	if err != nil {
		s.log.Error("Failed to insert task", zap.Error(err))
		return dto.InternalServerError(ctx)
	}

	// Формирование ответа
	response := dto.Response{
		Status: "success",
		Data: fiber.Map{
			"task_id": taskID,
		},
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (s *service) GetTask(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")

	if err != nil {
		s.log.Error("Invalid request body", zap.Error(err))
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "Invalid request body")
	}

	userId, ok := ctx.Locals("userId").(int64)

	if !ok {
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "Invalid token")
	}

	task, err := s.repo.GetTaskById(ctx.Context(), repo.GetTaskByIdParams{
		TaskID: int64(id),
		UserID: userId,
	})

	if err != nil {
		s.log.Error("Failed to get task", zap.Error(err))
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "Task not found")
	}

	response := dto.Response{
		Status: "success",
		Data:   task,
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

// GetTasks - обработчик запроса на получение всех задач
// TODO: Добавить пагинацию
func (s *service) GetTasks(ctx *fiber.Ctx) error {

	userId, ok := ctx.Locals("userId").(int64)

	if !ok {
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "Invalid token")
	}

	tasks, err := s.repo.GetTasks(ctx.Context(), repo.GetTasksParams{
		UserID: userId,
	})

	if err != nil {
		s.log.Error("Failed to get tasks", zap.Error(err))
		return dto.InternalServerError(ctx)
	}

	response := dto.Response{
		Status: "success",
		Data:   tasks,
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (s *service) DeleteTask(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")

	if err != nil {
		s.log.Error("Invalid request body", zap.Error(err))
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "Invalid request body")
	}

	userId, ok := ctx.Locals("userId").(int64)

	if !ok {
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "Invalid token")
	}

	err = s.repo.DeleteTaskById(ctx.Context(), repo.DeleteTaskByIdParams{
		TaskID: int64(id),
		UserID: userId,
	})

	if err != nil {
		s.log.Error("Failed to delete task", zap.Error(err))
		return dto.InternalServerError(ctx)
	}

	response := dto.Response{
		Status: "success",
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (s *service) UpdateTask(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")

	if err != nil {
		s.log.Error("Invalid request body", zap.Error(err))
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "Invalid request body")
	}

	userId, ok := ctx.Locals("userId").(int64)

	if !ok {
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "Invalid token")
	}

	var req TaskRequest

	// Десериализация JSON-запроса json.Unmarshal(ctx.Body(), &req)
	if err := ctx.BodyParser(&req); err != nil {
		s.log.Error("Invalid request body", zap.Error(err))
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "Invalid request body")
	}

	// Валидация входных данных
	if vErr := validator.Validate(ctx.Context(), req); vErr != nil {
		return dto.BadResponseError(ctx, dto.FieldIncorrect, vErr.Error())
	}

	err = s.repo.UpdateTaskById(ctx.Context(), repo.UpdateTaskParams{
		TaskID:      int64(id),
		UserID:      userId,
		Title:       req.Title,
		Description: req.Description,
	})

	if err != nil {
		s.log.Error("Failed to update task", zap.Error(err))
		return dto.InternalServerError(ctx)
	}

	response := dto.Response{
		Status: "success",
		Data: repo.Task{
			Title:       req.Title,
			Description: req.Description,
		},
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}
