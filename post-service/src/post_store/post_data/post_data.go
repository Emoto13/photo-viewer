package post_data

import "time"

type PostData struct {
	Name      string
	Path      string
	Owner     string
	CreatedOn time.Time
}
