package model

type User struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

type UserCache struct {
	ID   string `redis:"id"`
	Name string `redis:"name"`
}
