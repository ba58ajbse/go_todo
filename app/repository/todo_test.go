package repository

import (
	"go_todo/model"
	"go_todo/testutil"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetAll(t *testing.T) {
	db, mock := testutil.GetMockDB()
	defer db.Close()

	data := testutil.GetTodoTestData()
	rows := mock.NewRows([]string{"id", "todo", "completed", "created_at", "updated_at"}).
		AddRow(data[0].Id, data[0].Todo, data[0].Completed, data[0].CreatedAt, data[0].UpdatedAt).
		AddRow(data[1].Id, data[1].Todo, data[1].Completed, data[1].CreatedAt, data[1].UpdatedAt).
		AddRow(data[2].Id, data[2].Todo, data[2].Completed, data[2].CreatedAt, data[2].UpdatedAt).
		AddRow(data[3].Id, data[3].Todo, data[3].Completed, data[3].CreatedAt, data[3].UpdatedAt)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM todos")).
		WillReturnRows(rows)

	res, err := NewTodoRepository().GetAll(*db)

	assert.Equal(t, nil, err)
	assert.Equal(t, len(data), len(res))
}

func TestGet(t *testing.T) {
	db, mock := testutil.GetMockDB()
	defer db.Close()

	data := testutil.GetTodoTestData()[0]
	row := mock.NewRows([]string{"id", "name", "email", "created_at", "updated_at"}).
		AddRow(data.Id, data.Todo, data.Completed, data.CreatedAt, data.UpdatedAt)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM todos WHERE id = ?")).
		WithArgs(data.Id).
		WillReturnRows(row)

	res, err := NewTodoRepository().Get(*db, data.Id)

	assert.Equal(t, nil, err)
	assert.Equal(t, data.Id, res.Id)
	assert.Equal(t, data.Todo, res.Todo)
	assert.Equal(t, data.Completed, res.Completed)
	assert.Equal(t, data.CreatedAt, res.CreatedAt)
	assert.Equal(t, data.UpdatedAt, res.UpdatedAt)
}

func TestCreate(t *testing.T) {
	db, mock := testutil.GetMockDB()
	defer db.Close()

	lastId := 2
	todo := model.Todo{Id: lastId, Todo: "Ruby", Completed: false, CreatedAt: "2021-04-12 12:04:45", UpdatedAt: "2021-04-12 12:04:45"}

	mock.ExpectBegin()
	mock.ExpectPrepare(regexp.QuoteMeta("INSERT INTO todos (todo) VALUES (?)")).
		ExpectExec().
		WithArgs(todo.Todo).
		WillReturnResult(sqlmock.NewResult(int64(lastId), 1))
	mock.ExpectCommit()

	res, err := NewTodoRepository().Create(*db, todo)

	assert.Equal(t, nil, err)
	assert.Equal(t, todo, res)
}

func TestUpdate(t *testing.T) {
	db, mock := testutil.GetMockDB()
	defer db.Close()

	todo := model.Todo{Id: 2, Todo: "Ruby", Completed: true, CreatedAt: "2021-04-12 12:04:45", UpdatedAt: "2021-04-12 20:04:45"}

	mock.ExpectBegin()
	mock.ExpectPrepare(regexp.QuoteMeta("UPDATE todos SET todo = ?, completed = ? WHERE id = ?")).
		ExpectExec().
		WithArgs(todo.Todo, todo.Completed, todo.Id).
		WillReturnResult(sqlmock.NewResult(int64(todo.Id), 1))
	mock.ExpectCommit()

	res, rowCnt, err := NewTodoRepository().Update(*db, todo, 2)

	assert.Equal(t, nil, err)
	assert.Equal(t, 1, rowCnt)
	assert.Equal(t, todo, res)
}

func TestDelete(t *testing.T) {
	db, mock := testutil.GetMockDB()
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectPrepare(regexp.QuoteMeta("DELETE FROM todos WHERE id = ?")).
		ExpectExec().
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	rowCnt, err := NewTodoRepository().Delete(*db, 1)

	assert.Equal(t, nil, err)
	assert.Equal(t, 1, rowCnt)
}
