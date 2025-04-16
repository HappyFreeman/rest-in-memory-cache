package service

import (
	"bytes"
	"encoding/json"
	"github.com/HappyFreeman/rest-in-memory-cache/internal/api/middleware"
	"github.com/HappyFreeman/rest-in-memory-cache/internal/config"
	"github.com/golang/mock/gomock"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"net/http/httptest"
	"testing"

	"github.com/HappyFreeman/rest-in-memory-cache/internal/dto"
	"github.com/HappyFreeman/rest-in-memory-cache/internal/repo"
	"github.com/HappyFreeman/rest-in-memory-cache/internal/repo/mocks"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

const (
	token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDQwMzEwNTIsImlkIjo1LCJuYW1lIjoicmFtemlrIn0.9guUFoovFplIeHGGZNQYzf_BC2FX6vyl41urZigdlhY"
)

func loadJWT(t *testing.T) config.JWT {
	if err := godotenv.Load("../../.env"); err != nil {
		require.NoError(t, err)
	}

	// Загружаем конфигурацию из переменных окружения
	var jwtCfg config.JWT
	if err := envconfig.Process("", &jwtCfg); err != nil {
		require.NoError(t, err)
	}

	return jwtCfg
}

// TestCreateTask - тестирование метода CreateTask
func TestCreateTask(t *testing.T) {

	jwtCfg := loadJWT(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish() // Завершаем мок после теста

	// Создаем мок репозитория
	mockRepo := mocks.NewMockRepository(ctrl)
	logger := zap.NewNop().Sugar() // Без вывода логов

	// Создаем экземпляр сервиса с мок-репозиторием
	s := NewService(mockRepo, logger)

	// Инициализируем Fiber-контекст
	app := fiber.New()
	//apiGroup := app.Group("/v1", middleware.Authorization(config.JWT{Secret: "secret"}))
	app.Post("/tasks", middleware.Authorization(config.JWT{Secret: jwtCfg.Secret}), s.CreateTask)

	t.Run("успешное создание задачи", func(t *testing.T) {
		task := TaskRequest{
			Title:       "Test Task",
			Description: "Test Description",
		}
		body, _ := json.Marshal(task)

		mockTask := repo.Task{
			Title:       task.Title,
			Description: task.Description,
		}

		// Ожидаем, что вызов метода `CreateTask` в репозитории вернёт ID = 1
		mockRepo.EXPECT().CreateTask(gomock.Any(), mockTask, gomock.Any()).Return(1, nil).Times(1)

		// Отправляем запрос
		req := httptest.NewRequest("POST", "/tasks", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		// Выполняем запрос
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		// Проверяем ответ
		var response dto.Response
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "success", response.Status)
	})

}
