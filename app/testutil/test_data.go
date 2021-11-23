package testutil

import "go_todo/model"

func GetTodoTestData() []model.Todo {
	todos := []model.Todo{
		{Id: 1, Todo: "JavaScript", Completed: false},
		{Id: 2, Todo: "TypeScript", Completed: false},
		{Id: 3, Todo: "PHP", Completed: false},
		{Id: 4, Todo: "Golang", Completed: false},
	}

	return todos
}
