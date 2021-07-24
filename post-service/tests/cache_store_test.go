package unit

import (
	"fmt"
	"testing"

	"github.com/Emoto13/photo-viewer-rest/post-service/tests/setup"
)

func TestCacheStore(t *testing.T) {
	store := setup.NewMockCacheStore()
	fmt.Println(store)
}
