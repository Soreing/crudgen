//go:build ignore

package main

import (
	m "github.com/soreing/crudgen/example/models"
	r "github.com/soreing/crudgen/example/repo"
	"reflect"

	gen "github.com/soreing/crudgen"
)

var configs = []gen.CRUDDescription{
	{
		TableName:   "users",
		RepoType:    reflect.TypeOf(r.UserRepoImpl{}),
		ModelType:   reflect.TypeOf(m.User{}),
		PackageName: "repo",
		SrcFilePath: "userCrud.go",
	},
}

func main() {
	gen.Run(configs)
}
