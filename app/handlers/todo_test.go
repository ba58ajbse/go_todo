package handlers

import (
	"database/sql"
	"go_todo/testutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type TodoRepoStub struct{}

func (r *TodoRepoStub) getTodo(db sql.DB) {}

func TestGetAllTodos(t *testing.T) {
	e := echo.New()
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

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	NewTodoHandler(*db).GetAllTodos(c)

	expected, actual := testutil.FormatBodyForTest(data, rec.Body)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, expected, actual)
}

func TestGetTodo(t *testing.T) {
	e := echo.New()
	db, mock := testutil.GetMockDB()
	defer db.Close()

	data := testutil.GetTodoTestData()[0]
	row := mock.NewRows([]string{"id", "name", "email", "created_at", "updated_at"}).
		AddRow(data.Id, data.Todo, data.Completed, data.CreatedAt, data.UpdatedAt)
	paramValue := strconv.Itoa(data.Id)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM todos WHERE id = ?")).
		WithArgs(data.Id).
		WillReturnRows(row)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/todos/:id")
	c.SetParamNames("id")
	c.SetParamValues(paramValue)

	NewTodoHandler(*db).GetTodo(c)

	expected, actual := testutil.FormatBodyForTest(data, rec.Body)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, expected, actual)
}

func TestCreateTodo(t *testing.T) {
	e := echo.New()
	db, mock := testutil.GetMockDB()
	defer db.Close()

	todo := testutil.GetTodoTestData()[1]
	lastId := todo.Id + 1
	todo.Id = lastId
	reqBody := testutil.FormatModelDataToJsonStr(todo)

	mock.ExpectBegin()
	mock.ExpectPrepare(regexp.QuoteMeta("INSERT INTO todos (todo) VALUES (?)")).
		ExpectExec().
		WithArgs(todo.Todo).
		WillReturnResult(sqlmock.NewResult(int64(lastId), 1))
	mock.ExpectCommit()

	req := httptest.NewRequest(http.MethodPost, "/todos", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	NewTodoHandler(*db).CreateTodo(c)

	recBody := testutil.RemoveLFForRecBody(rec.Body)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, reqBody, recBody)
}

func TestUpdateTodo(t *testing.T) {
	e := echo.New()
	db, mock := testutil.GetMockDB()
	defer db.Close()

	todo := testutil.GetTodoTestData()[2]
	id := todo.Id
	reqBody := testutil.FormatModelDataToJsonStr(todo)

	mock.ExpectBegin()
	mock.ExpectPrepare(regexp.QuoteMeta("UPDATE todos SET todo = ?, completed = ? WHERE id = ?")).
		ExpectExec().
		WithArgs(todo.Todo, todo.Completed, todo.Id).
		WillReturnResult(sqlmock.NewResult(int64(todo.Id), 1))
	mock.ExpectCommit()

	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/todos/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))

	NewTodoHandler(*db).UpdateTodo(c)

	recBody := testutil.RemoveLFForRecBody(rec.Body)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.Equal(t, "1", recBody)
}

func TestDeleteTodo(t *testing.T) {
	e := echo.New()
	db, mock := testutil.GetMockDB()
	defer db.Close()

	delId := 2

	mock.ExpectBegin()
	mock.ExpectPrepare(regexp.QuoteMeta("DELETE FROM todos WHERE id = ?")).
		ExpectExec().
		WithArgs(delId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/todos/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(delId))

	NewTodoHandler(*db).DeleteTodo(c)

	recBody := testutil.RemoveLFForRecBody(rec.Body)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.Equal(t, "1", recBody)
}
