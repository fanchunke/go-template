package model

type Book struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}
