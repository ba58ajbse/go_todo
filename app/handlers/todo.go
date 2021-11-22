package handlers

import (
	"database/sql"
	"go_todo/repository"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type TodoHandler interface {
	GetAllTodos(c echo.Context) error
	GetTodo(c echo.Context) error
}

type todoHandler struct {
	db       sql.DB
	todoRepo repository.TodoRepository
}

func NewTodoHandler(db sql.DB) TodoHandler {
	return &todoHandler{
		db:       db,
		todoRepo: repository.NewTodoRepository(),
	}
}

func (h *todoHandler) GetAllTodos(c echo.Context) error {
	todos, err := h.todoRepo.GetAll(h.db)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, todos)
}

func (h *todoHandler) GetTodo(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	todo, err := h.todoRepo.Get(h.db, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, todo)
}
