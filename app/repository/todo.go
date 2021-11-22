package repository

import (
	"database/sql"
	"go_todo/model"
)

type TodoRepository interface {
	GetAll(db sql.DB) ([]model.Todo, error)
}
type todoRepository struct{}

func NewTodoRepository() TodoRepository {
	return &todoRepository{}
}

// 全件取得
func (r *todoRepository) GetAll(db sql.DB) ([]model.Todo, error) {
	todos := []model.Todo{}
	rows, err := db.Query("SELECT * FROM todos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		t := model.Todo{}
		err := rows.Scan(&t.Id, &t.Todo, &t.Completed)
		if err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}

	return todos, err
}
