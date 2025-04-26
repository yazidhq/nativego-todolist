package postgres

import (
	"context"
	"database/sql"
	"log"
	"restapi-native-go/internal/domain/entity"
	"restapi-native-go/internal/domain/repository"
)

type todoRepository struct {
	db *sql.DB
}

func NewTodoRepository(db *sql.DB) repository.TodoRepository {
	return &todoRepository{
		db: db,
	}
}

func (r *todoRepository) GetAll(ctx context.Context) ([]*entity.Todo, error) {
	query := `SELECT * FROM todos ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		log.Printf("Error querying todo: %v", err)
		return nil, err
	}
	defer rows.Close()

	var todos []*entity.Todo
	for rows.Next() {
		var todo entity.Todo
		if err := rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Description,
			&todo.Completed,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		); err != nil {
			log.Printf("Error scanning todo row: %v", err)
			return nil, err
		}
		todos = append(todos, &todo)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating todo rows: %v", err)
		return nil, err
	}

	return todos, nil
}

func (r *todoRepository) GetByID(ctx context.Context, id int64) (*entity.Todo, error) {
	query := `SELECT * FROM todos WHERE id = $1`

	row := r.db.QueryRowContext(ctx, query)

	var todo entity.Todo
	err := row.Scan(
		&todo.ID,
		&todo.Title,
		&todo.Description,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		log.Printf("Error querying todoL %v", err)
		return nil, err
	}

	return &todo, nil
}

func (r *todoRepository) Create(ctx context.Context, todo *entity.Todo) error {
	query := `INSERT INTO todos (title, description, completed, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := r.db.QueryRowContext(
		ctx,
		query,
		todo.Title,
		todo.Description,
		todo.Completed,
		todo.CreatedAt,
		todo.UpdatedAt,
	).Scan(&todo.ID)

	if err != nil {
		log.Printf("Error creating todo: %v", err)
		return err
	}

	return nil
}

func (r *todoRepository) Update(ctx context.Context, todo *entity.Todo) error {
	query := `UPDATE todos SET title = $1, description = $2, completed = $3, updated_at = $4 WHERE id = $5`

	res, err := r.db.ExecContext(
		ctx,
		query,
		todo.Title,
		todo.Description,
		todo.Completed,
		todo.UpdatedAt,
		todo.ID,
	)

	if err != nil {
		log.Printf("Error updating todo: %v", err)
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *todoRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM todos WHERE id = $1`

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		log.Printf("Error deleting todo: %v", err)
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
