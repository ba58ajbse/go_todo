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
	rows := mock.NewRows([]string{"id", "todo", "completed"}).
		AddRow(data[0].Id, data[0].Todo, data[0].Completed).
		AddRow(data[1].Id, data[1].Todo, data[1].Completed).
		AddRow(data[2].Id, data[2].Todo, data[2].Completed).
		AddRow(data[3].Id, data[3].Todo, data[3].Completed)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM todos")).
		WillReturnRows(rows)

	res, err := NewTodoRepository().GetAll(*db)

	assert.Equal(t, nil, err)
	assert.Equal(t, len(data), len(res))
}

func TestGet(t *testing.T) {
	db, mock := testutil.GetMockDB()
	defer db.Close()

	id := 1
	row := mock.NewRows([]string{"id", "name", "email"}).
		AddRow(1, "JavaScript", false)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM todos WHERE id = ?")).
		WithArgs(id).
		WillReturnRows(row)

	res, err := NewTodoRepository().Get(*db, id)

	assert.Equal(t, nil, err)
	assert.Equal(t, id, res.Id)
	assert.Equal(t, "JavaScript", res.Todo)
	assert.Equal(t, false, res.Completed)
}

func TestCreate(t *testing.T) {
	db, mock := testutil.GetMockDB()
	defer db.Close()

	lastId := 2
	todo := model.Todo{Id: lastId, Todo: "Ruby", Completed: false}

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

	todo := model.Todo{Id: 2, Todo: "Ruby", Completed: true}

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
