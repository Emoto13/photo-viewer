package setup

import (
	"context"

	"github.com/Emoto13/photo-viewer-rest/post-service/src/post_store/cache_store"
	"github.com/Emoto13/photo-viewer-rest/post-service/src/post_store/models"
)

type mockCacheStore struct {
}

func NewMockCacheStore() cache_store.PostCacheStore {
	return &mockCacheStore{}
}

func (s *mockCacheStore) Get(ctx context.Context, key string) ([]*models.Post, error) {
	return []*models.Post{}, nil
}

func (s *mockCacheStore) Set(ctx context.Context, key string, value []*models.Post) error {
	return nil
}
