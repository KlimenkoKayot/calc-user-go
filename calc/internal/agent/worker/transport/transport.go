package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/klimenkokayot/calc-net-go/internal/shared/models"
)

// Получает новые задания по адресу оркестратора и клиенту-воркеру
func GetTask(client *http.Client, baseURL string) (*models.Task, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("неверный базовый URL: %w", err)
	}
	u.Path = path.Join(u.Path, "internal/task")

	// Создаем новый запрос
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("Ошибка при создании запроса: %w", err)
	}
	// Отправляем запрос
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Ошибка при отправке запроса: %w", err)
	}
	// Проверка статус кода
	switch resp.StatusCode {
	case http.StatusNotFound:
		return nil, ErrNotFound
	case http.StatusInternalServerError:
		return nil, ErrInternalServer
	case http.StatusOK:
		break
	case http.StatusUnauthorized:
		return nil, ErrNotAuth
	default:
		return nil, ErrInternalServer
	}
	defer resp.Body.Close()
	// Парсинг данных
	body, err := io.ReadAll(resp.Body)
	task := &models.Task{}
	err = json.Unmarshal(body, task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

// Возвращаем результат решенной подзадачи в оркестратор
func PostTask(client *http.Client, baseURL string, result *models.TaskResult) error {
	u, err := url.Parse(baseURL)
	if err != nil {
		return fmt.Errorf("неверный базовый URL: %w", err)
	}
	u.Path = path.Join(u.Path, "internal/task")

	data, _ := json.Marshal(result)
	// Создание нового запроса
	req, err := http.NewRequest("POST", u.String(), bytes.NewReader(data))
	if err != nil {
		return err
	}
	_, err = client.Do(req)
	return err
}
