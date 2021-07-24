package cache_store

import (
	"context"
	"time"

	"github.com/Emoto13/photo-viewer-rest/post-service/src/post_store/post_data"
	"github.com/go-redis/cache/v8"
)

type PostCacheStore interface {
	Set(ctx context.Context, key string, value []*post_data.PostData) error
	Get(ctx context.Context, key string) ([]*post_data.PostData, error)
}

type postCacheStore struct {
	store *cache.Cache
}

func NewPostCacheStore(store *cache.Cache) PostCacheStore {
	return &postCacheStore{store: store}
}

func (s *postCacheStore) Get(ctx context.Context, key string) ([]*post_data.PostData, error) {
	var result []*post_data.PostData
	err := s.store.Get(ctx, key, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *postCacheStore) Set(ctx context.Context, key string, value []*post_data.PostData) error {
	err := s.store.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: value,
		TTL:   15 * time.Minute,
	})
	if err != nil {
		return err
	}
	return err
}
