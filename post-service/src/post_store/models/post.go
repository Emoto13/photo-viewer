package models

import "time"

type Post struct {
	Name      string    `cql:"name"`
	Path      string    `cql:"image_path"`
	Username  string    `cql:"username"`
	CreatedOn time.Time `cql:"created_on"`
}
