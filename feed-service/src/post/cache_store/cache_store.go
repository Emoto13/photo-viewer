package cache_store

import (
	"context"
	"time"

	"github.com/Emoto13/photo-viewer-rest/feed-service/src/post/models"
	"github.com/go-redis/cache/v8"
)

type PostCacheStore interface {
	Set(ctx context.Context, key string, value []*models.Post) error
	Get(ctx context.Context, key string) ([]*models.Post, error)
}

type postCacheStore struct {
	store *cache.Cache
}

func NewPostCacheStore(store *cache.Cache) PostCacheStore {
	return &postCacheStore{store: store}
}

func (s *postCacheStore) Get(ctx context.Context, key string) ([]*models.Post, error) {
	var result []*models.Post
	err := s.store.Get(ctx, key, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *postCacheStore) Set(ctx context.Context, key string, value []*models.Post) error {
	err := s.store.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: value,
		TTL:   3 * time.Minute,
	})
	if err != nil {
		return err
	}
	return err
}
