package agent

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/klimenkokayot/calc-net-go/internal/shared/models"
)

// Получает новые задания по адресу оркестратора и клиенту-воркеру
func GetTask(client *http.Client, url string) (*models.Task, error) {
	// Создаем новый запрос
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	// Отправляем запрос
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	// Проверка статус кода
	switch resp.StatusCode {
	case http.StatusNotFound:
		return nil, ErrNotFound
	case http.StatusInternalServerError:
		return nil, ErrInternalServer
	case http.StatusOK:
		break
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
func PostTask(client *http.Client, url string, result *models.TaskResult) error {
	data, _ := json.Marshal(result)
	// Создание нового запроса
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return err
	}
	_, err = client.Do(req)
	return err
}
