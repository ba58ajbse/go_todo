package handlers

import (
	"database/sql"
	"encoding/json"
	"go_todo/testutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type TodoRepoStub struct{}

func (r *TodoRepoStub) getTodo(db sql.DB) {}

func TestGetAllTodos(t *testing.T) {
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

	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	NewTodoHandler(*db).GetAllTodos(c)

	res := rec.Result()
	defer res.Body.Close()

	out, _ := json.Marshal(data)
	expected := string(out)
	actual := strings.TrimRight(rec.Body.String(), "\n")

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, expected, actual)

}
