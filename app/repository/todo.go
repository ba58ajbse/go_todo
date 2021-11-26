package repository

import (
	"database/sql"
	"go_todo/model"
)

type TodoRepository interface {
	GetAll(db sql.DB) ([]model.Todo, error)
	Get(db sql.DB, id int) (model.Todo, error)
	Create(db sql.DB, t model.Todo) (model.Todo, error)
	Update(db sql.DB, t model.Todo, id int) (model.Todo, int, error)
	Delete(db sql.DB, id int) (int, error)
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
		err := rows.Scan(&t.Id, &t.Todo, &t.Completed, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}

	return todos, err
}

// 1件取得
func (r *todoRepository) Get(db sql.DB, id int) (model.Todo, error) {
	todo := model.Todo{}
	err := db.QueryRow("SELECT * FROM todos WHERE id = ?", id).
		Scan(&todo.Id, &todo.Todo, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)

	return todo, err
}

// 新規作成
func (r *todoRepository) Create(db sql.DB, t model.Todo) (model.Todo, error) {
	tx, err := db.Begin()
	if err != nil {
		return t, err
	}
	defer tx.Rollback()

	stmt, err := db.Prepare("INSERT INTO todos (todo) VALUES (?)")
	if err != nil {
		return t, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(t.Todo)
	if err != nil {
		return t, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return t, err
	}
	t.Id = int(lastId)

	return t, err
}

// 更新処理
func (r *todoRepository) Update(db sql.DB, t model.Todo, id int) (model.Todo, int, error) {
	tx, err := db.Begin()
	if err != nil {
		return t, 0, err
	}
	defer tx.Rollback()

	stmt, err := db.Prepare("UPDATE todos SET todo = ?, completed = ? WHERE id = ?")
	if err != nil {
		return t, 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(t.Todo, t.Completed, id)
	if err != nil {
		return t, 0, err
	}

	rowCnt, err := res.RowsAffected()
	if err != nil {
		return t, 0, err
	}

	t.Id = id

	return t, int(rowCnt), err
}

func (r *todoRepository) Delete(db sql.DB, id int) (int, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	stmt, err := db.Prepare("DELETE FROM todos WHERE id = ?")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		return 0, err
	}

	rowCnt, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rowCnt), err
}
