package testutil

import "go_todo/model"

func GetTodoTestData() []model.Todo {
	todos := []model.Todo{
		{Id: 1, Todo: "JavaScript", Completed: false, CreatedAt: "2021-11-26 11:42:48", UpdatedAt: "2021-11-26 11:42:48"},
		{Id: 2, Todo: "TypeScript", Completed: false, CreatedAt: "2021-11-26 11:42:48", UpdatedAt: "2021-11-26 11:42:48"},
		{Id: 3, Todo: "PHP", Completed: false, CreatedAt: "2021-11-26 11:42:48", UpdatedAt: "2021-11-26 11:42:48"},
		{Id: 4, Todo: "Golang", Completed: false, CreatedAt: "2021-11-26 11:42:48", UpdatedAt: "2021-11-26 11:42:48"},
	}

	return todos
}
