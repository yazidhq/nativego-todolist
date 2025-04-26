package handler

import (
	"encoding/json"
	"net/http"
	"restapi-native-go/internal/domain/usecase"
	"restapi-native-go/internal/utils/errors"
	"restapi-native-go/pkg/response"
	"strconv"
	"strings"
)

type TodoHandler struct {
	todoUseCase usecase.TodoUseCase
}

func NewTodoHandler(todoUseCase usecase.TodoUseCase) *TodoHandler {
	return &TodoHandler{
		todoUseCase: todoUseCase,
	}
}

type createTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type updateTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

func (h *TodoHandler) ListTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := h.todoUseCase.ListTodos(r.Context())
	if err != nil {
		response.Error(w, err)
	}

	response.JSON(w, http.StatusOK, todos)
}

func (h *TodoHandler) GetTodo(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/todos/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, errors.NewBadRequestError("Invalid ID."))
		return
	}

	todo, err := h.todoUseCase.GetTodo(r.Context(), id)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, todo)
}

func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var req createTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, errors.NewBadRequestError("Invalid JSON format."))
		return
	}

	todo, err := h.todoUseCase.CreateTodo(r.Context(), req.Title, req.Description)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, todo)
}

func (h *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/todos/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, errors.NewBadRequestError("Invalid ID."))
		return
	}

	var req updateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, errors.NewBadRequestError("Invalid JSON format."))
		return
	}

	todo, err := h.todoUseCase.UpdateTodo(r.Context(), id, req.Title, req.Description, req.Completed)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, todo)
}

func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/todos/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, errors.NewBadRequestError("Invalid ID."))
		return
	}

	if err := h.todoUseCase.DeleteTodo(r.Context(), id); err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}
