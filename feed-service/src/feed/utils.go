package feed

import (
	"sort"

	postModels "github.com/Emoto13/photo-viewer-rest/feed-service/src/post/models"
)

func (p []*postModels.Post) Len() int {
	return len(p)
}

func (p []*postModels.Post) Less(i, j int) bool {
	return p[i].CreatedOn.Before(p[j].CreatedOn)
}

func (p []*postModels.Post) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func sortPosts(posts []*postModels.Post) {
	sort.Sort(posts)
}
