package route

import (
	"database/sql"
	"go_todo/handlers"
	"net/http"

	"github.com/labstack/echo/v4"
)

func NewRouter(e *echo.Echo, db *sql.DB) {
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "Hello, World!")
	})

	h := handlers.NewTodoHandler(*db)

	e.GET("/todos", h.GetAllTodos)
	e.GET("/todos/:id", h.GetTodo)
	e.POST("/todos", h.CreateTodo)
	e.PUT("/todos/:id", h.UpdateTodo)
	e.DELETE("/todos/:id", h.DeleteTodo)
}
