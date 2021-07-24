package setup

import (
	"context"

	"github.com/Emoto13/photo-viewer-rest/post-service/src/post_store/cache_store"
	"github.com/Emoto13/photo-viewer-rest/post-service/src/post_store/post_data"
)

type mockCacheStore struct {
}

func NewMockCacheStore() cache_store.PostCacheStore {
	return &mockCacheStore{}
}

func (s *mockCacheStore) Get(ctx context.Context, key string) ([]*post_data.PostData, error) {
	return []*post_data.PostData{}, nil
}

func (s *mockCacheStore) Set(ctx context.Context, key string, value []*post_data.PostData) error {
	return nil
}
