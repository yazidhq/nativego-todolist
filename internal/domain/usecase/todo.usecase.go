package usecase

import (
	"context"
	"restapi-native-go/internal/domain/entity"
	"restapi-native-go/internal/domain/repository"
	"restapi-native-go/internal/utils/errors"
	"time"
)

type TodoUseCase interface {
	ListTodos(ctx context.Context) ([]*entity.Todo, error)
	GetTodo(ctx context.Context, id int64) (*entity.Todo, error)
	CreateTodo(ctx context.Context, title, description string) (*entity.Todo, error)
	UpdateTodo(ctx context.Context, id int64, title, description string, completed bool) (*entity.Todo, error)
	DeleteTodo(ctx context.Context, id int64) error
}

type todoUseCase struct {
	todoRepo repository.TodoRepository
	timeout  time.Duration
}

func NewTodoUseCase(todoRepo repository.TodoRepository, timeout time.Duration) TodoUseCase {
	return &todoUseCase{
		todoRepo: todoRepo,
		timeout:  timeout,
	}
}

func (uc *todoUseCase) ListTodos(ctx context.Context) ([]*entity.Todo, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	return uc.todoRepo.GetAll(ctx)
}

func (uc *todoUseCase) GetTodo(ctx context.Context, id int64) (*entity.Todo, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	todo, err := uc.todoRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.NewNotFoundError("Todo not found.")
	}

	return todo, nil
}

func (uc *todoUseCase) CreateTodo(ctx context.Context, title, description string) (*entity.Todo, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	if title == "" {
		return nil, errors.NewBadRequestError("Title cannot be empty or null.")
	}

	now := time.Now()
	todo := &entity.Todo{
		Title:       title,
		Description: description,
		Completed:   false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := uc.todoRepo.Create(ctx, todo); err != nil {
		return nil, errors.NewInternalServerError("Failed to craete todo.")
	}

	return todo, nil
}

func (uc *todoUseCase) UpdateTodo(ctx context.Context, id int64, title, description string, completed bool) (*entity.Todo, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	todo, err := uc.todoRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.NewNotFoundError("Todo not found.")
	}

	if title != "" {
		todo.Title = title
	}

	if description != "" {
		todo.Description = description
	}

	now := time.Now()
	todo.Completed = completed
	todo.UpdatedAt = now

	if err := uc.todoRepo.Update(ctx, todo); err != nil {
		return nil, errors.NewInternalServerError("Failed to update todo.")
	}

	return todo, nil
}

func (uc *todoUseCase) DeleteTodo(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	if _, err := uc.todoRepo.GetByID(ctx, id); err != nil {
		return errors.NewNotFoundError("Todo not found.")
	}

	if err := uc.todoRepo.Delete(ctx, id); err != nil {
		return errors.NewInternalServerError("Failed to delete todo.")
	}

	return nil
}
