package models

import "time"

type Post struct {
	Username  string
	Name      string
	Path      string
	CreatedOn time.Time
}
