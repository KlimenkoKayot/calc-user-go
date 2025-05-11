package orchestrator

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	service "github.com/klimenkokayot/calc-net-go/internal/orchestrator/service"
	"github.com/klimenkokayot/calc-net-go/internal/shared/models"
	"github.com/klimenkokayot/calc-net-go/internal/shared/utils"
	"github.com/klimenkokayot/calc-user-go/config"
)

type Expressions struct {
	List []models.Expression `json:"expressions"`
}

// Структура обработчика, требует новый сервис
type OrchestratorHandler struct {
	Service *service.OrchestratorService
	config  *config.Config
}

// Создает экземпляр обработчика
func NewOrchestratorHandler(config *config.Config, service *service.OrchestratorService) *OrchestratorHandler {
	return &OrchestratorHandler{
		Service: service,
		config:  config,
	}
}

// Обработка получения новой задачи в оркестратор
func (h *OrchestratorHandler) NewExpression(w http.ResponseWriter, r *http.Request) {
	// Считываем тело запроса
	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.Writer(w).Write(utils.ErrorResponse(ErrInternalServer))
		return
	}
	defer r.Body.Close()

	// Создаем экземпляр выражения, пытаемся распарсить
	expression := &models.Expression{}
	err = json.Unmarshal(data, expression)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		io.Writer(w).Write(utils.ErrorResponse(ErrInvalidBodyDecode))
		return
	}

	// Попытка добавления новой задачи в сервис
	hash, err := h.Service.AddExpression(expression.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.Writer(w).Write(utils.ErrorResponse(ErrInternalServer))
		return
	}

	// Под id мы берем хэш, полученный SHA512 от арифметического выражения
	json, err := json.Marshal(models.Expression{
		Id: utils.EncodeToString(hash),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.Writer(w).Write(utils.ErrorResponse(ErrInternalServer))
		return
	}
	w.WriteHeader(http.StatusCreated)
	io.Writer(w).Write(json)
}

// Обработка запросов на получения списка всех полученных и обработанных задач
func (h *OrchestratorHandler) Expressions(w http.ResponseWriter, r *http.Request) {
	// Запрос в сервис

	expressions := h.Service.GetAllExpressions()

	json, err := json.Marshal(expressions)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.Writer(w).Write(utils.ErrorResponse(ErrInternalServer))
		return
	}

	w.WriteHeader(http.StatusOK)
	io.Writer(w).Write(json)
}

// Обработка запроса на получения статуса конкретного выражения
func (h *OrchestratorHandler) Expression(w http.ResponseWriter, r *http.Request) {
	// id это hash выражения в формате base64
	id := mux.Vars(r)["id"]

	// base64 конвертируем в SHA512
	hash, err := utils.EncodedToSHA512(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.Writer(w).Write(utils.ErrorResponse(err))
		return
	}

	// Пытаемся найти выражение в списке необработанных выражений

	if val, found := h.Service.Expressions[hash]; found {
		json, _ := json.Marshal(models.Expression{
			Id:     id,
			Status: val.Status,
		})
		w.WriteHeader(http.StatusOK)
		io.Writer(w).Write(json)

		return
	}

	// Пытаемся найти выражение в списке обработанных выражений
	if val, found := h.Service.Answers[hash]; found {
		json, _ := json.Marshal(models.Expression{
			Id:     id,
			Status: "Выполнено.",
			Result: val,
		})
		w.WriteHeader(http.StatusOK)
		io.Writer(w).Write(json)

		return
	}

	w.WriteHeader(http.StatusNotFound)
}

// Для агента, обработка запроса на получение новой задачи
func (h *OrchestratorHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	// Попытка получения новой подзадачи из сервиса
	task, err := h.Service.GetTask()
	if err == service.ErrHaveNoTask {
		w.WriteHeader(http.StatusNotFound)
		io.Writer(w).Write(utils.ErrorResponse(err))
		return
	} else if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		io.Writer(w).Write(utils.ErrorResponse(err))
		return
	}
	data, _ := json.Marshal(task)
	w.WriteHeader(http.StatusOK)
	io.Writer(w).Write(data)
}

// Для агента, обработка запроса на решение подзадачи
func (h *OrchestratorHandler) PostTask(w http.ResponseWriter, r *http.Request) {
	// Читаем тело запроса
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Ошибка при получении результата task: %e\n", err)
		return
	}
	// Экземпляр для парсинга тела запроса
	taskAnswer := &models.TaskResult{}
	err = json.Unmarshal(data, taskAnswer)
	if err != nil {
		log.Printf("Ошибка при попытке парсинга json TaskAnswer: %e\n", err)
		return
	}
	if taskAnswer.Error != "" {
		log.Printf("Подзадача id %d, error: %s\n", taskAnswer.Id, taskAnswer.Error)

		h.Service.ProcessErrorAnswer(taskAnswer)

		return
	}
	log.Printf("Обработка ответа id: %d, ответ: %f", taskAnswer.Id, taskAnswer.Result)
	// Обрабатываем решение подзадачи в сервисе
	h.Service.ProcessAnswer(taskAnswer)
}

func (h *OrchestratorHandler) Index(w http.ResponseWriter, r *http.Request) {
	temp := template.Must(template.ParseFiles(filepath.Join(".", "web", "template", "index.html")))
	w.WriteHeader(http.StatusFound)
	w.Header().Set("Content-Type", "text/html")
	temp.Execute(w, nil)
}
