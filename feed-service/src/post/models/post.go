package models

import "time"

type Post struct {
	Username  string    `cql:"username"`
	Name      string    `cql:"name"`
	Path      string    `cql:"image_path"`
	CreatedOn time.Time `cql:"created_on"`
}
