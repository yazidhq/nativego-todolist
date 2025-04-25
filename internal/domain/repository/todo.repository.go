package repository

import (
	"context"
	"restapi-native-go/internal/domain/entity"
)

type TodoRepository interface {
	GetAll(ctx context.Context) ([]*entity.Todo, error)
	GetByID(ctx context.Context, id int64) (*entity.Todo, error)
	Create(ctx context.Context, todo *entity.Todo) error
	Update(ctx context.Context, todo *entity.Todo) error
	Delete(ctx context.Context, id int64) error
}
