package feed

import (
	"sort"

	postModels "github.com/Emoto13/photo-viewer-rest/feed-service/src/post/models"
)

type FeedSlice []*postModels.Post

func (p FeedSlice) Len() int {
	return len(p)
}

func (p FeedSlice) Less(i, j int) bool {
	return p[i].CreatedOn.Before(p[j].CreatedOn)
}

func (p FeedSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func sortFeed(feed FeedSlice) {
	sort.Sort(feed)
}
