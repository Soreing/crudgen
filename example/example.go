//go:build ignore

package main

import (
	"fmt"

	r "github.com/soreing/crudgen/example/repo"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	driv := "postgres"
	dsn := "host=localhost port=5432 user=postgres password=secret dbname=test sslmode=disable"

	db, err := sqlx.Open(driv, dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	ur := r.MakeUserRepo(db)
	usr, err := ur.Read("userid")
	fmt.Println(usr, err)
}
