package handlers

import (
	"database/sql"
	"go_todo/model"
	"go_todo/repository"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type TodoHandler interface {
	GetAllTodos(c echo.Context) error
	GetTodo(c echo.Context) error
	CreateTodo(c echo.Context) error
	UpdateTodo(c echo.Context) error
	DeleteTodo(c echo.Context) error
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

func (h *todoHandler) CreateTodo(c echo.Context) error {
	t := new(model.Todo)
	if err := c.Bind(t); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	todo, err := h.todoRepo.Create(h.db, *t)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusCreated, todo)
}

func (h *todoHandler) UpdateTodo(c echo.Context) error {
	t := new(model.Todo)
	id, _ := strconv.Atoi(c.Param("id"))
	if err := c.Bind(t); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	todo, rowCnt, err := h.todoRepo.Update(h.db, *t, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if rowCnt == 0 {
		return c.JSON(http.StatusConflict, todo)
	}

	return c.JSON(http.StatusNoContent, todo)
}

func (h *todoHandler) DeleteTodo(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	rowCnt, err := h.todoRepo.Delete(h.db, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if rowCnt == 0 {
		return c.JSON(http.StatusConflict, rowCnt)
	}

	return c.JSON(http.StatusNoContent, rowCnt)
}
