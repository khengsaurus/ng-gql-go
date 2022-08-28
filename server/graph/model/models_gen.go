// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type NewTodo struct {
	Text   string `json:"text" bson:"text"`
	UserID string `json:"userId" bson:"objectId"`
}

type NewUser struct {
	Username string  `json:"username" bson:"username"`
	Email    *string `json:"email" bson:"email"`
}

type Todo struct {
	ID     string  `json:"id" bson:"_id"`
	UserID string  `json:"userId" bson:"objectId"`
	Text   string  `json:"text" bson:"text"`
	Color  *string `json:"color" bson:"color"`
	Done   bool    `json:"done" bson:"done"`
}

type User struct {
	ID       string  `json:"id" bson:"_id"`
	Username string  `json:"username" bson:"username"`
	Email    *string `json:"email" bson:"email"`
}
