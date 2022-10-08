// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"time"
)

type Board struct {
	ID        string    `json:"id" bson:"_id"`
	UserID    string    `json:"userId" bson:"userId"`
	Name      string    `json:"name" bson:"name"`
	Todos     []*Todo   `json:"todos" bson:"todos"`
	TodoIds   []*string `json:"todoIds" bson:"todoIds"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type GetBoardsRes struct {
	Boards []*Board `json:"boards" bson:"boards"`
	Cache  bool     `json:"cache" bson:"cache"`
}

type GetTodosRes struct {
	Todos []*Todo `json:"todos" bson:"todos"`
	Cache bool    `json:"cache" bson:"cache"`
}

type NewBoard struct {
	UserID string `json:"userId" bson:"userId"`
	Name   string `json:"name" bson:"name"`
}

type NewTodo struct {
	Text   string `json:"text" bson:"text"`
	UserID string `json:"userId" bson:"userId"`
}

type NewUser struct {
	Email    string  `json:"email" bson:"email"`
	Username *string `json:"username" bson:"username"`
}

type Todo struct {
	ID        string    `json:"id" bson:"_id"`
	UserID    string    `json:"userId" bson:"userId"`
	BoardID   string    `json:"boardId" bson:"boardId"`
	Text      string    `json:"text" bson:"text"`
	Priority  int       `json:"priority" bson:"priority"`
	Tag       string    `json:"tag" bson:"tag"`
	Markdown  bool      `json:"markdown" bson:"markdown"`
	Done      bool      `json:"done" bson:"done"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type UpdateBoard struct {
	ID     string    `json:"id" bson:"_id"`
	UserID string    `json:"userId" bson:"userId"`
	Name   string    `json:"name" bson:"name"`
	Todos  []*string `json:"todos" bson:"todos"`
}

type UpdateTodo struct {
	ID       string  `json:"id" bson:"_id"`
	UserID   string  `json:"userId" bson:"userId"`
	Text     *string `json:"text" bson:"text"`
	Priority *int    `json:"priority" bson:"priority"`
	Tag      *string `json:"tag" bson:"tag"`
	Markdown *bool   `json:"markdown" bson:"markdown"`
	Done     *bool   `json:"done" bson:"done"`
	BoardID  *string `json:"boardId" bson:"boardId"`
}

type User struct {
	ID       string    `json:"id" bson:"_id"`
	Username string    `json:"username" bson:"username"`
	Email    *string   `json:"email" bson:"email"`
	BoardIds []*string `json:"boardIds" bson:"boardIds"`
}
