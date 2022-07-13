package models

import "time"

type User struct {
	Id        string `gen:"order='1',pk,omitup"`
	Email     string
	FullName  string    `db:"full_name" gen:"name='full_name'"`
	CreatedAt time.Time `db:"created_at" gen:"name='created_at',omitup"`
	UpdatedAt time.Time `db:"updated_at" gen:"name='updated_at'"`
}
