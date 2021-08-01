package post_data

import "time"

type PostData struct {
	Name      string    `cql:"name"`
	Path      string    `cql:"image_path"`
	Owner     string    `cql:"username"`
	CreatedOn time.Time `cql:"created_on"`
}
