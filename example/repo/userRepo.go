package repo

import "github.com/jmoiron/sqlx"

//go:generate go run generator.go

type UserRepo interface {
	UserCRUD
}

type UserRepoImpl struct {
	db *sqlx.DB
}

func MakeUserRepo(db *sqlx.DB) UserRepo {
	return &UserRepoImpl{db}
}
